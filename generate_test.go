package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {

	base, err := createFs()

	if len(base) > 0 {
		defer os.RemoveAll(base)
	}

	if err != nil {
		t.Fatalf("unable to cerate fs %v", err)
	}

	config := New()
	config.Output = filepath.Join(base, "assets", "files")
	fpath := filepath.Join(base, "www")
	l := len(fpath) - 4
	fpath = fpath[:l] + PrefixMarker + fpath[l:]
	config.Files = []string{fpath, filepath.Join(base, "single", "settings.html")}
	config.Ignore = `\.go$`
	config.Include = `\.html$`

	for _, v := range []struct {
		name   string
		hasErr bool
		doFunc func() func()
	}{
		{
			name:   "No files",
			hasErr: true,
			doFunc: func() func() {
				files := config.Files
				config.Files = nil
				return func() { config.Files = files }
			},
		},
		{
			name:   "Bad Int",
			hasErr: true,
			doFunc: func() func() {
				config.ModifyTime = "bad int"
				return func() { config.ModifyTime = "" }
			},
		},
		{
			name:   "Duplicates",
			hasErr: true,
			doFunc: func() func() {
				files := config.Files
				config.Files = append(config.Files, filepath.Join(base, "repeat", "settings.html"))
				return func() { config.Files = files }
			},
		},
		{
			name: "Success",
			doFunc: func() func() {
				return func() {}
			},
		},
		{
			name: "Fixed Time",
			doFunc: func() func() {
				config.ModifyTime = fmt.Sprintf("%d", time.Now().Unix())
				return func() { config.ModifyTime = "" }
			},
		},
		{
			name: "Binary",
			doFunc: func() func() {
				config.Binary = true
				return func() { config.Binary = false }
			},
		},
		{
			name: "Binary Relative",
			doFunc: func() func() {
				oldOutput := config.Output
				config.Output = filepath.Join("assets", "files")
				old, _ := os.Getwd()
				os.Chdir(base)
				config.Binary = true
				return func() { config.Output = oldOutput; config.Binary = false; os.Chdir(old) }
			},
		},
		{
			name: "Go",
			doFunc: func() func() {
				config.Go = true
				return func() { config.Go = false }
			},
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			post := v.doFunc()
			err := config.Generate()
			post()
			if err == nil {
				if v.hasErr {
					t.Errorf("Generate did not return an error")
				}
			} else {
				if !v.hasErr {
					t.Errorf("Generate returned unexpected error %v", err)
				}
			}
		})
	}
}

func createFs() (string, error) {
	base, err := ioutil.TempDir("", "generate-test")

	for _, v := range []struct {
		path    string
		content []byte
	}{
		{"www/index.html", []byte("<html></html>")},
		{"www/scripts/init.js", []byte("var data")},
		{"www/code/process.go", []byte("package process")},
		{"single/settings.html", []byte("<html></html>")},
		{"repeat/settings.html", []byte("<html></html>")},
	} {
		if err == nil {
			path := filepath.Join(base, v.path)

			if err == nil {
				err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			}
			if err == nil {
				err = ioutil.WriteFile(path, v.content, os.ModePerm)
			}
		}
	}

	return base, err
}
