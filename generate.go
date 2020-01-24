package embed

// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/inabyte/embed/embedded"
	"github.com/inabyte/embed/internal/templates"
)

// New create new config
func New() *Config {
	return &Config{
		Output: "embed",
		Minify: "application/javascript,text/javascript,text/css,text/html,text/html; charset=utf-8",
	}
}

// Generate create embedded files.
func (config *Config) Generate() error {
	var gen generate

	return gen.generate(config)
}

// Config contains all information needed to run embed.
type Config struct {
	// Output is the file to write output.
	Output string
	// Package name for the generated file.
	Package string
	// Ignore is the regexp for files we should ignore (for example `\.DS_Store`).
	Ignore string
	// Include is the regexp for files to include. If provided, only files that
	// match will be included.
	Include string
	// Minify is comma seperated list of mime type to minify.
	Minify string
	// ModifyTime is the Unix timestamp to override as modification time for all files.
	ModifyTime string
	// DisableCompression, if true, does not compress files.
	DisableCompression bool
	// Binary, if true, produce self-contained extracter/http server binary.
	Binary bool
	// Local, if true, creates all file locally.
	Local bool
	// Go, if true, creates only go files.
	Go bool
	// BuildTags, if set, adds a build tags entry to file.
	BuildTags string
	// Files is the list of files or directories to embed.
	Files []string
}

const (
	header       = "// Code generated by embed. DO NOT EDIT."
	prefixMarker = "<->"
)

var (
	tmpl     = template.Must(template.New("").Parse(fileTemplate))
	testTmpl = template.Must(template.New("").Parse(testTemplate))
)

type generate struct {
	Remote      bool
	PackageName string
	BuildTags   string
	Imports     []string
	TestImports []string
	Main        bool
	Go          bool
	Files       []*file
	Dirs        []*dir
	ignore      *regexp.Regexp
	include     *regexp.Regexp
	imports     map[string]bool
	testImports map[string]bool
	minify      map[string]bool
	modifyTime  *int64
	prefix      string
	compress    bool
	Offset      int64
	processed   map[string]bool
	config      *Config
	last        time.Time
}

func (gen *generate) init(config *Config) (err error) {

	gen.Remote = !config.Local
	gen.config = config

	gen.Main = config.Binary

	gen.Go = config.Go

	gen.minify = make(map[string]bool)
	for _, e := range strings.Split(config.Minify, ",") {
		s := strings.TrimSpace(e)
		if len(s) > 0 {
			gen.minify[s] = true
		}
	}

	gen.BuildTags = config.BuildTags

	gen.PackageName = config.Package
	if len(gen.PackageName) == 0 {
		pkg := filepath.Base(filepath.Dir(config.Output))

		gen.PackageName = pkg
	}

	gen.Files = make([]*file, 0, 10)
	gen.Dirs = make([]*dir, 0, 10)
	gen.processed = make(map[string]bool, 10)
	gen.compress = !config.DisableCompression

	if config.ModifyTime != "" {
		if i, e := strconv.ParseInt(config.ModifyTime, 10, 64); e != nil {
			err = fmt.Errorf("ModifyTime must be an integer: %v", e)
		} else {
			gen.modifyTime = &i
		}
	}

	if err == nil && config.Ignore != "" {
		gen.ignore, err = regexp.Compile(config.Ignore)
	}

	if err == nil && config.Include != "" {
		gen.include, err = regexp.Compile(config.Include)
	}

	gen.imports = make(map[string]bool)
	gen.testImports = map[string]bool{"crypto/sha1": true, "encoding/base64": true, "testing": true}

	if gen.Main {
		gen.imports["flag"] = true
		gen.imports["net/http"] = true
		gen.imports["os"] = true
	}

	if gen.Go {
		gen.imports["reflect"] = true
		gen.imports["unsafe"] = true
	}

	if gen.Remote {
		gen.imports["github.com/inabyte/embed/embedded"] = true
		gen.testImports["github.com/inabyte/embed/embedded"] = true
	}

	sort.Strings(gen.Imports)
	sort.Strings(gen.TestImports)

	return
}

func (gen *generate) skip(name string) bool {
	if gen.ignore != nil && gen.ignore.MatchString(name) {
		return true
	}
	if gen.include == nil || gen.include.MatchString(name) {
		return false
	}
	return true
}

func (gen *generate) setLast(fi os.FileInfo) {
	if gen.last.Before(fi.ModTime()) {
		gen.last = fi.ModTime()
	}
}

func (gen *generate) getModTime(t time.Time) (m int64) {
	m = t.Unix()
	if gen.modifyTime != nil {
		m = *gen.modifyTime
	}
	return
}

