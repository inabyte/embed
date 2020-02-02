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
	offset() int64
	footer() error
}

type asmWriter struct {
	name        string
	buf         [lineSize]byte
	strBuf      [lineSize * 4]byte
	index       int
	f           *os.File
	dataOffset  int64
	writeOffset int64
}

type goWriter struct {
	name        string
	buf         [lineSize]byte
	strBuf      [lineSize * 4]byte
	index       int
	f           *os.File
	dataOffset  int64
	writeOffset int64
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

func createWriteHeader(path string, name string, extension string, tags ...string) (file *os.File, err error) {

	file, err = createFile(path, name, extension)

	if err == nil {
		_, err = file.WriteString(header)

		if err == nil {
			buildTags := strings.TrimSpace(strings.Join(tags, " "))
			if len(buildTags) > 0 {
				_, err = file.WriteString("\n// +build ")
				if err == nil {
					_, err = file.WriteString(buildTags)
				}
			}
		}
	}

	return
}

func createWriteHeaderInclude(path string, name string, extension string, tags ...string) (file *os.File, err error) {

	file, err = createWriteHeader(path, name, extension, tags...)

	if err == nil {
		if err == nil {
			_, err = file.WriteString("\n\n#include \"textflag.h\"\n\n")
		}
	}

	return
}

func createWriter(name string, path string, tags ...string) (w writer, err error) {
	var file *os.File

	if file, err = createWriteHeaderInclude(path, "", ".s", tags...); err == nil {
		w = &asmWriter{name: name, f: file}
	}

	return
}

func (w *asmWriter) offset() int64 {
	return w.dataOffset
}

func (w *asmWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	for _, b := range p {
		w.buf[w.index] = b
		w.index++
		w.dataOffset++
		if w.index == len(w.buf) {
			err = w.flush()
		}
	}

	return
}

func (w *asmWriter) Close() error {
	if w == nil || w.f == nil {
		return os.ErrInvalid
	}

	f := w.f
	w.f = nil

	return f.Close()
}

func (w *asmWriter) footer() (err error) {
	err = w.flush()

	if err == nil {
		_, err = fmt.Fprintf(w.f, "GLOBL ·%sData(SB),(NOPTR+RODATA),$%d\n", w.name, w.writeOffset)
	}
	return
}

func (w *asmWriter) flush() (err error) {

	if w.index > 0 {
		var sbuf = w.strBuf[0:0]

		for i := 0; i < w.index; i++ {
			sbuf = append(sbuf, []byte("\\x")...)
			if w.buf[i] < 0x10 {
				sbuf = append(sbuf, '0')
			}
			sbuf = strconv.AppendUint(sbuf, uint64(w.buf[i]), 16)
		}

		_, err = fmt.Fprintf(w.f, "DATA ·%sData+%d(SB)/%d,$\"%s\"\n", w.name, w.writeOffset, w.index, sbuf)
		w.writeOffset += int64(w.index)
		w.index = 0
	}

	return
}

func createGoWriter(name string, file *os.File) (w writer, err error) {
	w = &goWriter{name: name, f: file}

	_, err = fmt.Fprintf(file, `
func %sBytes() []byte {
	str := %sData
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

const (
	%sData = `, name, name, name)

	return
}

func (w *goWriter) offset() int64 {
	return w.dataOffset
}

func (w *goWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	for _, b := range p {
		w.buf[w.index] = b
		w.index++
		w.dataOffset++
		if w.index == len(w.buf) {
			err = w.flush()
		}
	}

	return
}

func (w *goWriter) Close() error {

	if w == nil || w.f == nil {
		return os.ErrInvalid
	}

	f := w.f
	w.f = nil

	return f.Close()
}

func (w *goWriter) footer() (err error) {
	err = w.flush()

	if err == nil {
		_, err = fmt.Fprintln(w.f, "\n)")
	}

	return
}

func (w *goWriter) flush() (err error) {

	if w.index > 0 {
		var sbuf = w.strBuf[0:0]

		for i := 0; i < w.index; i++ {
			sbuf = append(sbuf, []byte("\\x")...)
			if w.buf[i] < 0x10 {
				sbuf = append(sbuf, '0')
			}
			sbuf = strconv.AppendUint(sbuf, uint64(w.buf[i]), 16)
		}

		if w.writeOffset != 0 {
			_, err = fmt.Fprint(w.f, " +\n\t\t")
		}
		_, err = fmt.Fprintf(w.f, "\"%s\"", sbuf)
		w.writeOffset += int64(w.index)
		w.index = 0
	}

	return
}
