package conf

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithOptionsDecryptsIncludedConfigAndValidatesSchema(t *testing.T) {
	dir := t.TempDir()
	common := filepath.Join(dir, "common.setting")
	main := filepath.Join(dir, "main.setting")
	secret := base64.StdEncoding.EncodeToString([]byte("s3cr3t"))
	if err := os.WriteFile(common, []byte("[server]\nhost=127.0.0.1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include=common.setting\nname=main\nsecret=ENC(base64:"+secret+")\n[server]\nport=8080"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadWithOptions(main, LoadOptions{AllowInclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("secret"); got != "s3cr3t" {
		t.Fatalf("decrypted secret = %q", got)
	}
	if err := c.ValidateSchema(Schema{Fields: []FieldRule{
		{Key: "name", Required: true, Type: TypeString},
		{Group: "server", Key: "port", Required: true, Type: TypeInt},
		{Group: "server", Key: "host", Required: true},
	}}); err != nil {
		t.Fatalf("ValidateSchema() error = %v", err)
	}
	if err := c.ValidateSchema(Schema{Fields: []FieldRule{{Group: "server", Key: "debug", Required: true}}}); err == nil {
		t.Fatal("ValidateSchema() missing required error = nil")
	}
}
