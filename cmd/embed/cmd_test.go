// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestBinary(t *testing.T) {
	base, err := testContents{
		{"www/index.html", []byte("<html></html>")},
		{"www/scripts/init.js", []byte("var data")},
	}.Create()

	if len(base) > 0 {
		defer os.RemoveAll(base)
	}

	if err == nil {
		t.Run("binary", func(t *testing.T) {
			os.Chdir(base)
			os.Args = []string{"go-embed", "-binary", "-o=" + filepath.Join("bin", "srv"), "www"}
			f = flag.NewFlagSet("go-embed", flag.PanicOnError)
			f.SetOutput(ioutil.Discard)

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Binary failed %v", r)
				}
				return
			}()

			main()

			if f, err := os.Open(filepath.Join(base, "bin", "srv")); err == nil {
				f.Close()
			} else {
				t.Errorf("Counld not find generated target %v", err)
			}

		})
	}
}

type mocLogger int

func (mocLogger) Print(v ...interface{}) {}

type testContent struct {
	path    string
	content []byte
}

type testContents []testContent

func (list testContents) Create() (string, error) {

	base, err := ioutil.TempDir("", "generate-test")

	for _, v := range list {
		path := filepath.Join(base, v.path)

		if err == nil {
			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		}
		if err == nil {
			err = ioutil.WriteFile(path, v.content, os.ModePerm)
		}
	}

	if err != nil {
		return base, err
	}

	return base, nil
}
