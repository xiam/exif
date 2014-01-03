package exif

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	err := exif.Open("_examples/resources/test.jpg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("----- Open")
	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}

func TestWriteAndParse(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	file, err := os.Open("_examples/resources/test.jpg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	defer file.Close()

	_, err = io.Copy(exif, file)

	if err != nil && err != FoundExifInData {
		t.Fatalf("Error: %s", err.Error())
	}

	err = exif.Parse()

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("----- Write and Parse")
	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}
