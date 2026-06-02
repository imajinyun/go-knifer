package vnet

import (
	"crypto/tls"
	"crypto/x509"
	"io/fs"
	"math/big"
	"mime/multipart"
	stdnet "net"
	"net/http"
	"net/url"
	"time"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
)

const (
	LocalIP         = netimpl.LocalIP
	IPSplitMark     = netimpl.IPSplitMark
	IPMaskSplitMark = netimpl.IPMaskSplitMark
	IPMaskMax       = netimpl.IPMaskMax
	PortRangeMin    = netimpl.PortRangeMin
	PortRangeMax    = netimpl.PortRangeMax
	SSL             = netimpl.SSL
	SSLv2           = netimpl.SSLv2
	SSLv3           = netimpl.SSLv3
	TLS             = netimpl.TLS
	TLSv1           = netimpl.TLSv1
	TLSv11          = netimpl.TLSv11
	TLSv12          = netimpl.TLSv12
	TLSv13          = netimpl.TLSv13
)

type (
	URLBuilder         = netimpl.URLBuilder
	TLSConfigBuilder   = netimpl.TLSConfigBuilder
	UploadSetting      = netimpl.UploadSetting
	UploadSaveOption   = netimpl.UploadSaveOption
	MultipartFormData  = netimpl.MultipartFormData
	LocalPortGenerator = netimpl.LocalPortGenerator
)

func LongToIPv4(longIP uint32) string         { return netimpl.LongToIPv4(longIP) }
func IPv4ToLong(strIP string) (uint32, error) { return netimpl.IPv4ToLong(strIP) }
func IPv4ToLongDefault(strIP string, defaultValue uint32) uint32 {
	return netimpl.IPv4ToLongDefault(strIP, defaultValue)
}
func IPv6ToBigInt(ipv6Str string) (*big.Int, error)      { return netimpl.IPv6ToBigInt(ipv6Str) }
func BigIntToIPv6(n *big.Int) (string, error)            { return netimpl.BigIntToIPv6(n) }
func IsIP(s string) bool                                 { return netimpl.IsIP(s) }
func IsIPv4(s string) bool                               { return netimpl.IsIPv4(s) }
func IsIPv6(s string) bool                               { return netimpl.IsIPv6(s) }
func IsInnerIP(ipAddress string) bool                    { return netimpl.IsInnerIP(ipAddress) }
func FormatIPBlock(ip, mask string) (string, error)      { return netimpl.FormatIPBlock(ip, mask) }
func BeginIP(ip string, maskBit int) (string, error)     { return netimpl.BeginIP(ip, maskBit) }
func BeginIPLong(ip string, maskBit int) (uint32, error) { return netimpl.BeginIPLong(ip, maskBit) }
func EndIP(ip string, maskBit int) (string, error)       { return netimpl.EndIP(ip, maskBit) }
func EndIPLong(ip string, maskBit int) (uint32, error)   { return netimpl.EndIPLong(ip, maskBit) }
func MaskBitByMask(mask string) (int, error)             { return netimpl.MaskBitByMask(mask) }
func CountByMaskBit(maskBit int, isAll bool) (uint64, error) {
	return netimpl.CountByMaskBit(maskBit, isAll)
}
func MaskByMaskBit(maskBit int) (string, error)         { return netimpl.MaskByMaskBit(maskBit) }
func MaskByIPRange(fromIP, toIP string) (string, error) { return netimpl.MaskByIPRange(fromIP, toIP) }

func CountByIPRange(fromIP, toIP string) (uint64, error) {
	return netimpl.CountByIPRange(fromIP, toIP)
}
func IsMaskValid(mask string) bool    { return netimpl.IsMaskValid(mask) }
func IsMaskBitValid(maskBit int) bool { return netimpl.IsMaskBitValid(maskBit) }

func ListIPs(ipRange string, isAll bool) ([]string, error) { return netimpl.ListIPs(ipRange, isAll) }

