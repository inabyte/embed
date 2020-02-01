# Embed static files into go binaries

[![Build Status](https://travis-ci.com/inabyte/embed.svg?branch=master)](https://travis-ci.com/inabyte/embed)

Takes a list for file or folders (likely at `go generate` time) and
generates Go code that statically implements the a http.FileSystem.

Features:

-	Efficient generated code without unneccessary overhead.

-	Minimizes html css and js files.

-	Uses gzip compression internally (selectively, only for files that compress well).

-	Outputs `gofmt`ed Go code.


## Installation

First, get the command-line tool:

    go get github.com/inabyte/embed/cmd/embed 
