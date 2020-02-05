# Embed static files into go binaries

[![Build Status](https://travis-ci.com/inabyte/embed.svg?branch=master)](https://travis-ci.com/inabyte/embed)
[![Coverage Status](https://coveralls.io/repos/github/inabyte/embed/badge.svg?branch=master)](https://coveralls.io/github/inabyte/embed?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/inabyte/embed)](https://goreportcard.com/report/github.com/inabyte/embed)
![GitHub](https://img.shields.io/github/license/inabyte/embed)

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
