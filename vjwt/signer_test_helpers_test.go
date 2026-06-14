package vjwt_test

import "io"

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 1
	}
	return len(p), nil
}

var _ io.Reader = zeroReader{}