func (gen *generate) checkProcessed(name string, fpath string) (err error) {
	if gen.processed[name] {
		err = fmt.Errorf("%s, %s: duplicate Name after prefix removal", name, fpath)
	} else {
		gen.processed[name] = true
	}

	return
}

func (gen *generate) canonicalName(fname string) string {
	fpath := filepath.ToSlash(fname)
	return path.Join("/", strings.TrimPrefix(fpath, gen.prefix))
}

func (gen *generate) process(fpath string) error {
	var fi os.FileInfo

	gen.prefix = fpath
	if n := strings.Index(gen.prefix, prefixMarker); n >= 0 {
		gen.prefix = gen.prefix[:n]
		fpath = strings.Replace(fpath, prefixMarker, "", 1)
	}

	f, err := os.Open(fpath)
	if err == nil {
		fi, err = f.Stat()
		if err == nil {
			if fi.IsDir() {
				err = gen.processDir(fpath, f, fi)
			} else {
				err = gen.processFile(fpath, f, fi)
			}
		}
		f.Close()
	}

	return err
}

func (gen *generate) processDir(fpath string, f *os.File, fi os.FileInfo) error {
	var (
		d   *dir
		sub *os.File
	)

	n := gen.canonicalName(fpath)
	fis, err := f.Readdir(0)

	if err == nil {
		err = gen.checkProcessed(n, fpath)
	}

	if err == nil {
		gen.setLast(fi)
		d = &dir{
			name:     n,
			baseName: path.Base(n),
			local:    fpath,
			ModTime:  gen.getModTime(fi.ModTime()),
			files:    make(map[string]bool, len(fis)),
		}
		for _, fi := range fis {
			if err == nil {
				name := filepath.Join(fpath, fi.Name())
				if sub, err = os.Open(name); err == nil {
					n := gen.canonicalName(name)
					if fi.IsDir() {
						d.files[n] = true
						gen.processDir(name, sub, fi)
					} else {
						if !gen.skip(name) {
							d.files[n] = true
							gen.processFile(name, sub, fi)
						}
					}
					sub.Close()
				}
			}
		}
	}

	if err == nil {
		gen.Dirs = append(gen.Dirs, d)
	}

	return err
}

func (gen *generate) processFile(fpath string, f *os.File, fi os.FileInfo) error {
	n := gen.canonicalName(fpath)
	b, err := ioutil.ReadAll(f)

	if err == nil {
		err = gen.checkProcessed(n, fpath)
	}

	if err == nil {
		gen.setLast(fi)
		file := &file{
			name:     n,
			baseName: path.Base(n),
			data:     b,
			local:    fpath,
			ModTime:  gen.getModTime(fi.ModTime()),
		}

		gen.Files = append(gen.Files, file)
		gen.processed[n] = true
	}

	return err
}

func (gen *generate) processFiles() (err error) {

	var offset int64

	for _, entry := range gen.Files {

		if err == nil {
			entry.setMimeType()

			if gen.minify[entry.mimeType] {
				entry.minify()
			}

			entry.fill()
			entry.set()

			if gen.compress {
				err = entry.compress()
			}
		}
	}

	if err == nil {
		for _, entry := range gen.Dirs {
			entry.set()
		}

		stringer.process()
	}

	stringer.offset = int(offset)
	offset += stringer.len()

	for _, entry := range gen.Files {
		entry.Offset = offset
		offset += int64(len(entry.data))
	}

	gen.Offset = offset

	return
}

func (gen *generate) writeData(file *os.File) (err error) {

	var writer writer

	if gen.Go {
		writer, err = createGoWriter(file)
	} else {
		writer, err = createWriter(gen.config.Output, gen.config.BuildTags)
	}

	if err == nil {
		if !gen.Go {
			defer writer.Close()
		}

		if err == nil {
			_, err = writer.Write(stringer.bytes())
		}

		for _, entry := range gen.Files {
			if err == nil {
				if err == nil {
					_, err = writer.Write(entry.data)
				}
			}
		}

		if err == nil {
			writer.footer()
		}
	}

	return err
}

func (gen *generate) writeGoFiles() (err error) {
	var file *os.File

	file, err = createFile(gen.config.Output, "", ".go")

	if err == nil {
		err = tmpl.Execute(file, gen)
	}

	if err == nil && gen.config.Local && gen.Go {
		gen.appendFiles(file, false)
	}

	if err == nil && gen.Go {
		err = gen.writeData(file)
	}

	if file != nil {
		file.Close()
		file = nil
	}

	if err == nil {
		file, err = createFile(gen.config.Output, "_test", ".go")
	}

	if err == nil {
		err = testTmpl.Execute(file, gen)
	}

	if err == nil && gen.config.Local && gen.Go {
		gen.appendFiles(file, true)
	}

	if file != nil {
		file.Close()
		file = nil
	}

	if gen.config.Local && !gen.Go {
		gen.createFiles()
	}

	return
}

