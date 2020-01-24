// Copyright 2020 Inabyte Inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	syslog "log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/inabyte/embed"
)

type logger interface {
	Print(v ...interface{})
}

var (
	f             = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	log    logger = syslog.New(os.Stderr, "", syslog.LstdFlags)
	myExit        = os.Exit
)

// main generate the code
func main() {
	conf := embed.New()

	f.Usage = func() {
		fmt.Fprintf(f.Output(), `Usage:  %s [<options>] <files>
Where: <files> list of files and/or folders to embed
       <options> one or more of the following
`, os.Args[0])
		f.PrintDefaults()
	}

	f.StringVar(&conf.Output, "o", conf.Output, "Output files base.")
	f.StringVar(&conf.Package, "pkg", conf.Package, "Package name.")
	f.StringVar(&conf.BuildTags, "tags", conf.BuildTags, "Build tags.")
	f.StringVar(&conf.Ignore, "ignore", conf.Ignore, "Regexp for files we should ignore (for example \\\\.DS_Store).")
	f.StringVar(&conf.Include, "include", conf.Include, "Regexp for files to include. Only files that match will be included.")
	f.StringVar(&conf.Minify, "minify", conf.Minify, "Comma list of mimetypes to minify")
	f.StringVar(&conf.ModifyTime, "modifytime", conf.ModifyTime, "Unix timestamp to override as modification time for all files.")
	f.BoolVar(&conf.DisableCompression, "no-compress", conf.DisableCompression, "If true, do not compress files.")
	f.BoolVar(&conf.Go, "go", conf.Go, "write only go files")
	f.BoolVar(&conf.Binary, "binary", conf.Binary, "produce self-contained extracter/http server binary (<output> will become the binary name)")
	f.BoolVar(&conf.Local, "local", conf.Local, "If true, do not reference external files.")
	f.Parse(os.Args[1:])

	conf.Files = f.Args()

	if len(f.Args()) < 1 {
		showError(f, "No files/folders specified")
	} else {
		var err error

		if conf.Binary {
			err = createBinary(conf)
		} else {
			err = conf.Generate()
		}

		if err != nil {
			showError(f, err)
		}
	}
}

func showError(f *flag.FlagSet, v ...interface{}) {
	log.Print(v...)
	f.Usage()
	myExit(2)
}

func createBinary(conf *embed.Config) error {
	var (
		err     error
		outfile string
		base    string
	)

	outfile = conf.Output
	if !filepath.IsAbs(outfile) {
		outfile, err = os.Getwd()
		if err == nil {
			outfile = filepath.Join(outfile, conf.Output)
		}
	}

	if err == nil {
		conf.Local = true
		conf.DisableCompression = false
		conf.Package = "main"

		base, err = ioutil.TempDir("", "embed-generate")
	}

	if err == nil {
		defer os.RemoveAll(base)

		conf.Output = filepath.Join(base, "main")
		err = conf.Generate()
	}

	if err == nil {
		cmd := exec.Command("go", "build", "-o", outfile, "-ldflags", "-s")
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		cmd.Dir = base
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	return err
}