func ListIPCIDR(ip string, maskBit int, isAll bool) ([]string, error) {
	return netimpl.ListIPCIDR(ip, maskBit, isAll)
}
func ListIPRange(fromIP, toIP string) ([]string, error) { return netimpl.ListIPRange(fromIP, toIP) }
func MatchesWildcard(wildcard, ipAddress string) bool {
	return netimpl.MatchesWildcard(wildcard, ipAddress)
}
func IsInRange(ip, cidr string) bool { return netimpl.IsInRange(ip, cidr) }

func Decode(s string) (string, error)        { return netimpl.Decode(s) }
func DecodeForPath(s string) (string, error) { return netimpl.DecodeForPath(s) }
func DecodePlus(s string, plusToSpace bool) (string, error) {
	return netimpl.DecodePlus(s, plusToSpace)
}
func EncodeAll(s string) string         { return netimpl.EncodeAll(s) }
func Encode(s string) string            { return netimpl.Encode(s) }
func EncodeQuery(s string) string       { return netimpl.EncodeQuery(s) }
func EncodePathSegment(s string) string { return netimpl.EncodePathSegment(s) }
func EncodePath(s string) string        { return netimpl.EncodePath(s) }
func EncodeFragment(s string) string    { return netimpl.EncodeFragment(s) }
func FormURLEncode(s string) string     { return netimpl.FormURLEncode(s) }

func NewURLBuilder() *URLBuilder                      { return netimpl.NewURLBuilder() }
func NewHTTPURLBuilder(host string) *URLBuilder       { return netimpl.NewHTTPURLBuilder(host) }
func ParseURLBuilder(raw string) (*URLBuilder, error) { return netimpl.ParseURLBuilder(raw) }

func IsValidPort(port int) bool                       { return netimpl.IsValidPort(port) }
func IsUsableLocalPort(port int) bool                 { return netimpl.IsUsableLocalPort(port) }
func GetUsableLocalPort() (int, error)                { return netimpl.GetUsableLocalPort() }
func GetUsableLocalPortFrom(minPort int) (int, error) { return netimpl.GetUsableLocalPortFrom(minPort) }

func GetUsableLocalPortInRange(minPort, maxPort int) (int, error) {
	return netimpl.GetUsableLocalPortInRange(minPort, maxPort)
}

func GetUsableLocalPorts(numRequested, minPort, maxPort int) ([]int, error) {
	return netimpl.GetUsableLocalPorts(numRequested, minPort, maxPort)
}

func NewLocalPortGenerator(beginPort int) *LocalPortGenerator {
	return netimpl.NewLocalPortGenerator(beginPort)
}
func HideIPPart(ip string) string     { return netimpl.HideIPPart(ip) }
func HideIPPartLong(ip uint32) string { return netimpl.HideIPPartLong(ip) }
func BuildInetSocketAddress(host string, defaultPort int) (*stdnet.TCPAddr, error) {
	return netimpl.BuildInetSocketAddress(host, defaultPort)
}

func CreateAddress(host string, port int) *stdnet.TCPAddr { return netimpl.CreateAddress(host, port) }
func GetIPByHost(hostName string) string                  { return netimpl.GetIPByHost(hostName) }
func GetNetworkInterface(name string) (*stdnet.Interface, error) {
	return netimpl.GetNetworkInterface(name)
}
func GetNetworkInterfaces() ([]stdnet.Interface, error) { return netimpl.GetNetworkInterfaces() }
func LocalIPv4s() []string                              { return netimpl.LocalIPv4s() }
func LocalIPv6s() []string                              { return netimpl.LocalIPv6s() }
func LocalIPs() []string                                { return netimpl.LocalIPs() }
func ToIPList(addressList []stdnet.IP) []string         { return netimpl.ToIPList(addressList) }
func LocalAddressList(addressFilter func(stdnet.IP) bool) []stdnet.IP {
	return netimpl.LocalAddressList(addressFilter)
}

func LocalAddressListByInterface(interfaceFilter func(stdnet.Interface) bool, addressFilter func(stdnet.IP) bool) []stdnet.IP {
	return netimpl.LocalAddressListByInterface(interfaceFilter, addressFilter)
}