func (gen *generate) appendFiles(out *os.File, tests bool) {
	templates.FS.Walk(func(path string, info embedded.FileInfo) error {

		if !info.IsDir() {
			var (
				s string
			)

			file, err := templates.FS.Open(path)

			if err == nil {
				defer file.Close()
				read := bufio.NewReader(file)

				for s, err = read.ReadString('\n'); err == nil && !strings.HasPrefix(s, "import"); s, err = read.ReadString('\n') {
				}

				for s, err = read.ReadString('\n'); err == nil && !strings.HasPrefix(s, ")"); s, err = read.ReadString('\n') {
				}

				if err == nil {
					if strings.HasSuffix(info.Name(), "_test.go") {
						if tests {
							_, err = io.Copy(out, read)
						}
					} else {
						if !tests {
							_, err = io.Copy(out, read)
						}
					}
				}
			}
			return err
		}

		return nil
	})
}

func (gen *generate) createFiles() {
	templates.FS.Walk(func(path string, info embedded.FileInfo) error {

		if !info.IsDir() {
			var (
				s   string
				out *os.File
			)

			file, err := templates.FS.Open(path)

			if err == nil {
				defer file.Close()
				read := bufio.NewReader(file)

				for s, err = read.ReadString('\n'); err == nil && !strings.HasPrefix(s, "package"); s, err = read.ReadString('\n') {
				}

				if err == nil {
					out, err = createWriteHeader(gen.config.Output, "_"+info.Name(), "")
				}

				if err == nil {
					_, err = fmt.Fprintf(out, "\n\npackage %s\n", gen.PackageName)
				}

				if err == nil {
					_, err = io.Copy(out, read)
				}
			}
			return err
		}

		return nil
	})
}

func (gen *generate) scanImports() {
	templates.FS.Walk(func(path string, info embedded.FileInfo) error {

		if !info.IsDir() {
			var (
				s string
			)

			file, err := templates.FS.Open(path)

			if err == nil {
				defer file.Close()
				read := bufio.NewReader(file)

				for s, err = read.ReadString('\n'); err == nil && !strings.HasPrefix(s, "import"); s, err = read.ReadString('\n') {
				}

				for s, err = read.ReadString('\n'); err == nil && !strings.HasPrefix(s, ")"); s, err = read.ReadString('\n') {
					s = strings.TrimSpace(s)
					s = s[1 : len(s)-1] // Remove double quotes
					if strings.HasSuffix(info.Name(), "_test.go") {
						gen.testImports[s] = true
					} else {
						gen.imports[s] = true
					}
				}
			}
			return err
		}

		return nil
	})
}

func (gen *generate) addVirtualDirs() error {
	list := make(map[string]*dir, len(gen.Dirs))

	for _, k := range gen.Dirs {
		list[k.name] = k
	}

	for _, k := range gen.Dirs {
		names := strings.Split(k.name, "/")

		fpath := "/"

		for _, l := range names {
			fpath = path.Join(fpath, l)

			if entry := list[fpath]; entry == nil {

				dir := &dir{
					name:     fpath,
					baseName: path.Base(fpath),
					local:    "",
					ModTime:  gen.getModTime(gen.last),
					files:    make(map[string]bool),
				}

				list[fpath] = dir
				gen.Dirs = append(gen.Dirs, dir)

			} else {
				entry = list[path.Dir(fpath)]
				if fpath != "/" && entry != nil && !entry.files[fpath] {
					entry.files[fpath] = true
				}
			}
		}
	}
	return nil
}

func (gen *generate) generate(config *Config) error {
	err := gen.init(config)

	if err == nil {
		for _, entry := range config.Files {
			if err == nil {
				err = gen.process(entry)
			}
		}
	}

	if err == nil && len(gen.Files) == 0 {
		err = errors.New("Files empty")
	}

	if err == nil {

		if gen.Go && gen.config.Local {
			gen.scanImports()
		}

		gen.Imports = make([]string, len(gen.imports))
		i := 0
		for k := range gen.imports {
			gen.Imports[i] = k
			i++
		}

		gen.TestImports = make([]string, len(gen.testImports))
		i = 0
		for k := range gen.testImports {
			gen.TestImports[i] = k
			i++
		}

		sort.Strings(gen.Imports)
		sort.Strings(gen.TestImports)

		gen.addVirtualDirs()

		sort.Slice(gen.Files, func(i, j int) bool {
			return strings.Compare(gen.Files[i].name, gen.Files[j].name) == -1
		})

		sort.Slice(gen.Dirs, func(i, j int) bool {
			return strings.Compare(gen.Dirs[i].name, gen.Dirs[j].name) == -1
		})
	}

	if err == nil {
		err = gen.processFiles()
	}

	if err == nil && !gen.Go {
		err = gen.writeData(nil)
	}

	if err == nil {
		err = gen.writeGoFiles()
	}

	if err == nil && !gen.Go {
		err = assmemblerFiles.output(gen.config.Output, config.BuildTags)
	}

	return err
}

