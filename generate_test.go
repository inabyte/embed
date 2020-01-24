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

var mocFs = testContents{
	{
		"www/index.html", []byte("<html></html>"),
	},
	{
		"www/scripts/init.js", []byte("var data"),
	},
	{
		"www/code/process.go", []byte("package process"),
	},
	{
		"single/settings.html", []byte("<html></html>"),
	},
	{
		"repeat/settings.html", []byte("<html></html>"),
	},
}

func TestGenerate(t *testing.T) {

	base, err := mocFs.Create()

	if len(base) > 0 {
		defer os.RemoveAll(base)
	}

	if err != nil {
		t.Fatalf("unable to cerate fs %v", err)
	}

	cfg := New()

	if err = cfg.Generate(); err == nil {
		t.Errorf("Generate with no files did not err")
	}

	cfg.Output = filepath.Join(base, "assets/files")
	fpath := filepath.Join(base, "www")
	l := len(fpath) - 4
	fpath = fpath[:l] + prefixMarker + fpath[l:]
	cfg.Files = []string{fpath, filepath.Join(base, "single", "settings.html")}
	cfg.Local = true
	cfg.ModifyTime = fmt.Sprintf("%d", time.Now().Unix())
	cfg.Ignore = `\.go$`
	cfg.Include = `\.html$`

	if err = cfg.Generate(); err != nil {
		t.Errorf("Generate returned error %v", err)
	}

	cfg.ModifyTime = "bad int"
	if err = cfg.Generate(); err == nil {
		t.Errorf("Generate bad modify time did no return an error")
	}

	cfg.ModifyTime = ""

	cfg.Files = append(cfg.Files, filepath.Join(base, "repeat", "settings.html"))
	if err = cfg.Generate(); err == nil {
		t.Errorf("Generate dup output did not return error")
	}

}

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
