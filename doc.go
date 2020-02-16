/*
Package embed takes a list for file or folders (likely at `go generate` time) and
generates Go code that statically implements the a http.FileSystem.

Basics

Embeds files into go programs and provides http.FileSystem interfaces
to them.

It adds all named files or files recursively under named directories at the
path specified. The output file provides an http.FileSystem interface with
optionally zero dependencies on packages outside the standard library.

Features

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

Source

The source paramaters are a list of folder and/or files, with an optional marker for prefix start.
The prefix marker is <->, if you specify a path of ./files/html<->/embed it will import
everything under the folder ./files/html/embed as /embed/... for example
./files/html/embed/index.html will be stored as /embed/index.html.
There is no limit to the number you can specify but the resultant file system must be unique,
if the processing produces two files in the same location Generate() exit with an error.

Example

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

	package main

	import (
	    "log"
	    "net/http"
	)

	func main() {
	    // FileHandler() is created by embed and returns a http.Handler.
	    log.Fatal(http.ListenAndServe(":8080", FileHandler()))
	}


And embed_static.go contains:

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

1. Generate the embedded data:
	go run embed_static.go
2. Start the server:
	go run .
3. Access http://localhost:8080/ to view the content.

*/
package embed
