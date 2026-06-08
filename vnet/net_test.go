package vnet_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	stdnet "net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vnet"
)

type recordingDialer struct {
	network string
	address string
	data    chan []byte
}

type stubListener struct{}

func (stubListener) Accept() (stdnet.Conn, error) { return nil, errors.New("stub listener") }
func (stubListener) Close() error                 { return nil }
func (stubListener) Addr() stdnet.Addr {
	return &stdnet.TCPAddr{IP: stdnet.ParseIP("127.0.0.1"), Port: 12345}
}

func (d *recordingDialer) DialContext(_ context.Context, network, address string) (stdnet.Conn, error) {
	d.network = network
	d.address = address
	client, server := stdnet.Pipe()
	go func() {
		defer func() { _ = server.Close() }()
		payload, _ := io.ReadAll(server)
		d.data <- payload
	}()
	return client, nil
}

func TestVNetFacade(t *testing.T) {
	v, err := vnet.IPv4ToLong("127.0.0.1")
	if err != nil || vnet.LongToIPv4(v) != "127.0.0.1" {
		t.Fatalf("IPv4 facade failed: %d %v", v, err)
	}
	if !vnet.IsIPv4("192.168.1.1") || !vnet.IsIPv6("::1") || !vnet.IsInnerIP("10.0.0.1") {
		t.Fatal("IP validators failed")
	}
	if !vnet.IsValidPort(80) || vnet.HideIPPart("192.168.1.2") != "192.168.1.*" {
		t.Fatal("port or hide helper failed")
	}
	if vnet.CreateTLSConfig() == nil || vnet.NewUploadSetting().MemoryThreshold == 0 {
		t.Fatal("TLS/upload helpers failed")
	}
}

func TestVNetFacadeOptions(t *testing.T) {
	ln, err := stdnet.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen local port: %v", err)
	}
	defer func() { _ = ln.Close() }()

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := ln.Accept()
		if err == nil {
			_ = conn.Close()
		}
	}()

	_, portStr, err := stdnet.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("split listener address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse listener port: %v", err)
	}
	if !vnet.PingWithOptions("127.0.0.1", vnet.WithPingPorts(port), vnet.WithPingTimeout(time.Second), vnet.WithPingNetwork("tcp")) {
		t.Fatal("PingWithOptions should reach local listener")
	}
	<-done

	if vnet.IsUsableLocalPortWithOptions(port, vnet.WithPortHost("127.0.0.1")) {
		t.Fatal("IsUsableLocalPortWithOptions should reject occupied port")
	}
	g := vnet.NewLocalPortGeneratorWithOptions(port, vnet.WithPortHost("127.0.0.1"))
	generated, err := g.Generate()
	if err != nil {
		t.Fatalf("LocalPortGenerator.Generate with options: %v", err)
	}
	if generated <= port || generated > vnet.PortRangeMax {
		t.Fatalf("LocalPortGenerator generated %d, want > %d", generated, port)
	}
	freePort, err := vnet.GetUsableLocalPortInRangeWithOptions(port+1, port+20, vnet.WithPortHost("127.0.0.1"))
	if err != nil || freePort < port+1 || freePort > port+20 {
		t.Fatalf("GetUsableLocalPortInRangeWithOptions = %d, %v", freePort, err)
	}
	ports, err := vnet.GetUsableLocalPortsWithOptions(1, port+1, port+20, vnet.WithPortHost("127.0.0.1"))
	if err != nil || len(ports) != 1 {
		t.Fatalf("GetUsableLocalPortsWithOptions = %v, %v; want one port", ports, err)
	}
	ips, err := vnet.GetIPByHostWithOptions("localhost", vnet.WithResolveNetwork("ip4"), vnet.WithResolveTimeout(time.Second))
	if err != nil || len(ips) == 0 {
		t.Fatalf("GetIPByHostWithOptions = %v, %v; want at least one IPv4", ips, err)
	}
	dns, err := vnet.GetDNSInfoWithOptions("localhost", vnet.WithDNSTypes("A"), vnet.WithResolveTimeout(time.Second))
	if err != nil || len(dns) == 0 {
		t.Fatalf("GetDNSInfoWithOptions = %v, %v; want at least one A record", dns, err)
	}
}

