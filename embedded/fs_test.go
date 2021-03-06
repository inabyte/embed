package embedded

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

const (
	index     = "<html><header></header><body></body</html>"
	mimeType  = "text/html; charset=utf-8"
	setTime   = 1579282495
	indexSize = int64(len(index))
)

var (
	indexBytes      = []byte(index)
	indexCompressed = compress(indexBytes)
	indexTag        = tag(indexBytes)
)

func TestWriteFile(t *testing.T) {
	dir, f := makeFs()
	defer os.RemoveAll(dir)

	for _, test := range []struct {
		name     string
		file     string
		hasError bool
		data     []byte
	}{
		{"Replace", "/settings.html", false, []byte("some data")},
		{"Add", "/images/picture.svg", false, []byte("some image data")},
		{"Bad Path", "/settings.html/js/utils.js", true, []byte("some script data")},
	} {
		t.Run(test.name, func(t *testing.T) {

			err := f.WriteFile(test.file, test.data, os.ModePerm)

			if err == nil {
				if test.hasError {
					t.Errorf("Test %s did not error as expected", test.name)
				}
			} else {
				if !test.hasError {
					t.Errorf("Test %s returned unexpected error %v", test.name, err)
				}
			}

			if !test.hasError {
				if file, err := f.Open(test.file); err != nil {
					t.Errorf("Test %s open file returned unexpected error %v", test.name, err)
				} else {
					if b, err := ioutil.ReadAll(file); err != nil {
						t.Errorf("Test %s ReadAll returned unexpected error %v", test.name, err)
					} else {
						if !reflect.DeepEqual(b, test.data) {
							t.Errorf("Did not get expected contents (%s) expect (%s)", b, test.data)
						}
					}
				}

			}
		})
	}
}

func TestWalk(t *testing.T) {
	dir, f := makeFs()
	defer os.RemoveAll(dir)

	for _, test := range []struct {
		name   string
		file   string
		err    error
		expect []string
	}{
		{"Bad Path", "bad", nil, []string{"bad"}},
		{"Good", "/", nil, []string{"/", "/index.html", "/settings.html", "/files", "/files/js", "/files/js/index.html"}},
		{"Skip", "/", filepath.SkipDir, []string{"/"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var list []string
			f.Walk(test.file, func(path string, info FileInfo, err error) error {
				list = append(list, path)
				return test.err
			})

			if !reflect.DeepEqual(list, test.expect) {
				t.Errorf("Did not Walk file as expected got(%v) expect (%v)", list, test.expect)
			}
		})
	}

}

func TestAdd(t *testing.T) {

	dir, f := makeFs()
	defer os.RemoveAll(dir)

	for _, test := range []struct {
		name       string
		file       string
		isFile     bool
		hasError   bool
		data       []byte
		compressed bool
		size       int64
	}{
		{
			name:     "Duplicate folder",
			file:     "/files/js",
			hasError: true,
		},
		{
			name:       "Duplicate file",
			file:       "/index.html",
			isFile:     true,
			hasError:   true,
			data:       indexCompressed,
			compressed: true,
			size:       int64(len(index)),
		},
		{
			name:       "Bad path",
			file:       "/index.html/file.html",
			isFile:     true,
			hasError:   true,
			data:       indexCompressed,
			compressed: true,
			size:       int64(len(index)),
		},
		{
			name:       "Duplicate file",
			file:       "/test/files/index.html",
			isFile:     true,
			data:       indexCompressed,
			compressed: true,
			size:       int64(len(index)),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var err error

			if test.isFile {
				str := *(*string)(unsafe.Pointer(&test.data))
				err = f.AddFile(test.file, path.Base(test.file), "", test.size, time.Now().Unix(), "", "",
					test.compressed, test.data, str)
			} else {
				err = f.AddFolder(test.file, path.Base(test.file), "", time.Now().Unix())
			}

			if err == nil {
				if test.hasError {
					t.Errorf("Did no get error as expected when adding %s", test.file)
				}
			} else {
				if !test.hasError {
					t.Errorf("Got unexpected when adding when adding %s %v", test.file, err)
				}
			}
		})
	}
}

