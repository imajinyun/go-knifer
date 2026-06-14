package resty

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestJSONMarshalProvider(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	called := false
	resp := Post(srv.URL, WithJSONMarshalFunc(func(any) ([]byte, error) {
		called = true
		return []byte(`{"provided":true}`), nil
	})).BodyJSONValue(map[string]any{"ignored": true}).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !called || resp.Body() != `{"provided":true}` {
		t.Fatalf("marshal provider called=%v body=%q", called, resp.Body())
	}
}

func TestRequestJSONUnmarshalProvider(t *testing.T) {
	type result struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"ignored"}`))
	}))
	defer srv.Close()

	called := false
	out := &result{}
	resp := Get(srv.URL, WithJSONUnmarshalFunc(func(_ []byte, dst any) error {
		called = true
		return json.Unmarshal([]byte(`{"name":"provided"}`), dst)
	})).Result(out).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !called || out.Name != "provided" || resp.Result() == nil {
		t.Fatalf("unmarshal provider called=%v result=%+v raw=%v", called, out, resp.Result())
	}
}

func TestRequestJSONDecodeReadOptions(t *testing.T) {
	type result struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"abcdef"}`))
	}))
	defer srv.Close()

	tooLarge := &result{}
	resp := Get(srv.URL,
		WithMaxDecodeBytes(3),
		WithJSONUnmarshalFunc(json.Unmarshal),
	).Result(tooLarge).Execute()
	if resp.Err() == nil {
		t.Fatal("Execute() with max decode bytes error = nil")
	}

	readAllCalled := false
	out := &result{}
	resp = Get(srv.URL,
		WithJSONDecodeReadAllFunc(func(io.Reader) ([]byte, error) {
			readAllCalled = true
			return []byte(`{"name":"provided"}`), nil
		}),
		WithJSONUnmarshalFunc(json.Unmarshal),
	).Result(out).Execute()
	if resp.Err() != nil || !readAllCalled || out.Name != "provided" {
		t.Fatalf("custom decode readAll called=%v out=%+v err=%v", readAllCalled, out, resp.Err())
	}
}
