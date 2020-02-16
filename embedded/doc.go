/*
Package embedded implements the a http.FileSystem and http.Handler.


FileSystem

Implement a the embedded FileSystem, and provides the following methods.

	AddFile(path string, name string, local string, size int64, modtime int64, mimeType string, tag string, compressed bool, data []byte, str string) error
Add a file to embedded filesystem, the file must not exist.

	AddFolder(path string, name string, local string, modtime int64, paths ...string) error
Add an directory to the embedded filesystem, the folder must not exist.

	http.FileSystem
		Open(name string) (http.File, error)
Implement a http.FileSystem Open method returning a http.File.
If the embedded content is being served the the file file will
also implement the FileInfo and http.Handler interfaces.

	Walk(root string, walkFn WalkFunc) error
Walk walks the file tree rooted at root, calling walkFn for each file or
directory in the tree, including root. All errors that arise visiting files
and directories are filtered by walkFn. The files are walked in lexical
order.

	Copy(target string, mode os.FileMode) error
Copy or extract all files to target directory, creates file with the specfified mode

	WriteFile(filename string, data []byte, perm os.FileMode) error
Writes data to a file named by filename.
If the file does not exist, WriteFile creates it with permissions perm;
otherwise WriteFile truncates it before writing.

	UseLocal(bool)
Use on disk copy instead of embedded data intended for development.
When developing javascript or html you want the origion source served (uncompressed and not minified)
this will only work with file that have a local path specified. Open will return a os.File using the local path.


Handler

Implement a handler the serve data in response to http requests.

	http.Handler
		ServeHTTP(http.ResponseWriter, *http.Request)
Serves data in response to a http request, this will serve gzip files to clients that accept compressed content.
also serve Etag-based conditional requests.


	SetNotFoundHandler(http.Handler)
Set a hander to be called when request asset cannot be found.

	SetPermissionHandler(http.Handler)
Set a hander to be called for no permission http.StatusForbidden


	SetRenderFolders(enable bool)
If true and the folder does not contain index.html render folder
otherwise the PermissionHandler will be called usually serving 403 http.StatusForbidden

FileInfo

Internal file info and access file contents

	os.FileInfo
Implements the os.FileInfo methods.

	Compressed() bool
Is this file compressed

	Tag() string
Etag for the file contents, used to conditionally serve the content based on If-None-Match header.

	MimeType() string
Mimetype for file contents when this file is served the mime type header will be filled in with this value.

	String() string
file contents as string, if not compress this will just be the internal string,
if it is compress it will uncompress it.

	Bytes() []byte
file contents as byte array uncompressing if necessary, this is always a copy of the content event if it not compressed

	Raw() []byte
raw bytes (if compressed the compressed bytes) this is in readonly memory
any attemt to write to this will result in a segmentation fault
}

*/
package embedded
