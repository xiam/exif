package exif

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	exif := New()

	err := exif.Open("test.jpg")

	if err != nil {
		panic(err)
	}

	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}
