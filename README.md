# Embed static files into go binaries

[![GoDoc](https://godoc.org/github.com/inabyte/embed?status.svg)](https://godoc.org/github.com/inabyte/embed)
[![Build Status](https://travis-ci.com/inabyte/embed.svg?branch=master)](https://travis-ci.com/inabyte/embed)
[![Coverage Status](https://coveralls.io/repos/github/inabyte/embed/badge.svg?branch=master)](https://coveralls.io/github/inabyte/embed?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/inabyte/embed)](https://goreportcard.com/report/github.com/inabyte/embed)
![GitHub](https://img.shields.io/github/license/inabyte/embed)

Takes a list for file or folders (likely at `go generate` time) and
generates Go code that statically implements the a http.FileSystem.

Features:

`embed`:

- outputs `gofmt`ed and `lint`ed Go code.
- produces go-gettable go and go assembly sources with `go generate`.
- keeps data and strings in read-only section of the binary.
- minify HTML, CSS, and JavaScript files.
- compress compressible files with `gzip`.
- provides [http.FileSystem](https://golang.org/pkg/net/http/#FileSystem) API.
- provides [http.Handler](https://golang.org/pkg/net/http/#Handler) handler
  (if requested),
- generates a test code. 
- allows access to compressed data.
- serving gzip files to clients that accept compressed content.
- calculate checksums for serving Etag-based conditional requests.
- zero dependencies on packages outside the standard library
  (if requested),


The minifier used
is [github.com/tdewolff/minify](https://github.com/tdewolff/minify).

## Source file/folder specification

The source paramaters are a list of folder and/or files, with an optional marker for prefix start. 
The prefix marker is `<->`, if you specify a path of `./files/html<->/embed` it will import 
everything under the folder `./files/html/embed` as `/embed/...` for example 
`./files/html/embed/index.html` will be stored as `/embed/index.html`. 
There is no limit to the number you can specify but the resultant file system must be unique, 
if the processing produces two files in the same location `embed` will error out.

# Usage as Binary

## Installation

First, get the command-line tool:

    go get github.com/inabyte/embed/cmd/embed 

## Usage

`embed [flags] [files...]`

Flags Are:

```
-o="embed"
  Output files base name.
-pkg=""
  Package name defauls to directory of output.
-tags=""
  Build tags added to output files.
-ignore=""
  Regexp for files we should ignore (for example \\\\.DS_Store).
-include=""
  Regexp for files to include. Only files that match will be included.
-minify="application/javascript,text/javascript,text/css,text/html,text/html; charset=utf-8"
  Comma list of mimetypes to minify.
-modifytime=""
  Unix timestamp to override as modification time for all files.
-no-compress
  If set, do not compress files.
-go
  If set, write only go files
-fileserver
  If set, produce http server code
-binary
  If set, produce self-contained extractor/http server binary (-o will become the binary name)
-noremote
  If set, force zero dependencies on packages outside the standard library.
-nolocalfs
  If set, do not store local file system paths.
```

## Example

Embedded assets can be served with HTTP using the `http.Server`.
Assuming a directory structure similar to below:

	.
	├── main.go
	└── html
	    ├── css
	    │   └── style.css
	    ├── scripts
	    │   └── utils.js
	    └── index.html


Where main.go contains:

```
package main

import (
  "log"
  "net/http"
)

func main() {
  // FileHandler() is created by embed and returns a http.Handler.
  log.Fatal(http.ListenAndServe(":8080", FileHandler()))
}
```

1. Generate the embedded data:
	`embed -o static -pkg=main html`
2. Start the server:
	`go run .`
3. Access http://localhost:8080/ to view the content.

You can see a worked example in [examples/binary](examples/binary) dir
just run it as
`go run ./examples/binary`



# Usage with Import


## Config

```
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
	// Minify is comma separated list of mime type to minify.
	Minify string
	// ModifyTime is the Unix timestamp to override as modification time for all files.
	ModifyTime string
	// DisableCompression, if true, does not compress files.
	DisableCompression bool
	// Binary, if true, produce self-contained extractor/http server binary.
	Binary bool
	// NoRemote, if true, zero dependencies on packages outside the standard library.
	NoRemote bool
	// Go, if true, creates only go files.
	Go bool
	// NoLocalFS, if true, do not store local file system paths.
	NoLocalFS bool
	// FileServer, if true, add http.Handler to serve files.
	FileServer bool
	// BuildTags, if set, adds a build tags entry to file.
	BuildTags string
	// Files is the list of files or directories to embed.
	Files []string
}
```

## Example

Embedded assets can be served with HTTP using the `http.Server`.
Assuming a directory structure similar to below:

	.
	├── main.go
	└── html
	    ├── css
	    │   └── style.css
	    ├── scripts
	    │   └── utils.js
	    └── index.html


Where main.go contains:

```
package main

import (
  "log"
  "net/http"
)

func main() {
  // FileHandler() is created by embed and returns a http.Handler.
  log.Fatal(http.ListenAndServe(":8080", FileHandler()))
}
```

And embed_static.go contains:

```
// +build ignore
package main

import (
  "github.com/inabyte/embed"
)

func main() {
  config := embed.New()

  config.Output = "static"
  config.Package = "main"
  config.Files = []string{"html"}
  config.FileServer = true
  config.NoLocalFS = true

  config.Generate()
}
```

1. Generate the embedded data:
	`go run embed_static.go`
2. Start the server:
	`go run .`
3. Access http://localhost:8080/ to view the content.

You can see a worked example in [examples/code](examples/code) dir
just run it as
`go run ./examples/code`

