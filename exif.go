// Copyright (c) 2012-2015 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package exif provides bindings for libexif.
package exif

/*
#include <stdlib.h>
#include <libexif/exif-data.h>
#include <libexif/exif-loader.h>
#include "_cgo/types.h"

exif_value_t* pop_exif_value(exif_stack_t *);
void free_exif_value(exif_value_t* n);
exif_stack_t* exif_dump(ExifData *);
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// Error messages.
var (
	ErrNoExifData      = errors.New(`No EXIF data found.`)
	ErrFoundExifInData = errors.New(`Found EXIF header. OK to call Parse.`)
)

// Data stores the EXIF tags of a file.
type Data struct {
	exifLoader *C.ExifLoader
	Tags       map[string]string
}

// New creates and returns a new exif.Data object.
func New() *Data {
	data := &Data{
		Tags: make(map[string]string),
	}
	return data
}

// Read attempts to read EXIF data from a file.
func Read(file string) (*Data, error) {
	data := New()
	if err := data.Open(file); err != nil {
		return nil, err
	}
	return data, nil
}

// Open opens a file path and loads its EXIF data.
func (d *Data) Open(file string) error {

	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	exifData := C.exif_data_new_from_file(cfile)

	if exifData == nil {
		return ErrNoExifData
	}
	defer C.exif_data_unref(exifData)

	return d.parseExifData(exifData)
}

func (d *Data) parseExifData(exifData *C.ExifData) error {
	values := C.exif_dump(exifData)
	defer C.free(unsafe.Pointer(values))

	for {
		value := C.pop_exif_value(values)
		if value == nil {
			break
		} else {
			d.Tags[strings.Trim(C.GoString((*value).name), " ")] = strings.Trim(C.GoString((*value).value), " ")
		}
		C.free_exif_value(value)
	}

	return nil
}

// Write writes bytes to the exif loader. Sends ErrFoundExifInData error when
// enough bytes have been sent.
func (d *Data) Write(p []byte) (n int, err error) {
	if d.exifLoader == nil {
		d.exifLoader = C.exif_loader_new()
		runtime.SetFinalizer(d, (*Data).cleanup)
	}

	res := C.exif_loader_write(d.exifLoader, (*C.uchar)(unsafe.Pointer(&p[0])), C.uint(len(p)))

	if res == 1 {
		return len(p), nil
	}
	return len(p), ErrFoundExifInData
}

// Parse finalizes the data loader and sets the tags
func (d *Data) Parse() error {
	defer d.cleanup()

	exifData := C.exif_loader_get_data(d.exifLoader)
	if exifData == nil {
		return fmt.Errorf(ErrNoExifData.Error(), "")
	}

	defer func() {
		C.exif_data_unref(exifData)
	}()

	return d.parseExifData(exifData)
}

func (d *Data) cleanup() {
	if d.exifLoader != nil {
		C.exif_loader_unref(d.exifLoader)
		d.exifLoader = nil
	}
}
