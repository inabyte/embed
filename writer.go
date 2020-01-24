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
	buf         [lineSize]byte
	strBuf      [lineSize * 4]byte
	index       int
	f           *os.File
	dataOffset  int64
	writeOffset int64
}

type goWriter struct {
	buf         [lineSize]byte
	strBuf      [lineSize * 4]byte
	index       int
	f           *os.File
	dataOffset  int64
	writeOffset int64
}

type asmFile struct {
	name      string
	buildTags string
	code      string
}

type asmFiles []asmFile

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

func createWriter(path string, tags ...string) (w writer, err error) {
	var file *os.File

	if file, err = createWriteHeaderInclude(path, "", ".s", tags...); err == nil {
		w = &asmWriter{f: file}
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
		_, err = fmt.Fprintf(w.f, "GLOBL ·data(SB),RODATA,$%d\n", w.writeOffset)
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

		_, err = fmt.Fprintf(w.f, "DATA ·data+%d(SB)/%d,$\"%s\"\n", w.writeOffset, w.index, sbuf)
		w.writeOffset += int64(w.index)
		w.index = 0
	}

	return
}

func createGoWriter(file *os.File) (w writer, err error) {
	w = &goWriter{f: file}

	_, err = fmt.Fprint(file, `
func dataBytes() []byte {
	str := dataString()
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

func dataString() string {
	return data
}

const (
	data = `)

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

func (list asmFiles) output(path string, buildTags string) (err error) {
	var (
		file *os.File
	)

	for _, entry := range list {
		if err == nil {
			file, err = createWriteHeaderInclude(path, entry.name, ".s", buildTags, entry.buildTags)
		}

		if err == nil {
			_, err = file.WriteString(entry.code)
		}

		if file != nil {
			file.Close()
			file = nil
		}
	}

	return
}

var assmemblerFiles = asmFiles{
	{"_386", "", `
TEXT ·file_bytes(SB),NOSPLIT,$0-4
	LEAL	·data(SB), AX
	MOVL	AX, ret+4(FP)
	MOVL	len+0(FP), AX
	MOVL	AX, ret+8(FP)
	MOVL	AX, ret+12(FP)
	RET

TEXT ·file_string(SB),NOSPLIT,$0-4
	LEAL	·data(SB), AX
	MOVL	AX, ret+4(FP)
	MOVL	len+0(FP), AX
	MOVL	AX, ret+8(FP)
	RET
`},
	{"_amd64", "", `
TEXT ·file_bytes(SB),NOSPLIT,$0-4
	LEAQ	·data(SB), AX
	MOVQ	AX, ret+8(FP)
	MOVL	len+0(FP), AX
	MOVLQSX	AX, AX
	MOVQ	AX, ret+16(FP)
	MOVQ	AX, ret+24(FP)
	RET

TEXT ·file_string(SB),NOSPLIT,$0-4
	LEAQ	·data(SB), AX
	MOVQ	AX, ret+8(FP)
	MOVL	len+0(FP), AX
	MOVLQSX	AX, AX
	MOVQ	AX, ret+16(FP)
	RET
`},
	{"_arm", "", `
TEXT ·file_bytes(SB),NOSPLIT,$0-4
	MOVW	$·data(SB), R0
	MOVW	R0, ret+4(FP)
	MOVW	len+0(FP), R0
	MOVW	R0, ret+8(FP)
	MOVW	R0, ret+12(FP)
	RET

TEXT ·file_string(SB),NOSPLIT,$0-4
	MOVW	$·data(SB), R0
	MOVW	R0, ret+4(FP)
	MOVW	len+0(FP), R0
	MOVW	R0, ret+8(FP)
	RET
`},
	{"_arm64", "", `
TEXT ·file_bytes(SB),NOSPLIT,$0-8
	MOVD	$·data(SB), R0
	MOVD	R0, ret+8(FP)
	MOVW	len+0(FP), R0
	MOVD	R0, ret+16(FP)
	MOVD	R0, ret+24(FP)
	RET

TEXT ·file_string(SB),NOSPLIT,$0-8
	MOVD	$·data(SB), R0
	MOVD	R0, ret+8(FP)
	MOVD	len+0(FP), R0
	MOVD	R0, ret+16(FP)
	RET
`},
	{"_mipsx", "mips mipsle", `
TEXT ·file_bytes(SB),NOSPLIT,$0-4
	MOVW	$·data(SB), R1
	MOVW	R1, ret+4(FP)
	MOVW	len+0(FP), R1
	MOVW	R1, ret+8(FP)
	MOVW	R1, ret+12(FP)
	JMP	(R31)

TEXT ·file_string(SB),NOSPLIT,$0-4
	MOVW	$·data(SB), R1
	MOVW	R1, ret+4(FP)
	MOVW	len+0(FP), R1
	MOVW	R1, ret+8(FP)
	JMP	(R31)
`},
	{"_mips64x", "mips64 mips64le", `
TEXT ·file_bytes(SB),NOSPLIT,$0-8
	MOVV	$·data(SB), R1
	MOVV	R1, ret+8(FP)
	MOVV	len+0(FP), R1
	MOVV	R1, ret+16(FP)
	MOVV	R1, ret+24(FP)
	JMP	(R31)

TEXT ·file_string(SB),NOSPLIT,$0-8
	MOVV	$·data(SB), R1
	MOVV	R1, ret+8(FP)
	MOVV	len+0(FP), R1
	MOVV	R1, ret+16(FP)
	JMP	(R31)
`},
	{"_ppc64", "", `
TEXT ·file_bytes(SB),NOSPLIT,$0-8
	MOVD	$·data(SB), R3
	MOVD	R3, ret+8(FP)
	MOVD	len+0(FP), R3
	MOVD	R3, ret+16(FP)
	MOVD	R3, ret+24(FP)
	RET

TEXT ·file_string(SB),NOSPLIT,$0-8
	MOVD	$·data(SB), R3
	MOVD	R3, ret+8(FP)
	MOVD	len+0(FP), R3
	MOVD	R3, ret+16(FP)
	RET
`},
	{"_s390x", "", `
TEXT ·file_bytes(SB),NOSPLIT|NOFRAME,$0-8
	MOVD	$·data(SB), R0
	MOVW	len+0(FP), R1
	MOVD	R1, R2
	STMG	R0, R2, ret+8(FP)
	JMP	R14

TEXT ·file_string(SB),NOSPLIT|NOFRAME,$0-8
	MOVD	$·data(SB), R0
	MOVW	len+0(FP), R1
	STMG	R0, R1, ret+8(FP)
	JMP	R14
`},
}
