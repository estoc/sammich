package main

import (
  "net/http"
  "github.com/gorilla/feeds"
)

func DecorateMdw(h HttpHandler) HttpHandler {
  return func (w http.ResponseWriter, r *http.Request) {
    // attach an id to the request
    uuidv4 := feeds.NewUUID().String()
    r.Header.Set("X-Request-Id", uuidv4)
    w.Header().Set("X-Request-Id", uuidv4)

    h(w, r)
  }
}