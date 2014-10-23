package main

import (
  "encoding/json"
  "net/http"
  "github.com/gorilla/feeds"
  "github.com/gorilla/context"
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

    // attach logger to request context
    reqLog, err := serverLog.Child(uuidv4)
    if err != nil {
      reqLog.Error("failed to attach request logger: %s", err)
      http.Error(w, "failed to attach request logger", http.StatusInternalServerError)
      return
    }
    context.Set(r, "log", reqLog) // pass dereferenced reqLog so that requests can't modify logger

    // log incoming request
    headers, err := json.MarshalIndent(r.Header, "", "\t")
    if err != nil {
      reqLog.Error("an error occurred when marshalling request headers: %s", err)
    }
    reqLog.Info("[%s] [%s %s]\n%s", r.Host, r.Method, r.RequestURI, headers)

    next(w, r)
  }
}

// Returns a middleware that cleans up that request context and wraps http.ResponseWriter so that
// we can log outgoing requests
func CleanupMdw(next HttpHandler) HttpHandler {
  return func (w http.ResponseWriter, r *http.Request) {
    log := context.Get(r, "log").(*logger)

    // wrap our response writer so that we can log when the request leaves the system
    next(&loggedResponseWriter{w: w, r: r, log: *log}, r)

    context.Clear(r)
  }
}

// Implements the http.ResponseWriter interface, providing us a way to log outgoing requests
type loggedResponseWriter struct {
  w    http.ResponseWriter
  r    *http.Request
  log  logger
  status int
}
func (w *loggedResponseWriter) Header() http.Header {
  return w.w.Header()
}
func (w *loggedResponseWriter) Write(d []byte) (int, error) {
  // if i never explicitly called WriteHeader, status code will be 200
  if w.status == 0 {
    w.WriteHeader(200)
  }

  // log outgoing request
  headers, err := json.MarshalIndent(w.w.Header(), "", "\t")
  if err != nil {
    w.log.Error("an error occurred when marshalling response headers: %s", err)
  }
  w.log.Info("[%s] [%s %s] [%v]\n%s", w.r.Host, w.r.Method, w.r.RequestURI, w.status, headers)

  return w.w.Write(d)
}
func (w *loggedResponseWriter) WriteHeader(status int) {
  w.status = status
  w.w.WriteHeader(status)
}