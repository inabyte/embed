package main

//go:generate go run ../../cmd/embed/cmd.go -pkg=main -nolocalfs -fileserver -o static ../../testdata

import (
	"fmt"
	"log"
	"net/http"
)

var srv *http.Server

func main() {
	fmt.Println("Open link http://localhost:8080/")
	// FileHandler() is created by embed and returns a http.Handler.
	srv = &http.Server{Addr: ":8080", Handler: FileHandler()}
	log.Print(srv.ListenAndServe())
}
