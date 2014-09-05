/*
  Copyright (c) 2012-2013 José Carlos Nieto, https://menteslibres.net/xiam

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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <libexif/exif-data.h>

#include "_cgo/types.h"

#define EXIF_VALUE_MAXLEN 256

void import_entry(ExifEntry*, void*);
void import_ifds(ExifContent*, void*);
exif_value_t *new_exif_value(void);
void push_exif_value(exif_stack_t*, exif_value_t*);
exif_value_t* pop_exif_value(exif_stack_t *);
void free_exif_value(exif_value_t* n);
exif_stack_t* exif_dump(ExifData *);

void import_entry(ExifEntry* entry, void* user_data) {
  exif_value_t* value;
  char exif_text[EXIF_VALUE_MAXLEN];

  value = new_exif_value();

  ExifIfd ifd = exif_entry_get_ifd(entry);

  strncpy(value->name, exif_tag_get_title_in_ifd(entry->tag, ifd), EXIF_VALUE_MAXLEN);

  strncpy(value->value, exif_entry_get_value(entry, exif_text, EXIF_VALUE_MAXLEN), EXIF_VALUE_MAXLEN);

  push_exif_value(user_data, value);
}

void import_ifds(ExifContent* content, void* user_data) {
  exif_content_foreach_entry(content, import_entry, user_data);
}

exif_value_t* new_exif_value() {
  exif_value_t* n;
  n = (exif_value_t*) malloc(sizeof(exif_value_t));

  if (n == NULL) {
    return NULL;
  }

  n->name = (char *)malloc(sizeof(char)*EXIF_VALUE_MAXLEN);
  n->value = (char *)malloc(sizeof(char)*EXIF_VALUE_MAXLEN);

  n->name[0]  = '\0';
  n->value[0] = '\0';
  n->prev     = 0;
  return n;
}

void push_exif_value(exif_stack_t* stack, exif_value_t* n) {
  n->prev = stack->head;
  stack->head = n;
}

exif_value_t* pop_exif_value(exif_stack_t *stack) {
  exif_value_t *n;
  if (stack->head == NULL) {
    return NULL;
  }
  n = stack->head;
  stack->head = n->prev;
  return n;
}

void free_exif_value(exif_value_t* n) {
  free(n->name);
  free(n->value);
  free(n);
}

exif_stack_t* exif_dump(ExifData* data) {
  exif_stack_t* user_data;

  user_data = (exif_stack_t*)malloc(sizeof(exif_stack_t));
  user_data->head = NULL;

  exif_data_foreach_content(data, import_ifds, user_data);

  return user_data;
}
