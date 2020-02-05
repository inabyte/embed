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
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unsafe"
)

// FileSystem defines the FileSystem interface and builder
type FileSystem interface {
	http.FileSystem

	// Walk walks the file tree rooted at root, calling walkFn for each file or
	// directory in the tree, including root. All errors that arise visiting files
	// and directories are filtered by walkFn. The files are walked in lexical
	// order.
	Walk(root string, walkFn WalkFunc) error

	// Copy all files to target directory
	Copy(target string, mode os.FileMode) error

	// AddFile add a file to embedded filesystem
	AddFile(path string, name string, local string, size int64, modtime int64, mimeType string, tag string, compressed bool, data []byte, str string)
	// AddFolder add a file to embedded filesystem
	AddFolder(path string, name string, local string, modtime int64, paths ...string)

	// WriteFile writes data to a file named by filename.
	// If the file does not exist, WriteFile creates it with permissions perm;
	// otherwise WriteFile truncates it before writing.
	WriteFile(filename string, data []byte, perm os.FileMode) error

	// UseLocal use on disk copy instead of embedded data (for development)
	UseLocal(bool)
}

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = filepath.SkipDir

// WalkFunc is the type of the function called for each file or directory
// visited by Walk. The path argument contains the argument to Walk as a
// prefix; that is, if Walk is called with "dir", which is a directory
// containing the file "a", the walk function will be called with argument
// "dir/a". The info argument is the os.FileInfo for the named path.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how
// to handle that error (and Walk will not descend into that directory). In the
// case of an error, the info argument will be nil. If an error is returned,
// processing stops. The sole exception is when the function returns the special
// value SkipDir. If the function returns SkipDir when invoked on a directory,
// Walk skips the directory's contents entirely. If the function returns SkipDir
// when invoked on a non-directory file, Walk skips the remaining files in the
// containing directory.
type WalkFunc func(path string, info FileInfo, err error) error

// FileInfo internal file info and file contents
type FileInfo interface {
	os.FileInfo
	Compressed() bool // Is this file compressed
	Tag() string      // Etag for the file contents
	MimeType() string // file contens mimetype
	String() string   // file contents as string
	Bytes() []byte    // file contents as byte array
}

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
	subFiles   []FileInfo
}

type files struct {
	list  map[string]*file
	local bool
}

// New creates a new Files that loads embedded content.
func New(count int) FileSystem {
	return &files{list: make(map[string]*file, count)}
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

// walk recursively descends path, calling walkFn.
func (fs *files) walk(fpath string, info FileInfo, walkFn WalkFunc) (err error) {
	if !info.IsDir() {
		return walkFn(fpath, info, nil)
	}

	f := info.(*file)
	err1 := walkFn(fpath, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, entry := range f.subFiles {
		filename := path.Join(fpath, entry.Name())
		fileInfo := entry.(*file)
		err = fs.walk(filename, fileInfo, walkFn)
		if err != nil {
			if !fileInfo.IsDir() || err != SkipDir {
				return
			}
		}
	}

	return
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order.
func (fs *files) Walk(root string, walkFn WalkFunc) (err error) {
	if info, ok := fs.list[root]; !ok {
		err = walkFn(root, nil, os.ErrNotExist)
	} else {
		err = fs.walk(root, info, walkFn)
	}
	if err == SkipDir {
		err = nil
	}

	return
}

func (fs *files) Copy(target string, mode os.FileMode) error {
	mode = mode & 0777
	dirmode := os.ModeDir | ((mode & 0444) >> 2) | mode
	return fs.Walk("/", func(path string, info FileInfo, err error) error {
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
func (fs *files) AddFolder(path string, name string, local string, modtime int64, paths ...string) {
	subFiles := make([]FileInfo, len(paths))

	for i, e := range paths {
		subFiles[i] = fs.list[e]
	}

	fs.list[path] = &file{
		name:     name,
		local:    local,
		isDir:    true,
		modtime:  modtime,
		subFiles: subFiles,
	}
}

func (fs *files) addToFolder(filename string) (err error) {
	folder := path.Dir(filename)

	if f, ok := fs.list[folder]; !ok {
		fs.list[folder] = &file{
			name:    path.Base(folder),
			isDir:   true,
			modtime: time.Now().Unix(),
		}
		err = fs.addToFolder(folder)
		if err != nil {
			delete(fs.list, folder)
		}
	} else {
		if !f.isDir {
			err = &os.PathError{Op: "mkdir", Path: folder, Err: os.ErrInvalid}
		}
	}

	if err == nil {
		f := fs.list[folder]
		f.subFiles = append(f.subFiles, fs.list[filename])
		if len(f.subFiles) > 1 {
			sort.Slice(f.subFiles, func(i, j int) bool {
				return strings.Compare(f.subFiles[i].Name(), f.subFiles[j].Name()) == -1
			})
		}
	}

	return
}

func (fs *files) WriteFile(filename string, data []byte, perm os.FileMode) (err error) {

	// Make a copy of the byte data
	local := make([]byte, len(data))
	copy(local, data)
	localStr := *(*string)(unsafe.Pointer(&local))

	// If file exists just replace the data
	if f, ok := fs.list[filename]; ok {
		f.tag = ""
		f.mimeType = ""
		f.size = int64(len(local))
		f.compressed = false
		f.data = local
		f.str = localStr
	} else {
		fs.list[filename] = &file{
			name:    path.Base(filename),
			size:    int64(len(local)),
			modtime: time.Now().Unix(),
			data:    local,
			str:     localStr,
		}

		// Now add folder entries as required
		if err = fs.addToFolder(filename); err != nil {
			delete(fs.list, filename)
		}
	}

	return
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
				entries := r.subFiles[r.position : r.position+int64(count)]
				r.position += int64(count)
				list = make([]os.FileInfo, len(entries))
				for i, e := range entries {
					list[i] = e
				}
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
