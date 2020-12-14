package sinex_test

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/go-gnss/sinex"
)

func TestParse(t *testing.T) {
	f, err := os.Open("fixtures/apr20187.snx.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sinex.Parse(r)
	if err != nil {
		t.Fatal(err)
	}
}
