package main

import (
	"io"
	"net/http"

	context "github.com/gorilla/context"
)

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	log := context.Get(r, "log").(*Logger)
	log.Debug("Hello, World!")

	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Hello, World!")
}