func TestFiles(t *testing.T) {

	dir, f := makeFs()
	defer os.RemoveAll(dir)

	for _, test := range []struct {
		name       string
		file       string
		notFile    bool
		expect     []byte
		isDir      bool
		local      bool
		compressed bool
	}{
		{
			name:    "no file",
			file:    "/indexs.html",
			notFile: true,
		},
		{
			name:       "index",
			file:       "/index.html",
			expect:     indexBytes,
			compressed: true,
		},
		{
			name:       "index no path",
			file:       "index.html",
			expect:     indexBytes,
			compressed: true,
		},
		{
			name:   "settings",
			file:   "/settings.html",
			expect: indexBytes,
		},
		{
			name:   "index local",
			file:   "/index.html",
			expect: indexBytes,
			local:  true,
		},
		{
			name:  "folder",
			file:  "/",
			isDir: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {

			f.UseLocal(test.local)

			file, err := f.Open(test.file)

			if err == nil {
				if test.notFile {
					t.Errorf("Did not expect to get file %s", test.file)
				}

				testStat(t, test.file, file, test.isDir, test.local)

				if !test.isDir {
					testInfo(t, test.file, file, test.compressed, test.local)
				}

				if b, err := ioutil.ReadAll(file); err == nil {
					if test.isDir {
						t.Errorf("Read on folder did not return error for %s", test.file)
					} else {
						if !reflect.DeepEqual(b, test.expect) {
							t.Errorf("Returned data did not match get (%s) expectd(%s)", b, test.expect)
						}
					}
				} else {
					if !test.isDir {
						t.Errorf("Read on file return error for %s %v", test.file, err)
					}
				}

				if !test.isDir {
					testSeekStart(t, test.file, file, test.expect)
					testSeekCurrent(t, test.file, file, test.expect)
					testSeekEnd(t, test.file, file, test.expect)
					testSeekBad(t, test.file, file, test.expect)
				}

				testReadDir(t, test.file, file, test.isDir)

				file.Close()
				testClosed(t, test.file, file)

			} else {
				if !test.notFile {
					t.Errorf("Got error opening %s %v", test.file, err)
				}

				if file != nil {
					t.Errorf("Got bogus file for %s %v", test.file, file)
				}
			}
		})
	}
}

func makeFs() (string, FileSystem) {
	tmpdir, _ := ioutil.TempDir("", "fs-test")

	f := New(6)

	f.AddFile("/index.html", "index.html",
		filepath.Join(tmpdir, "index.html"),
		indexSize, setTime, mimeType, indexTag, true, indexCompressed, string(indexCompressed))

	f.AddFile("/settings.html", "settings.html",
		filepath.Join(tmpdir, "settings.html"),
		indexSize, setTime, mimeType, indexTag, false, indexBytes, index)

	f.AddFile("/files/fs/index.html", "index.html",
		filepath.Join(tmpdir, "files", "fs", "index.html"),
		indexSize, setTime, mimeType, indexTag, true, indexCompressed, string(indexCompressed))

	f.AddFolder("/files/js", "js",
		filepath.Join(tmpdir, "files", "js"),
		setTime,
		"/files/fs/index.html",
	)

	f.AddFolder("/files", "files",
		filepath.Join(tmpdir, "files"),
		setTime,
		"/files/js",
	)

	f.AddFolder("/", "/",
		tmpdir,
		setTime,
		"/index.html",
		"/settings.html",
		"/files",
	)

	f.Copy(tmpdir, os.ModePerm)
	return tmpdir, f
}

func testSeekStart(t *testing.T, name string, file http.File, expect []byte) {
	// Test SeekStart
	if n, err := file.Seek(10, io.SeekStart); err == nil {
		if b, err := ioutil.ReadAll(file); err == nil {
			if !reflect.DeepEqual(b, expect[10:]) {
				t.Errorf("Returned data did not match get (%s) expectd(%s)", b, expect[10:])
			}
		} else {
			t.Errorf("Read on file after io.SeekStart return error for %s %v", name, err)
		}

		if n != 10 {
			t.Errorf("Seek io.SeekStart on file return %d expected 10 for %s", n, name)
		}
	} else {
		t.Errorf("Seek io.SeekStart on file return error for %s %v", name, err)
	}
}

func testSeekCurrent(t *testing.T, name string, file http.File, expect []byte) {
	// Test SeekCurrent
	if _, err := file.Seek(0, io.SeekStart); err == nil {
		if n, err := file.Seek(10, io.SeekCurrent); err == nil {
			if b, err := ioutil.ReadAll(file); err == nil {
				if !reflect.DeepEqual(b, expect[10:]) {
					t.Errorf("Returned data did not match get (%s) expectd(%s)", b, expect[10:])
				}
			} else {
				t.Errorf("Read on file after io.SeekCurrent return error for %s %v", name, err)
			}
			if n != 10 {
				t.Errorf("Seek io.SeekStart on file return %d expected 10 for %s", n, name)
			}
		} else {
			t.Errorf("Seek io.SeekCurrent on file return error for %s %v", name, err)
		}
	}
}

func testSeekEnd(t *testing.T, name string, file http.File, expect []byte) {

	// Test SeekEnd
	if n, err := file.Seek(-10, io.SeekEnd); err == nil {
		if b, err := ioutil.ReadAll(file); err == nil {
			if !reflect.DeepEqual(b, expect[len(expect)-10:]) {
				t.Errorf("Returned data did not match get (%s) expectd(%s)", b, expect[len(expect)-10:])
			}
		} else {
			t.Errorf("Read on file after io.SeekEnd return error for %s %v", name, err)
		}
		l := int64(len(expect)) - 10
		if n != l {
			t.Errorf("Seek io.SeekStart on file return %d expected %d for %s", n, l, name)
		}
	} else {
		t.Errorf("Seek io.SeekCurrent on file return error for %s %v", name, err)
	}

}

