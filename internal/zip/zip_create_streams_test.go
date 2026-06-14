package zip

import (
	"bytes"
	"testing"
)

func TestZipStreamsEnforcesInputLimit(t *testing.T) {
	var buf bytes.Buffer
	err := ZipStreamsToWriterWithOptions(&buf, []StreamEntry{{Name: "a.txt", Reader: bytes.NewReader([]byte("abcd"))}}, WithMaxBytes(3))
	if err == nil {
		t.Fatal("ZipStreamsToWriterWithOptions should reject stream input over max bytes")
	}
}