const (
	fileTemplate = header + `
{{ if .BuildTags }}// +build {{ .BuildTags }} 
{{ end }}
package {{.PackageName}}{{ if .Imports }}

import ({{ range .Imports }}
	"{{.}}"{{ end }}
){{ end }}

// FS return file system
var FS {{ if .Remote }}embedded.{{ end }}FileSystem

// FileHandler return http file server implements http.Handler
func FileHandler() {{ if .Remote }}embedded.{{ end }}Handler {
	return {{ if .Remote }}embedded.{{ end }}GetFileServer(FS)
}

{{ if not .Go }}func file_bytes(uint32) []byte
func file_string(uint32) string
{{ end }}func init() {
{{ if .Go }}
	bytes := dataBytes()
	str := dataString()
{{ else }}
	bytes := file_bytes({{ .Offset }})
	str := file_string({{ .Offset }})
{{ end}}
	FS = {{ if .Remote }}embedded.{{ end }}New()
{{ range .Files }}
	FS.AddFile( {{ .Name }},
		{{ .BaseName }},
		{{ .Local }},
		{{ .Size  }}, {{ .ModTime }},
		{{ .MimeType }},
		{{ .Tag }},
		{{ .Compressed }}, bytes[{{ .Slice  }}], str[{{ .Slice  }}])
{{ end -}}
{{ range .Dirs }}
	FS.AddFolder( {{ .Name }},
		{{ .BaseName }},
		{{ .Local }},
		{{ .ModTime }})
{{ end -}}
{{ range .Dirs }}
	FS.SetFiles( {{ .Name }}, 
	{{- range .Files }}
		{{.}},
	{{- end }}
	){{ end }}
}
{{- if .Main }}

func main() {
	var (
		tls        bool
		err        error
		listenAddr string
		certFile   string
		keyFile    string
		extract    string
		show       bool
	)

	flag.BoolVar(&show, "show", false, "list contents and exit")
	flag.StringVar(&extract, "extract", "", "extract contents to the target directory and exit")
	flag.StringVar(&listenAddr, "listen", ":8080", "socket address to listen")
	flag.StringVar(&certFile, "tls-cert", "", "TLS certificate file to use")
	flag.StringVar(&keyFile, "tls-key", "", "TLS key file to use")

	flag.Parse()

	if show {
		FS.Walk(func(path string, info {{ if .Remote }}embedded.{{ end }}FileInfo) error {
			if !info.IsDir() {
				os.Stdout.WriteString(path)
				os.Stdout.WriteString("\n")
			}
			return nil
		})
		return
	}
	if extract != "" {
		if err = FS.Copy(extract, 0640); err != nil {
			os.Stderr.WriteString("error extracting content: ")
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
			os.Exit(1)
		}
		return
	}
	if certFile != "" && keyFile != "" {
		tls = true
	} else if certFile != "" || keyFile != "" {
		os.Stderr.WriteString("both certFile and keyFile must be supplied for HTTPS\n")
		os.Exit(1)
	}
	if tls {
		err = http.ListenAndServeTLS(listenAddr, certFile, keyFile, FileHandler())
	} else {
		err = http.ListenAndServe(listenAddr, FileHandler())
	}
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}{{ end }}
`

	testTemplate = header + `
{{ if .BuildTags }}// +build {{ .BuildTags }} 
{{ end }}
package {{.PackageName}}
{{ if .TestImports }}
import ({{ range .TestImports }}
	"{{.}}"{{ end }}
){{ end }}

func TestFileServer(t *testing.T) {
	if FileHandler() == nil {
		t.Errorf("Call to FileServer did no return a handler")
	}
}

func TestBytes(t *testing.T) {
	FS.Walk(func(path string, info {{ if .Remote }}embedded.{{ end }}FileInfo) error {
		if !info.IsDir() {
			if getTag(info.Bytes()) != info.Tag() {
				t.Errorf("checksum for file %s doesn't match recorded", path)
			}
		}
		return nil
	})
}

func TestString(t *testing.T) {
	FS.Walk(func(path string, info {{ if .Remote }}embedded.{{ end }}FileInfo) error {
		if !info.IsDir() {
			s := info.String()
			if getTag([]byte(s)) != info.Tag() {
				t.Errorf("checksum for file %s doesn't match recorded {%s}", path, s)
			}
		}
		return nil
	})
}

func getTag(data []byte) string {
	hash := sha1.Sum(data)
	return base64.RawURLEncoding.EncodeToString(hash[:]) + "-gz"
}
`
)
