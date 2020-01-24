// Package embedded embed files
//
// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package embedded

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type file struct {
	name       string
	size       int64
	modtime    int64
	local      string
	isDir      bool
	compressed bool
	mimeType   string
	tag        string
	data       []byte
	str        string
	subFiles   []os.FileInfo
}

type files struct {
	list  map[string]*file
	local bool
}

// New creates a new Files that loads embedded content.
func New() FileSystem {
	return &files{list: make(map[string]*file)}
}

func (fs *files) Open(name string) (file http.File, err error) {
	if f, ok := fs.list[name]; ok {
		if fs.local && len(f.local) > 0 {
			file, err = os.Open(f.local)
		} else {
			file = &reader{file: f, length: f.size}
		}
	} else {
		err = os.ErrNotExist
	}
	return
}

func (fs *files) Walk(walkFn WalkFunc) (err error) {
	for key, value := range fs.list {
		if err == nil {
			err = walkFn(key, value)
		}
	}

	return
}

func (fs *files) Copy(target string, mode os.FileMode) error {
	mode = mode & 0777
	dirmode := os.ModeDir | ((mode & 0444) >> 2) | mode
	return fs.Walk(func(path string, info FileInfo) error {
		targetPath := filepath.Join(target, path)
		if info.IsDir() {
			return os.MkdirAll(targetPath, dirmode)
		}

		file, err := fs.Open(path)
		if err == nil {
			defer file.Close()
			targetPathDir := filepath.Dir(targetPath)
			err = os.MkdirAll(targetPathDir, dirmode)
		}

		var out *os.File
		if err == nil {
			out, err = os.Create(targetPath)
		}

		if err == nil {
			_, err = io.Copy(out, file)

			out.Close()

			if err == nil {
				os.Chtimes(targetPath, info.ModTime(), info.ModTime())
				os.Chmod(targetPath, mode)
			}
		}

		return err
	})
}

// AddFile Adds a file to the file system
func (fs *files) AddFile(path string, name string, local string, size int64, modtime int64, mimeType string, tag string, compressed bool, data []byte, str string) {
	fs.list[path] = &file{
		name:       name,
		local:      local,
		size:       size,
		modtime:    modtime,
		mimeType:   mimeType,
		tag:        tag,
		compressed: compressed,
		data:       data,
		str:        str,
	}
}

// AddFolder Adds a folder to the file system
func (fs *files) AddFolder(path string, name string, local string, modtime int64) {
	fs.list[path] = &file{
		name:    name,
		local:   local,
		isDir:   true,
		modtime: modtime,
	}
}

// SetFiles Adds list of file to a folder
func (fs *files) SetFiles(path string, paths ...string) {
	subFiles := make([]os.FileInfo, len(paths))

	for i, e := range paths {
		subFiles[i] = fs.list[e]
	}

	fs.list[path].subFiles = subFiles
}

func (fs *files) UseLocal(value bool) {
	fs.local = value
}

// Name name of the file
func (f *file) Name() string {
	return f.name
}

// Size length in bytes for regular files
func (f *file) Size() int64 {
	return f.size
}

// Mode file mode bits
func (f *file) Mode() os.FileMode {
	if f.isDir {
		return os.ModePerm | os.ModeDir
	}
	return os.ModePerm
}

// ModTime file modification time
func (f *file) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

// IsDir abbreviation for Mode().IsDir()
func (f *file) IsDir() bool {
	return f.isDir
}

// Sys underlying data source (can return nil)
func (f *file) Sys() interface{} {
	return f
}

func (f *file) Compressed() bool {
	return f.compressed
}

func (f *file) Tag() string {
	return f.tag
}

func (f *file) MimeType() string {
	return f.mimeType
}

// String returns (uncompressed, if necessary) content of file as a string
func (f *file) String() (str string) {
	str = f.str
	if !f.IsDir() && f.compressed {
		var builder strings.Builder
		ungzip, _ := gzip.NewReader(bytes.NewReader(f.data))
		io.Copy(&builder, ungzip)
		ungzip.Close()
		str = builder.String()
	}
	return
}

// Bytes returns (uncompressed) content of file as a []byte
func (f *file) Bytes() (buf []byte) {
	if !f.IsDir() {
		if f.compressed {
			ungzip, _ := gzip.NewReader(bytes.NewReader(f.data))
			buf, _ = ioutil.ReadAll(ungzip)
			ungzip.Close()
		} else {
			buf = make([]byte, len(f.data))
			copy(buf, f.data)
		}
	}
	return
}

type reader struct {
	*file
	closed         bool
	position       int64
	seek           int64
	length         int64
	readCompressed bool
	decompressor   *gzip.Reader
	byteReader     *bytes.Reader
}

func (r *reader) Close() (err error) {
	if r.closed {
		err = os.ErrInvalid
	} else {
		r.closed = true
		if r.decompressor != nil {
			err = r.decompressor.Close()
			r.decompressor = nil
		}
		r.byteReader = nil
	}
	return
}

func (r *reader) Stat() (info os.FileInfo, err error) {
	if r.closed {
		err = os.ErrInvalid
	} else {
		info = r
	}
	return
}

func (r *reader) Readdir(count int) (list []os.FileInfo, err error) {
	if r.closed {
		err = os.ErrInvalid
	} else {
		if r.isDir {
			l := int64(len(r.subFiles))
			if r.position >= l && count > 0 {
				err = io.EOF
			} else {
				if count <= 0 || int64(count) > l-r.position {
					count = int(l - r.position)
				}
				list = r.subFiles[r.position : r.position+int64(count)]
				r.position += int64(count)
			}

		} else {
			err = errors.New("reader.Readdir: not valid on file")
		}
	}
	return
}

func (r *reader) Seek(offset int64, whence int) (n int64, err error) {
	if r.closed {
		err = os.ErrInvalid
	} else {
		if r.isDir {
			if offset == 0 && whence == io.SeekStart {
				r.position = 0
			} else {
				err = errors.New("reader.Seek: invalid Seek on directory")
			}
		} else {
			switch whence {
			case io.SeekStart:
				n = 0 + offset
			case io.SeekCurrent:
				n = int64(r.seek) + offset
			case io.SeekEnd:
				n = r.length + offset
			default:
				return 0, errors.New("reader.Seek: invalid whence")
			}

			if n < 0 {
				return 0, errors.New("reader.Seek: negative position")
			}

			r.seek = n
		}
	}
	return
}

func (r *reader) Read(b []byte) (n int, err error) {
	if r.closed {
		err = os.ErrInvalid
	} else {
		if r.isDir {
			err = errors.New("fileRead.Read: not valid on directory")
		} else {
			if r.compressed && !r.readCompressed {
				// Setup to decompress data
				if r.decompressor == nil {
					r.byteReader = bytes.NewReader(r.data)
					r.decompressor, err = gzip.NewReader(r.byteReader)
				}

				if err == nil {
					if r.position > r.seek {
						// Rewind to beginning.
						r.byteReader.Seek(io.SeekStart, 0)
						err = r.decompressor.Reset(r.byteReader)
						r.position = 0
					}

					if err == nil && r.position < r.seek {
						if _, err = io.CopyN(ioutil.Discard, r.decompressor, r.seek-r.position); err == nil {
							r.position = r.seek
						}
					}

					if err == nil {
						n, err = r.decompressor.Read(b)
						r.position += int64(n)
						r.seek = r.position
					}
				}
			} else {
				if r.seek >= r.length {
					return 0, io.EOF
				}
				n = copy(b, r.data[r.seek:])
				r.seek += int64(n)
			}
		}
	}
	return
}
