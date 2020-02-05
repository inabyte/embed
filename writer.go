package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	lineSize = 16
)

type writer interface {
	io.WriteCloser
	offset() int
}

type fileWriter struct {
	name        string
	isGo        bool
	buf         [lineSize]byte
	strBuf      [lineSize * 4]byte
	index       int
	f           *os.File
	dataOffset  int
	writeOffset int
}

func createFile(path string, name string, extension string) (file *os.File, err error) {

	basePath := filepath.Dir(path)

	if len(basePath) > 0 && basePath != "." && basePath != string(filepath.Separator) {
		err = os.MkdirAll(basePath, os.ModePerm)
	}

	if err == nil {
		file, err = os.Create(fmt.Sprintf("%s%s%s", path, name, extension))
	}

	return
}

func createWriter(isGo bool, pkg string, name string, path string, tags ...string) (w writer, err error) {
	var file *os.File

	ext := ".s"
	if isGo {
		ext = ".go"
	}

	buildTags := strings.TrimSpace(strings.Join(tags, " "))
	file, err = createFile(path, "_data", ext)

	if err == nil {
		_, err = file.WriteString(header)
	}

	if err == nil && len(buildTags) > 0 {
		_, err = file.WriteString("\n// +build " + buildTags)
	}

	if err == nil {
		if isGo {
			_, err = fmt.Fprintf(file, "\n\npackage %s\n\nconst (\n\t%sData = ", pkg, name)
		} else {
			_, err = file.WriteString("\n\n#include \"textflag.h\"\n\n")
		}
	}

	if err == nil {
		w = &fileWriter{name: name, isGo: isGo, f: file}
	}

	return
}

func (w *fileWriter) offset() int {
	return w.dataOffset
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	for _, b := range p {
		if err == nil {
			w.buf[w.index] = b
			w.index++
			w.dataOffset++
			if w.index == len(w.buf) {
				err = w.flush()
			}
		}
	}

	return
}

func (w *fileWriter) flush() (err error) {

	if w.index > 0 {
		var sbuf = w.strBuf[0:0]

		for i := 0; i < w.index; i++ {
			sbuf = append(sbuf, []byte("\\x")...)
			if w.buf[i] < 0x10 {
				sbuf = append(sbuf, '0')
			}
			sbuf = strconv.AppendUint(sbuf, uint64(w.buf[i]), 16)
		}

		if w.isGo {
			if w.writeOffset != 0 {
				_, err = fmt.Fprint(w.f, " +\n\t\t")
			}
			_, err = fmt.Fprintf(w.f, "\"%s\"", sbuf)
		} else {
			_, err = fmt.Fprintf(w.f, "DATA ·%sData+%d(SB)/%d,$\"%s\"\n", w.name, w.writeOffset, w.index, sbuf)
		}
		w.writeOffset += w.index
		w.index = 0
	}

	return
}

func (w *fileWriter) footer() (err error) {
	err = w.flush()

	if err == nil {
		if w.isGo {
			_, err = fmt.Fprintln(w.f, "\n)")
		} else {
			_, err = fmt.Fprintf(w.f, "GLOBL ·%sData(SB),(NOPTR+RODATA),$%d\n", w.name, w.writeOffset)
		}
	}
	return
}

func (w *fileWriter) Close() error {

	if w == nil || w.f == nil {
		return os.ErrInvalid
	}

	w.footer()

	f := w.f
	w.f = nil

	return f.Close()
}
