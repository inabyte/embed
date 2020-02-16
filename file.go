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
	"io/ioutil"
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
	path       string
	local      string
	Size       int
	ModTime    int64
	mimeType   string
	tag        string
	dataSize   int
	Compressed bool
	offset     int

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

func (f *file) Slice() string {
	return fmt.Sprintf("%d:%d", f.offset, f.offset+f.dataSize)
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

func (f *file) write(w writer) error {
	var (
		buf bytes.Buffer
		gw  *gzip.Writer
	)

	f.offset = w.offset()
	b, err := ioutil.ReadFile(f.path)

	if err == nil {
		// Determine mimetype
		f.mimeType = mime.TypeByExtension(filepath.Ext(f.name))
		if f.mimeType == "" {
			// read a chunk to decide between utf-8 text and binary
			f.mimeType = http.DetectContentType(b)
		}

		// Minify the data
		if m, e := minifier.Bytes(f.mimeType, b); e == nil {
			b = m
		}

		// Create eTag
		hash := sha1.Sum(b)
		f.tag = base64.RawURLEncoding.EncodeToString(hash[:]) + "-gz"

		f.Size = len(b)
		f.dataSize = f.Size

		gw, err = gzip.NewWriterLevel(&buf, gzip.BestCompression)
	}

	if err == nil {
		_, err = gw.Write(b)
	}

	if err == nil {
		err = gw.Close()
	}

	if err == nil {
		if buf.Len() < f.Size {
			b = buf.Bytes()
			f.dataSize = len(b)
			f.Compressed = true
		}
	}

	if err == nil {
		f.dataSize, err = w.Write(b)
	}

	f.set()

	return err
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
