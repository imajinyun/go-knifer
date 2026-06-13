package zip

import (
	"bytes"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestZipErrorContract(t *testing.T) {
	_, err := GetStream(nil)
	assertZipCode(t, err, knifer.ErrCodeInvalidInput)
	assertZipCode(t, UnzipReaderToLimit(nil, t.TempDir(), -1), knifer.ErrCodeInvalidInput)

	var buf bytes.Buffer
	err = ZipEntriesToWriter(&buf, EntryData{Name: "../evil.txt", Data: []byte("bad")})
	assertZipCode(t, err, knifer.ErrCodeInvalidInput)
}
