// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(t *testing.T) {

	log = mocLogger(0)
	myExit = func(int) {}

	for _, test := range []struct {
		name string
		args []string
	}{
		{"no args", []string{"go-embed"}},
		{"bad args", []string{"go-embed", "-h"}},
		{"file", []string{"go-embed", "files"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args
			f = flag.NewFlagSet("go-embed", flag.PanicOnError)
			f.SetOutput(ioutil.Discard)

			defer func() {
				if r := recover(); r != nil {
				}
				return
			}()

			main()
		})
	}
}

type mocLogger int

func (mocLogger) Print(v ...interface{}) {}
