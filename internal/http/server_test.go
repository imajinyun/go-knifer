package http

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Mirrors the action routing example from hutool-http server/SimpleServerTest.

func TestSimpleServerStartAndStop(t *testing.T) {
	port := pickFreePort(t)
	srv := NewSimpleServer(port)

	srv.AddAction("/get", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Path))
	})
	srv.AddAction("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write(body)
	})
	srv.AddAction("/zero", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("0"))
	})

	errCh := srv.StartAsync()
	defer func() {
		_ = srv.Stop(2 * time.Second)
		select {
		case err := <-errCh:
			if err != nil {
				t.Logf("server err: %v", err)
			}
		case <-time.After(2 * time.Second):
		}
	}()

	waitServerReady(t, port)
	base := "http://127.0.0.1:" + strconv.Itoa(port)

	if got := Get(base + "/get").Execute().Body(); got != "/get" {
		t.Fatalf("get: %q", got)
	}

	body := Post(base + "/echo").BodyString(`{"a":1}`).Execute().Body()
	if body != `{"a":1}` {
		t.Fatalf("echo: %q", body)
	}

	if got := Get(base + "/zero").Execute().Body(); got != "0" {
		t.Fatalf("zero: %q", got)
	}
}

func TestCreateServerHelper(t *testing.T) {
	srv := CreateServer(0)
	if srv == nil {
		t.Fatal("nil")
	}
}

// pickFreePort reserves a free port, releases it immediately, and returns the port number.
func pickFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return port
}

// waitServerReady polls until the port is connectable.
func waitServerReady(t *testing.T, port int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(port), 100*time.Millisecond)
		if err == nil {
			_ = c.Close()
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("server on %d not ready", port)
}
