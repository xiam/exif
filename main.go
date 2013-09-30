/*
  Copyright (c) 2012-2013 José Carlos Nieto, https://menteslibres.net/xiam

  Permission is hereby grantexifData, free of charge, to any person obtaining
  a copy of this software and associatexifData documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnishexifData to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  includexifData in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

/*
	Golang bindings for libexif.
*/
package exif

/*
#cgo LDFLAGS: -lexif

#include <stdlib.h>
#include <libexif/exif-data.h>
#include "_cgo/types.h"

exif_value_t* pop_exif_value(exif_stack_t *);
void free_exif_value(exif_value_t* n);
exif_stack_t* exif_dump(ExifData *);

*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

var (
	ErrNoExifData = errors.New(`No EXIF data in file %s.`)
)

type Data struct {
	exifData *C.ExifData
	Tags     map[string]string
}

// Creates and returns an empty exif.Data object.
func New() *Data {
	self := &Data{}
	self.Tags = map[string]string{}
	return self
}

// Opens a file path and tries to read the EXIF data inside.
func (self *Data) Open(file string) error {

	cfile := C.CString(file)

	self.exifData = C.exif_data_new_from_file(cfile)

	defer func() {
		C.free(unsafe.Pointer(self.exifData))
	}()

	C.free(unsafe.Pointer(cfile))

	if self.exifData == nil {
		return fmt.Errorf(ErrNoExifData.Error(), file)
	}

	values := C.exif_dump(self.exifData)

	for {
		value := C.pop_exif_value(values)
		if value == nil {
			break
		} else {
			self.Tags[strings.Trim(C.GoString((*value).name), " ")] = strings.Trim(C.GoString((*value).value), " ")
		}
		C.free_exif_value(value)
	}

	C.free(unsafe.Pointer(values))

	return nil
}
