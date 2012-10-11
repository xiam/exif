/*
  Copyright (c) 2012 José Carlos Nieto, http://xiam.menteslibres.org/

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package exif

/*
#cgo LDFLAGS: -lexif

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libexif/exif-data.h>

#define MAXLEN 256

typedef struct exif_value {
	char *name;
	char *value;
	struct exif_value* prev;
} EXIF_VALUE;

typedef struct exif_stack {
	struct exif_value* head;
} EXIF_STACK;

EXIF_VALUE *new_exif_value(void);
void push_exif_value(EXIF_STACK*, EXIF_VALUE*);
void import_entry(ExifEntry*, void*);
void import_ifds(ExifContent*, void*);
EXIF_STACK* exif_dump(ExifData*);

void import_entry(ExifEntry* entry, void* user_data) {
	EXIF_VALUE* value;
	char exif_text[MAXLEN];

	value = new_exif_value();

	strncpy(value->name, exif_tag_get_title(entry->tag), MAXLEN);
	strncpy(value->value, exif_entry_get_value(entry, exif_text, MAXLEN), MAXLEN);

	push_exif_value(user_data, value);
}

void import_ifds(ExifContent* content, void* user_data) {
	exif_content_foreach_entry(content, import_entry, user_data);
}

EXIF_VALUE* new_exif_value() {
	EXIF_VALUE* n;
	n = (EXIF_VALUE*) malloc(sizeof(EXIF_VALUE));

	if (n == NULL) {
		return NULL;
	}

	n->name = (char *)malloc(sizeof(char)*MAXLEN);
	n->value = (char *)malloc(sizeof(char)*MAXLEN);

	n->name[0] 	= '\0';
	n->value[0] = '\0';
	n->prev  		= 0;
	return n;
}

void push_exif_value(EXIF_STACK* stack, EXIF_VALUE* n) {
	n->prev = stack->head;
	stack->head = n;
}

EXIF_VALUE* pop_exif_value(EXIF_STACK *stack) {
	EXIF_VALUE *n;
	if (stack->head == NULL) {
		return NULL;
	}
	n = stack->head;
	stack->head = n->prev;
	return n;
}

void free_exif_value(EXIF_VALUE* n) {
	free(n->name);
	free(n->value);
	free(n);
}

EXIF_STACK* exif_dump(ExifData* data) {
	EXIF_STACK* user_data;

	user_data = (EXIF_STACK*)malloc(sizeof(EXIF_STACK));
	user_data->head = NULL;

	exif_data_foreach_content(data, import_ifds, user_data);

	return user_data;
}

*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

type Data struct {
	ed   *C.ExifData
	Tags map[string]string
}

func New() *Data {
	self := &Data{}
	self.Tags = make(map[string]string)
	return self
}

func (self *Data) Open(file string) error {

	self.ed = C.exif_data_new_from_file(C.CString(file))

	if self.ed == nil {
		return fmt.Errorf("No EXIF data in file: %s\n", file)
	}

	values := C.exif_dump(self.ed)

	for true {
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
