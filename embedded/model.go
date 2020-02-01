// Package embedded embed files
//
// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package embedded

import (
	"net/http"
	"os"
	"path/filepath"
)

// FileSystem defines the FileSystem interface and builder
type FileSystem interface {
	http.FileSystem

	// Walk all folders and files in the filesystem
	Walk(root string, walkFn filepath.WalkFunc) error

	// Copy all files to target directory
	Copy(target string, mode os.FileMode) error

	// AddFile add a file to embedded filesystem
	AddFile(path string, name string, local string, size int64, modtime int64, mimeType string, tag string, compressed bool, data []byte, str string)
	// AddFolder add a file to embedded filesystem
	AddFolder(path string, name string, local string, modtime int64, paths ...string)

	// UseLocal use on disk copy instead of embedded data (for development)
	UseLocal(bool)
}

// FileInfo internal file info and file contents
type FileInfo interface {
	os.FileInfo
	Compressed() bool // Is this file compressed
	Tag() string      // Etag for the file contents
	MimeType() string // file contens mimetype
	String() string   // file contents as string
	Bytes() []byte    // file contents as byte array
}

// Handler serves embedded handle to serve FileSystem
type Handler interface {
	http.Handler
	// SetNotFoundHandler set a hander to be called for no found
	SetNotFoundHandler(http.Handler)
	// If true and the folder does not contain index.html render folder
	// otherwise return 403 http.StatusForbidden
	SetRenderFolders(enable bool)
}
