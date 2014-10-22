package main

import (
  "io"
  "net/http"
)

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  io.WriteString(w, "Hello, World!")
}