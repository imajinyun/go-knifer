package errx

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEmptyFormatterSuppressesOutput(t *testing.T) {
	data, err := EmptyFormatter.Format(logrus.NewEntry(logrus.New()))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Fatalf("EmptyFormatter output length = %d, want 0", len(data))
	}
}
