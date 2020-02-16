package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
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
	base, err := ioutil.TempDir("", "file-test")

	w := mocWriter{}

	f := file{
		name:     "/scripts/test.html",
		baseName: "test.html",
		path:     "test.html",
		local:    "test.html",
		ModTime:  1579282495,
	}

	if err == nil {
		var curdir string

		defer os.RemoveAll(base)
		curdir, err = os.Getwd()
		if err == nil {
			defer os.Chdir(curdir)
			err = os.Chdir(base)
		}
	}

	if err == nil {
		err = ioutil.WriteFile(f.path, []byte(rawContents), os.ModePerm)
	}

	if err == nil {
		err = f.write(&w)
	}

	if err == nil {
		err = stringer.write(&mocWriter{})
	}

	t.Run("Contents", func(t *testing.T) {
		if !reflect.DeepEqual(w.bytes, compress(minifyContents)) {
			t.Errorf("Did not get expected for %s got (%s) expected(%s)", "Contents", unCompress(w.bytes), minifyContents)
		}
	})

	if err == nil {
		fileTests{
			{"Name", "/* /scripts/test.html */ str[54:72]"},
			{"BaseName", "/* test.html */ str[63:72]"},
			{"Local", "/* test.html */ str[63:72]"},
			{"MimeType", "/* text/html; charset=utf-8 */ str[30:54]"},
			{"Tag", "/* xwI1ooNerSnDfL_w9IZZoVz_A4Y-gz */ str[0:30]"},
			{"Slice", "0:98"},
		}.run(t, &f)
	} else {
		t.Errorf("Error setting up test %v", err)
	}
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

	w := mocWriter{}
	err := stringer.write(&w)

	if err == nil {
		fileTests{
			{"Name", "/* /scripts */ str[54:62]"},
			{"BaseName", "/* scripts */ str[55:62]"},
			{"Local", "/* embed/scripts */ str[91:104]"},
			{"Files", []string{"/* /scripts/index.html */ str[54:73]"}},
		}.run(t, &d)
	} else {
		t.Errorf("Error setting up test %v", err)
	}
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
