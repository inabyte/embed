package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/svg"
)

type file struct {
	name       string
	baseName   string
	data       []byte
	local      string
	Size       int64
	ModTime    int64
	mimeType   string
	tag        string
	Compressed bool
	Offset     int64

	fileinfo os.FileInfo
}

type dir struct {
	name     string
	baseName string
	local    string
	ModTime  int64
	files    map[string]bool
}

var (
	minifier *minify.M
	stringer builder
)

func (f *file) set() {
	stringer.add(f.name)
	stringer.add(f.baseName)
	stringer.add(f.local)
	stringer.add(f.mimeType)
	stringer.add(f.tag)
}

func (f *file) Data() []byte {
	return f.data
}

func (f *file) Slice() string {
	return fmt.Sprintf("%d:%d", f.Offset, f.Offset+int64(len(f.data)))
}

func (f *file) Name() string {
	return stringer.slice(f.name)
}

func (f *file) BaseName() string {
	return stringer.slice(f.baseName)
}

func (f *file) Local() string {
	return stringer.slice(f.local)
}
func (f *file) MimeType() string {
	return stringer.slice(f.mimeType)
}

func (f *file) Tag() string {
	return stringer.slice(f.tag)
}

func (f *file) setMimeType() {
	f.mimeType = mime.TypeByExtension(filepath.Ext(f.name))
	if f.mimeType == "" {
		// read a chunk to decide between utf-8 text and binary
		f.mimeType = http.DetectContentType(f.data)
	}
}

func (f *file) fill() {
	hash := sha1.Sum(f.data)
	f.tag = base64.RawURLEncoding.EncodeToString(hash[:]) + "-gz"
	f.Size = int64(len(f.data))
}

func (f *file) compress() error {
	var buf = &bytes.Buffer{}

	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)

	if err == nil {
		_, err = gw.Write(f.data)
	}

	if err == nil {
		err = gw.Close()
	}

	if err == nil {
		if buf.Len() < len(f.data) {
			f.data = buf.Bytes()
			f.Compressed = true
		}
	}

	return err
}

func (f *file) minify() {
	if b, err := minifier.Bytes(f.mimeType, f.data); err == nil {
		f.data = b
	}
}

func (d *dir) set() {
	stringer.add(d.name)
	stringer.add(d.baseName)
	stringer.add(d.local)

	for k := range d.files {
		stringer.add(k)
	}
}

func (d *dir) Name() string {
	return stringer.slice(d.name)
}

func (d *dir) BaseName() string {
	return stringer.slice(d.baseName)
}

func (d *dir) Local() string {
	return stringer.slice(d.local)
}

func (d *dir) Files() []string {
	res := make([]string, len(d.files))

	i := 0
	for k := range d.files {
		res[i] = k
		i++
	}

	sort.Strings(res)

	for i, entry := range res {
		res[i] = stringer.slice(entry)
	}

	return res
}

func init() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)
	minifier.AddFunc("application/javascript", js.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.Add("text/html", &html.Minifier{
		KeepConditionalComments: true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
	})
	minifier.Add("text/html; charset=utf-8", &html.Minifier{
		KeepConditionalComments: true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
	})
}