func testSeekBad(t *testing.T, name string, file http.File, expect []byte) {
	// test seek before start of file
	if _, err := file.Seek(-10, io.SeekStart); err == nil {
		t.Errorf("Did no get error for seek before start of file %s", name)
	}

	// Test bad whence
	if _, err := file.Seek(0, -1); err == nil {
		t.Errorf("Did no get error for seek bad whence of file %s", name)
	}

}

func testStat(t *testing.T, name string, file http.File, isDir bool, local bool) {
	if info, err := file.Stat(); err != nil {
		t.Errorf("Stat return error for file %s %v", name, err)
	} else {
		base := path.Base(name)
		if n := info.Name(); n != base {
			t.Errorf("Name did not return valid value got (%s) expected %s", n, base)
		}

		modTime := time.Unix(setTime, 0)
		if n := info.ModTime(); n != modTime {
			t.Errorf("ModTime did not return valid value got (%v) expected %v", n, modTime)
		}

		if is := info.IsDir(); is != isDir {
			t.Errorf("IsDir did not return valid value got (%v) expected %v", is, isDir)
		}

		if sys := info.Sys(); sys == nil {
			t.Errorf("Sys did not return underlying value for %s", name)
		}

		expectSize := indexSize
		if isDir {
			expectSize = 0
		}

		if size := info.Size(); size != expectSize {
			t.Errorf("Size unexpected got (%v) expected %v for %s", size, expectSize, name)
		}

		if !local {
			expectMode := os.ModePerm | os.ModeDir
			if !isDir {
				expectMode = os.ModePerm
			}

			if mode := info.Mode(); mode != expectMode {
				t.Errorf("Mode unexpected got (%v) expected %v for %s", mode, expectMode, name)
			}
		}

	}
}

func testReadDir(t *testing.T, name string, file http.File, isDir bool) {
	if _, err := file.Readdir(-1); err != nil {
		if isDir {
			t.Errorf("Got error Readdir for %s %v", name, err)
		}
	} else {
		if !isDir {
			t.Errorf("Did no get error Readdir for %s", name)
		} else {
			if _, err := file.Readdir(1); err == nil {
				t.Errorf("Did not get error recall Readdir for %s", name)
			}

			if _, err := file.Seek(10, io.SeekCurrent); err == nil {
				t.Errorf("Seek on folder did no return an err for %s", name)
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				t.Errorf("Seek to beginning of folder return error for %s %v", name, err)
			}
		}
	}
}

func testInfo(t *testing.T, name string, file http.File, compressed bool, local bool) {
	if info, ok := file.(FileInfo); ok {
		if compress := info.Compressed(); compress != compressed {
			t.Errorf("FileInfo.Compressed unexpected got (%v) expected %v for %s", compress, compressed, name)
		}

		if tag := info.Tag(); tag != indexTag {
			t.Errorf("FileInfo.Tag unexpected got (%v) expected %v for %s", tag, indexTag, name)
		}

		if mime := info.MimeType(); mime != mimeType {
			t.Errorf("FileInfo.MimeType unexpected got (%v) expected %v for %s", mime, mimeType, name)
		}

		if str := info.String(); str != index {
			t.Errorf("FileInfo.String unexpected got (%v) expected %v for %s", str, index, name)
		}

		if buf := info.Bytes(); !reflect.DeepEqual(buf, indexBytes) {
			t.Errorf("FileInfo.String unexpected got (%v) expected %v for %s", buf, indexBytes, name)
		}

		expect := indexBytes
		if info.Compressed() {
			expect = indexCompressed
		}

		if buf := info.Raw(); !reflect.DeepEqual(buf, expect) {
			t.Errorf("FileInfo.String unexpected got (%v) expected %v for %s", buf, expect, name)
		}

	} else {
		if !local {
			t.Errorf("Did not get FileInfo interface for %s", name)
		}
	}
}

func testClosed(t *testing.T, name string, file http.File) {
	_, err := file.Stat()
	if err == nil {
		t.Errorf("Stat on closed file did not error %s", name)
	}

	_, err = file.Readdir(-1)
	if err == nil {
		t.Errorf("Readdir on closed file did not error %s", name)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err == nil {
		t.Errorf("Seek on closed file did not error %s", name)
	}

	_, err = file.Read(make([]byte, 10))
	if err == nil {
		t.Errorf("Read on closed file did not error %s", name)
	}

	err = file.Close()
	if err == nil {
		t.Errorf("Close on closed file did not error %s", name)
	}
}

func compress(data []byte) []byte {
	var buf bytes.Buffer

	gw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)

	if err == nil {
		_, err = gw.Write(data)
	}

	if err == nil {
		err = gw.Close()
	}

	if err == nil {
		return buf.Bytes()
	}

	return data
}

func tag(data []byte) string {
	hash := sha1.Sum(data)
	return base64.RawURLEncoding.EncodeToString(hash[:]) + "-gz"
}
