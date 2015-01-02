package main

import (
	"net/http"
	"strings"
)

func httpAddress() string {
	// Take the http addr from the command line if possible
	if len(flagHttpAddr) > 0 {
		return strings.TrimSpace(flagHttpAddr)
	}

	// We weren't handed a path on a plate, so have to synthesize it as best we can
	return ":8080"
}

func serveHttp(httpAddr string) {
	go http.ListenAndServe(httpAddr, nil)
}
