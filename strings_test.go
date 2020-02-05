package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"reflect"
	"testing"
)

const (
	fullString  = "some test string"
	startString = "some test"
	endString   = "string"
)

func TestStringsBuilder(t *testing.T) {

	var builder builder

	builder.add(startString)
	builder.add(endString)
	builder.add(fullString)

	builder.write(&mocWriter{})

	buf := []byte(fullString)
	if b := []byte(builder.str); !reflect.DeepEqual(b, buf) {
		t.Errorf("Did not get expected buffer got (%v) expect (%v)", b, buf)
	}

	for _, test := range []struct {
		name   string
		str    string
		expect string
	}{
		{"full", fullString, "/* some test string */ str[0:16]"},
		{"start", startString, "/* some test */ str[0:9]"},
		{"end", endString, "/* string */ str[10:16]"},
		{"not in", "not", `"not"`},
		{"blank", "", `""`},
	} {
		t.Run(test.name, func(t *testing.T) {

			s := builder.slice(test.str)

			if s != test.expect {
				t.Errorf("Did not get expected for %s got (%s) expected(%s)", test.str, s, test.expect)
			}
		})
	}

}

type mocWriter struct {
	bytes []byte
}

func (m *mocWriter) offset() int {
	return 0
}

func (m *mocWriter) Write(p []byte) (n int, err error) {
	m.bytes = append(m.bytes, p...)
	return len(p), nil
}

func (m *mocWriter) Close() error {
	return nil
}
