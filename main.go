package main

import (
  "flag"
)

func main() {
  // parse command line args
  assetsDir := flag.String("assetsDir", ".", "the absolute path to static assets. defaults to working directory.")
  flag.Parse()

  // initialize router
  router := NewMethodRouter()

  // provide static assets
  // served from root path
  router.ServeStatic(*assetsDir)

  // register api routes
  // served from "/api" path
  router.HandleFunc("GET", "/", HandleHelloWorld) // GET /api/

  // TODO: move to http.Server instantiation if we need TLS
	router.ListenAndServe(":8080")

  return
}