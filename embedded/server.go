// Package embedded embed files
//
// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package embedded

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type server struct {
	http.FileSystem
	notFound      http.Handler
	sys           http.Handler
	renderFolders bool
}

// GetFileServer create a http.handler server
func GetFileServer(fs http.FileSystem) Handler {
	return &server{
		FileSystem:    fs,
		sys:           http.FileServer(fs),
		renderFolders: true,
	}
}

// SetNotFoundHandler set a hander to be called for no found
func (s *server) SetNotFoundHandler(h http.Handler) {
	s.notFound = h
}

// If true and the folder does not contain index.html render folder
// otherwise call the NotFound handler
func (s *server) SetRenderFolders(enable bool) {
	s.renderFolders = enable
}

// ServeHTTP implement http.Handler interface
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const indexPage = "/index.html"

	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	name := path.Clean(upath)

	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		localRedirect(w, r, "./")
		return
	}

	var d os.FileInfo

	f, err := s.Open(name)
	if err == nil {
		defer f.Close()
		d, err = f.Stat()
	}

	if err != nil {
		s.toHTTPError(w, r, err)
		return
	}

	// redirect to canonical path: / at end of directory url
	// r.URL.Path always begins with /
	url := r.URL.Path
	if d.IsDir() {
		if url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}
	} else {
		if url[len(url)-1] == '/' {
			localRedirect(w, r, "../"+path.Base(url))
			return
		}
	}

	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := s.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = index
				d = dd
				f = ff
			}
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		if s.renderFolders {
			s.sys.ServeHTTP(w, r)
		} else {
			s.toHTTPError(w, r, os.ErrPermission)
		}
		s.sys.ServeHTTP(w, r)
		return
	}

	if reader, ok := f.(*reader); ok {
		reader.serve(w, r)
	} else {
		// ServeContent will check modification time
		http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	}
}

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func (s *server) toHTTPError(w http.ResponseWriter, r *http.Request, err error) {

	httpStatus := http.StatusInternalServerError

	if os.IsNotExist(err) {
		httpStatus = http.StatusNotFound
		if s.notFound != nil {
			s.notFound.ServeHTTP(w, r)
			return
		}
	} else {
		if os.IsPermission(err) {
			httpStatus = http.StatusForbidden
		}
	}

	http.Error(w, http.StatusText(httpStatus), httpStatus)
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

// serve set various headers etag, content type
func (f *reader) serve(w http.ResponseWriter, r *http.Request) {
	tag := f.tag
	// Check is requesting compressed and we have it compressed
	if f.compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		f.readCompressed = true
		f.length = int64(len(f.data))
	} else {
		tag = tag[:len(tag)-3]
	}

	w.Header().Set("Content-Type", f.mimeType)
	w.Header().Set("Etag", strconv.Quote(tag))

	// ServeContent will check modification time
	http.ServeContent(w, r, f.Name(), f.ModTime(), f)
}
