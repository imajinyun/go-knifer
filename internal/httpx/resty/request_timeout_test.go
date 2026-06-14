package resty

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Timeout(time.Millisecond).Execute()
	if resp.Err() == nil {
		t.Fatal("Execute() error is nil, want timeout error")
	}
}
