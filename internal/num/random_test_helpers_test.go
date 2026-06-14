package num

import "errors"

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("forced random failure") }

type sequenceReader struct {
	next byte
}

func (r *sequenceReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = r.next
		r.next++
	}
	return len(p), nil
}
