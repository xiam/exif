# Go (golang) bindings for libexif

Provides basic support for reading EXIF tags on files using [libexif][1] and
CGO.

## How to install

Get libexif for your OS:

```
# OSX
brew install libexif

# Debian
apt-get install -y libexif-dev

# RedHat
dnf install -y libexif-devel
```

Then grab the exif package with `go get`:

```
go get github.com/xiam/exif
```

## Usage

Install the package with `go get` and use `import` to include it in your
project:

```
import "github.com/xiam/exif"
```

This is an example on how to read EXIF data from a file:

```go
package main

import (
  "fmt"
  "github.com/xiam/exif"
)

func main() {
    data, err := exif.Read("_examples/resources/test.jpg")
    ...
    for key, val := range data.Tags {
        fmt.Printf("%s = %s\n", key, val)
    }
}
```

If you just have the image available as an io.Reader, you can parse the EXIF
header like this:

```
reader := exif.New()

_, err = io.Copy(reader, data)

// exif.FoundExifInData is a signal that the EXIF parser has all it needs,
// it doesn't need to be given the whole image.
if err != nil && err != exif.FoundExifInData {
  t.Fatalf("Error loading bytes: %s", err.Error())
}

err := reader.Parse()

if err != nil {
  t.Fatalf("Error parsing EXIF: %s", err.Error())
}

for key, val := range reader.Tags {
  fmt.Printf("%s: %s\n", key, val)
}
```

## License

This is Open Source released under the terms of the MIT License:

> Copyright (c) 2012-2015 José Carlos Nieto, https://menteslibres.net/xiam
>
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[1]: http://libexif.sourceforge.net/
