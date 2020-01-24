package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"bytes"
	"compress/gzip"
	"io"
	"reflect"
	"strings"
	"testing"
)

const (
	rawContents = `<html>
<title>Test</title>

<body>
Some text to make this compressable Some text to make this compressable Some text to make this compressable Some text to make this compressable Some text to make this compressable
</body>
</html>
`

	minifyContents = `<html><title>Test</title><body>Some text to make this compressable Some text to make this compressable Some text to make this compressable Some text to make this compressable Some text to make this compressable</body></html>`
)

func TestFile(t *testing.T) {
	f := file{
		name:     "/scripts/index",
		baseName: "index",
		local:    "embed/scripts/index",
		ModTime:  1579282495,
		data:     []byte(rawContents),
	}

	f.fill()
	f.setMimeType()
	f.minify()
	f.compress()

	f.set()

	stringer.process()

	fileTests{
		{"Name", "/* /scripts/index */ str[59:73]"},
		{"BaseName", "/* index */ str[68:73]"},
		{"Local", "/* embed/scripts/index */ str[54:73]"},
		{"MimeType", "/* text/html; charset=utf-8 */ str[30:54]"},
		{"Tag", "/* J1A-DFNtnAw81oBuVFV4VjPNReo-gz */ str[0:30]"},
		{"Slice", "0:98"},
		{"Data", compress(minifyContents)},
	}.run(t, &f)
}

func TestDir(t *testing.T) {
	d := dir{
		name:     "/scripts",
		baseName: "scripts",
		local:    "embed/scripts",
		ModTime:  1579282495,
		files:    map[string]bool{"/scripts/index.html": true},
	}

	d.set()

	stringer.process()

	fileTests{
		{"Name", "/* /scripts */ str[59:67]"},
		{"BaseName", "/* scripts */ str[60:67]"},
		{"Local", "/* embed/scripts */ str[54:67]"},
		{"Files", []string{"/* /scripts/index.html */ str[73:92]"}},
	}.run(t, &d)
}

type fileTest struct {
	name   string
	expect interface{}
}

type fileTests []fileTest

func (list fileTests) run(t *testing.T, obj interface{}) {
	for _, test := range list {
		{
			t.Run(test.name, func(t *testing.T) {
				value := reflect.ValueOf(obj)
				m := value.MethodByName(test.name)

				if m.IsValid() {
					out := m.Call(nil)

					if len(out) != 1 {
						t.Errorf("Wrong number of return values(%d) for method %s", len(out), test.name)
					} else {
						data := out[0].Interface()

						if !reflect.DeepEqual(data, test.expect) {
							if d, ok := data.([]byte); ok {
								t.Errorf("Did not get expected for %s got (%s) expected(%s)", test.name, unCompress(d), unCompress(test.expect.([]byte)))
							} else {
								t.Errorf("Did not get expected for %s got (%v) expected(%v)", test.name, data, test.expect)
							}
						}
					}
				} else {
					t.Errorf("Function %s not found", test.name)
				}
			})
		}
	}
}

func unCompress(data []byte) (s string) {
	var buf strings.Builder

	ungzip, err := gzip.NewReader(bytes.NewReader(data))

	if err == nil {
		_, err = io.Copy(&buf, ungzip)
	}
	if err == nil {
		err = ungzip.Close()
	}
	if err == nil {
		s = buf.String()
	}
	return
}

func compress(data string) (b []byte) {
	var buf = &bytes.Buffer{}

	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)

	if err == nil {
		_, err = io.WriteString(gw, data)
	}
	if err == nil {
		err = gw.Close()
	}
	if err == nil {
		b = buf.Bytes()
	}

	return
}