func GetLocalhostStr() string                       { return netimpl.GetLocalhostStr() }
func GetLocalhost() stdnet.IP                       { return netimpl.GetLocalhost() }
func GetLocalHostName() string                      { return netimpl.GetLocalHostName() }
func GetLocalMACAddress(separator ...string) string { return netimpl.GetLocalMACAddress(separator...) }

func GetMACAddress(inetAddress stdnet.IP, separator ...string) string {
	return netimpl.GetMACAddress(inetAddress, separator...)
}

func GetHardwareAddress(inetAddress stdnet.IP) stdnet.HardwareAddr {
	return netimpl.GetHardwareAddress(inetAddress)
}

func GetLocalHardwareAddress() stdnet.HardwareAddr { return netimpl.GetLocalHardwareAddress() }
func NetCat(host string, port int, data []byte, timeout time.Duration) error {
	return netimpl.NetCat(host, port, data, timeout)
}

func Ping(ip string, timeout time.Duration) bool { return netimpl.Ping(ip, timeout) }
func IsOpen(address *stdnet.TCPAddr, timeout time.Duration) bool {
	return netimpl.IsOpen(address, timeout)
}

func IDNToASCII(unicode string) (string, error)    { return netimpl.IDNToASCII(unicode) }
func GetMultistageReverseProxyIP(ip string) string { return netimpl.GetMultistageReverseProxyIP(ip) }
func IsUnknown(checkString string) bool            { return netimpl.IsUnknown(checkString) }
func ParseCookies(cookieStr string) []*http.Cookie { return netimpl.ParseCookies(cookieStr) }
func GetDNSInfo(hostName string, attrNames ...string) ([]string, error) {
	return netimpl.GetDNSInfo(hostName, attrNames...)
}

func Connect(hostname string, port int, timeout time.Duration) (stdnet.Conn, error) {
	return netimpl.Connect(hostname, port, timeout)
}

func GetRemoteAddress(conn stdnet.Conn) string { return netimpl.GetRemoteAddress(conn) }
func IsConnected(conn stdnet.Conn) bool        { return netimpl.IsConnected(conn) }

func NewTLSConfigBuilder() *TLSConfigBuilder { return netimpl.NewTLSConfigBuilder() }
func CreateTLSConfig(insecureSkipVerify bool) *tls.Config {
	return netimpl.CreateTLSConfig(insecureSkipVerify)
}

func InsecureTLSConfig() *tls.Config    { return netimpl.InsecureTLSConfig() }
func TLSVersion(protocol string) uint16 { return netimpl.TLSVersion(protocol) }
func NewCertPool() *x509.CertPool       { return x509.NewCertPool() }

func NewUploadSetting() UploadSetting { return netimpl.NewUploadSetting() }
func ParseMultipartForm(r *http.Request, setting UploadSetting) (*MultipartFormData, error) {
	return netimpl.ParseMultipartForm(r, setting)
}

func WithUploadFilePerm(perm fs.FileMode) UploadSaveOption { return netimpl.WithUploadFilePerm(perm) }

func WithUploadDirPerm(perm fs.FileMode) UploadSaveOption { return netimpl.WithUploadDirPerm(perm) }

func WithUploadOverwrite(overwrite bool) UploadSaveOption {
	return netimpl.WithUploadOverwrite(overwrite)
}

func WithUploadCreateParents(create bool) UploadSaveOption {
	return netimpl.WithUploadCreateParents(create)
}

func SaveUploadedFile(file *multipart.FileHeader, destPath string, opts ...UploadSaveOption) error {
	return netimpl.SaveUploadedFile(file, destPath, opts...)
}

func UploadFileName(file *multipart.FileHeader) string { return netimpl.UploadFileName(file) }
func UploadFileSize(file *multipart.FileHeader) int64  { return netimpl.UploadFileSize(file) }
func UploadFileContentType(file *multipart.FileHeader) string {
	return netimpl.UploadFileContentType(file)
}

// URLValues creates a URL query value map.
func URLValues() url.Values { return url.Values{} }
