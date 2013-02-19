package exif

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	err := exif.Open("_examples/resources/test.jpg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}
