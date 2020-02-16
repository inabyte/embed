package main

import (
	"net/http"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	go closeWhenUp()

	main()
}

func closeWhenUp() {
	for {
		time.Sleep(time.Microsecond * 100)

		if resp, err := http.Get("http://localhost:8080/"); err == nil {
			if resp.Body != nil {
				resp.Body.Close()
			}
			if resp.StatusCode == http.StatusOK {
				srv.Close()
			}
		}
	}
}
