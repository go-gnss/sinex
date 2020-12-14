package sinex_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-gnss/sinex"
)

func TestParse(t *testing.T) {
	f, err := os.Open("fixtures/apr20187.snx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	sf, err := sinex.Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", sf)
}
