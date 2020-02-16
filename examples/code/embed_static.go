// +build ignore

package main

import (
	"github.com/inabyte/embed"
)

func main() {
	config := embed.New()

	config.Output = "static"
	config.Package = "main"
	config.Files = []string{"../../testdata"}
	config.FileServer = true
	config.NoLocalFS = true

	config.Generate()
}
