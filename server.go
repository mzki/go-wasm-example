//go:build server
// +build server

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./html", "direcotry to be served")
	flag.Parse()

	fs := http.FileServer(http.Dir(dir))
	log.Print("Serving " + dir + " on http://localhost:8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		fs.ServeHTTP(resp, req)
	}))
}
