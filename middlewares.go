package main

import (
  "net/http"
  "github.com/gorilla/feeds"
)

// Returns a middleware function that calls the next handler.
//
// This middleware function is executed for each incoming request and decorates the
// req/res combination with necessary facilities like unique ids and logging facilities.
func DecoratorMdw(next HttpHandler) HttpHandler {
  return func (w http.ResponseWriter, r *http.Request) {
    // attach a unique id to the request
    uuidv4 := feeds.NewUUID().String()
    r.Header.Set("X-Request-Id", uuidv4)
    w.Header().Set("X-Request-Id", uuidv4)

    // TODO: attach child logger to request context

    next(w, r)
  }
}