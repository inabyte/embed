package main

//go:generate go run embed_static.go

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