func TestVNetProviderOptionsFacade(t *testing.T) {
	var network, address string
	addr, err := vnet.BuildInetSocketAddressWithOptions("example.com", 8080, vnet.WithAddressNetwork("tcp4"), vnet.WithTCPAddrResolver(func(n, a string) (*stdnet.TCPAddr, error) {
		network, address = n, a
		return &stdnet.TCPAddr{IP: stdnet.ParseIP("10.0.0.2"), Port: 8080}, nil
	}))
	if err != nil || addr.Port != 8080 {
		t.Fatalf("BuildInetSocketAddressWithOptions = %#v %v", addr, err)
	}
	if network != "tcp4" || address != "example.com:8080" {
		t.Fatalf("address resolver target = %s %s", network, address)
	}

	if !vnet.IsUsableLocalPortWithOptions(23456, vnet.WithPortNetwork("tcp4"), vnet.WithPortHost("127.0.0.2"), vnet.WithPortListenerFactory(func(n, a string) (stdnet.Listener, error) {
		network, address = n, a
		return stubListener{}, nil
	})) {
		t.Fatal("IsUsableLocalPortWithOptions should use listener factory")
	}
	if network != "tcp4" || address != "127.0.0.2:23456" {
		t.Fatalf("listener target = %s %s", network, address)
	}

	iface := stdnet.Interface{Name: "vnet0", HardwareAddr: stdnet.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}}
	_, ipNet, err := stdnet.ParseCIDR("10.9.8.7/24")
	if err != nil {
		t.Fatal(err)
	}
	ipNet.IP = stdnet.ParseIP("10.9.8.7")
	opts := []vnet.InterfaceOption{
		vnet.WithInterfaceByNameFunc(func(name string) (*stdnet.Interface, error) { return &iface, nil }),
		vnet.WithInterfacesFunc(func() ([]stdnet.Interface, error) { return []stdnet.Interface{iface}, nil }),
		vnet.WithInterfaceAddrsFunc(func(stdnet.Interface) ([]stdnet.Addr, error) { return []stdnet.Addr{ipNet}, nil }),
		vnet.WithReverseLookupFunc(func(string) ([]string, error) { return []string{"vnet.local."}, nil }),
		vnet.WithNetHostnameFunc(func() (string, error) { return "fallback", nil }),
	}
	gotIface, err := vnet.GetNetworkInterfaceWithOptions("vnet0", opts...)
	if err != nil || gotIface.Name != "vnet0" {
		t.Fatalf("GetNetworkInterfaceWithOptions = %#v %v", gotIface, err)
	}
	if got := vnet.LocalIPv4sWithOptions(opts...); len(got) != 1 || got[0] != "10.9.8.7" {
		t.Fatalf("LocalIPv4sWithOptions = %#v", got)
	}
	if got := vnet.GetLocalHostNameWithOptions(opts...); got != "vnet.local" {
		t.Fatalf("GetLocalHostNameWithOptions = %q", got)
	}
	if got := vnet.GetLocalMACAddressWithOptions(opts, "-"); got != "01-02-03-04-05-06" {
		t.Fatalf("GetLocalMACAddressWithOptions = %q", got)
	}
}

func TestVNetConnectOptionsFacade(t *testing.T) {
	dialer := &recordingDialer{data: make(chan []byte, 1)}
	conn, err := vnet.ConnectWithOptions(
		"example.com", 8080,
		vnet.WithConnectNetwork("tcp4"),
		vnet.WithConnectTimeout(time.Second),
		vnet.WithConnectDialer(dialer),
	)
	if err != nil {
		t.Fatalf("ConnectWithOptions: %v", err)
	}
	_ = conn.Close()
	if dialer.network != "tcp4" || dialer.address != "example.com:8080" {
		t.Fatalf("dial target = %s %s", dialer.network, dialer.address)
	}

	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if err := vnet.NetCatWithOptions("127.0.0.1", 1234, []byte("hello"), vnet.WithConnectDialer(dialer)); err != nil {
		t.Fatalf("NetCatWithOptions: %v", err)
	}
	if got := string(<-dialer.data); got != "hello" {
		t.Fatalf("NetCatWithOptions wrote %q", got)
	}

	addr := &stdnet.TCPAddr{IP: stdnet.ParseIP("127.0.0.1"), Port: 4321}
	dialer = &recordingDialer{data: make(chan []byte, 1)}
	if !vnet.IsOpenWithOptions(addr, vnet.WithConnectDialer(dialer)) {
		t.Fatal("IsOpenWithOptions should report true")
	}
}

