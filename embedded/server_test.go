package embedded

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestLocalWebServe(t *testing.T) {
	dir, fs := makeFs()
	defer os.RemoveAll(dir)

	s := GetFileServer(fs)
	s.SetNotFoundHandler(http.HandlerFunc(http.NotFound))
	s.SetPermissionHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}))

	for _, tt := range []struct {
		name          string
		url           string
		override      string
		requestGZip   bool
		local         bool
		renderFolders bool
		expect        []byte
		status        int
	}{
		{"index compressed", "/", "", true, false, true, indexCompressed, http.StatusOK},
		{"index uncompressed", "/", "", false, false, true, indexBytes, http.StatusOK},
		{"index redirect", "/index.html?debug=true", "", true, false, true, nil, http.StatusMovedPermanently},
		{"index add slash", "/index.html", "index.html", false, false, true, nil, http.StatusMovedPermanently},
		{"not found", "/bad", "", false, false, true, nil, http.StatusNotFound},
		{"norender folder", "/files/js/", "", true, false, false, nil, http.StatusForbidden},
		{"folder", "/files/js/", "", true, false, true, nil, http.StatusOK},
		{"folder redirect", "/files/js", "", true, false, true, nil, http.StatusMovedPermanently},
		{"file redirect", "/settings.html/", "", true, false, true, nil, http.StatusMovedPermanently},
		{"settings", "/settings.html", "", true, false, true, indexBytes, http.StatusOK},
		{"local fs", "/settings.html", "", true, true, true, indexBytes, http.StatusOK},
	} {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tt.url, nil)
			if len(tt.override) > 0 {
				r.URL.Path = tt.override
			}

			if tt.requestGZip {
				r.Header.Add("Accept-Encoding", "gzip")
			}

			fs.UseLocal(tt.local)
			s.SetRenderFolders(tt.renderFolders)
			s.ServeHTTP(w, r)

			if tt.status != 0 && w.Result().StatusCode != tt.status {
				t.Errorf("%q. result = %d, want %d", tt.name, w.Result().StatusCode, tt.status)
			}

			if tt.expect != nil {
				if !reflect.DeepEqual(w.Body.Bytes(), tt.expect) {
					t.Errorf("%q. got = (%s), want (%s)", tt.name, w.Body.Bytes(), tt.expect)
				}
			}
		})
	}
}
func TestToHTTPError(t *testing.T) {

	s := GetFileServer(nil)

	srv := s.(*server)

	for _, tt := range []struct {
		name   string
		err    error
		status int
	}{
		{"not found", os.ErrNotExist, http.StatusNotFound},
		{"permission", os.ErrPermission, http.StatusForbidden},
		{"server", os.ErrClosed, http.StatusInternalServerError},
	} {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			srv.toHTTPError(w, nil, tt.err)

			if tt.status != 0 && w.Result().StatusCode != tt.status {
				t.Errorf("%q. result = %d, want %d", tt.name, w.Result().StatusCode, tt.status)
			}
		})
	}

}
