/*
Package embed takes a list for file or folders (likely at `go generate` time) and
generates Go code that statically implements the a http.FileSystem.

Features:

-	Efficient generated code without unneccessary overhead.

-	Minimizes html css and js files.

-	Uses gzip compression internally (selectively, only for files that compress well).

-	Outputs `gofmt`ed Go code.
*/
package embed