func TestVNetUploadSaveOptionsFacade(t *testing.T) {
	req := multipartRequest(t, "avatar", "a.txt", "hello")
	form, err := vnet.ParseMultipartForm(req, vnet.NewUploadSetting())
	if err != nil {
		t.Fatalf("ParseMultipartForm: %v", err)
	}
	file := form.GetFile("avatar")
	if file == nil {
		t.Fatal("uploaded file is nil")
	}
	if vnet.UploadFileName(file) != "a.txt" || vnet.UploadFileSize(file) != int64(len("hello")) || vnet.UploadFileContentType(file) == "" {
		t.Fatalf("upload metadata = name:%q size:%d type:%q", vnet.UploadFileName(file), vnet.UploadFileSize(file), vnet.UploadFileContentType(file))
	}

	dir := t.TempDir()
	dest := filepath.Join(dir, "nested", "a.txt")
	if err := vnet.SaveUploadedFile(file, dest, vnet.WithUploadFilePerm(0o600), vnet.WithUploadDirPerm(0o700)); err != nil {
		t.Fatalf("SaveUploadedFile: %v", err)
	}
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("stat saved file: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("saved file perm = %v", info.Mode().Perm())
	}
	if err := vnet.SaveUploadedFile(file, dest, vnet.WithUploadOverwrite(false)); err == nil {
		t.Fatal("SaveUploadedFile should reject overwrite when disabled")
	}
	missingParent := filepath.Join(dir, "missing", "b.txt")
	if err := vnet.SaveUploadedFile(file, missingParent, vnet.WithUploadCreateParents(false)); err == nil {
		t.Fatal("SaveUploadedFile should reject missing parent when parent creation is disabled")
	}

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	err = vnet.SaveUploadedFile(file, "/virtual/upload/a.txt",
		vnet.WithUploadMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		vnet.WithUploadOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		vnet.WithUploadDirPerm(0o700), vnet.WithUploadFilePerm(0o600),
	)
	if err != nil {
		t.Fatalf("SaveUploadedFile provider: %v", err)
	}
	if mkdirPath != "/virtual/upload" || mkdirPerm != 0o700 || openPath != "/virtual/upload/a.txt" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != "hello" {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}

func TestVNetTLSFileOptionsFacade(t *testing.T) {
	const certPEM = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIRAPWQSq0Qr7yZD5twH61BxFIwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHZ28tdGVzdDAeFw0yNjA2MDYwMDAwMDBaFw0yNzA2MDYwMDAwMDBa
MBIxEDAOBgNVBAoTB2dvLXRlc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASm
1YPqMC7UTw4R7ovbHYgk4+LALoU6hr61VnsBiKCdsMCMScpLob8ldIl+6o4f/ntM
5kmXvEFd9Mp6FfaHkgnbo0IwQDAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQUX90U1OkOXbGUzD2JNoWlqQtk3/0wCgYIKoZIzj0EAwID
SQAwRgIhANw7UzN0vtxOfygWqANg00uGOo7y98q1/Ac3N1wQxVBkAiEA7QjQRHtH
LA6wKo8yoCnW36b+nvxlhHvzrIxwWCgwCWM=
-----END CERTIFICATE-----`
	readPath := ""
	b := vnet.NewTLSConfigBuilder()
	if err := b.AddRootCAFileWithOptions("ca.pem", vnet.WithTLSReadFile(func(path string) ([]byte, error) {
		readPath = path
		return []byte(certPEM), nil
	})); err != nil {
		t.Fatalf("AddRootCAFileWithOptions: %v", err)
	}
	if readPath != "ca.pem" || b.Build().RootCAs == nil {
		t.Fatalf("TLS read provider not applied path=%q cfg=%#v", readPath, b.Build())
	}

	b = vnet.NewTLSConfigBuilder()
	if err := b.AddRootCAReader(strings.NewReader(certPEM)); err != nil || b.Build().RootCAs == nil {
		t.Fatalf("AddRootCAReader rootCAs=%#v err=%v", b.Build().RootCAs, err)
	}
}

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

func multipartRequest(t *testing.T, field, filename, content string) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}
