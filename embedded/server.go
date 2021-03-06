package embedded

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Handler serves embedded handle to serve FileSystem
type Handler interface {
	http.Handler
	// SetNotFoundHandler set a hander to be called for no found
	SetNotFoundHandler(http.Handler)
	// SetPermissionHandler set a hander to be called for permission http.StatusForbidden
	SetPermissionHandler(http.Handler)
	// If true and the folder does not contain index.html render folder
	// otherwise return 403 http.StatusForbidden
	SetRenderFolders(enable bool)
}

type server struct {
	http.FileSystem
	notFound      http.Handler
	permission    http.Handler
	sys           http.Handler
	renderFolders bool
}

// GetFileServer create a http.handler
// this will serve file from the embedded FileSystem.
// Serving gzip files to clients that accept compressed content if the content is compressed.
// Serving Etag-based conditional requests if specified.
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

// SetPermissionHandler set a hander to be called for no found
func (s *server) SetPermissionHandler(h http.Handler) {
	s.permission = h
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
		return
	}

	// If the files implements a http.Handler user that otherwise pass to http.ServerContent
	if handler, ok := f.(http.Handler); ok {
		handler.ServeHTTP(w, r)
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
			if s.permission != nil {
				s.permission.ServeHTTP(w, r)
				return
			}
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

// ServeHTTP set various headers etag, content type and serve file
func (f *reader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tag := f.tag

	// Check is requesting compressed and we have it compressed
	if f.compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		f.readCompressed = true
		f.length = int64(len(f.data))
	} else {
		if len(tag) > 3 {
			tag = tag[:len(tag)-3]
		}
	}

	if len(f.mimeType) > 0 {
		w.Header().Set("Content-Type", f.mimeType)
	}

	if len(tag) > 0 {
		w.Header().Set("Etag", strconv.Quote(tag))
	}

	// ServeContent will check modification time
	http.ServeContent(w, r, f.Name(), f.ModTime(), f)
}
