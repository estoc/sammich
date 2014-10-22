package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type MethodRouter struct {
  methods []string
  primaryRouter *mux.Router
  subRouters map[string]*mux.Router
}

func NewMethodRouter() *MethodRouter {
  methods := []string{"GET", "POST", "PUT", "DELETE"}
  router := &MethodRouter{methods, mux.NewRouter(), make(map[string]*mux.Router)}

  for _, method := range router.methods {
    router.subRouters[method] = router.primaryRouter.Methods(method).Subrouter()
  }

  return router
}

func (mr MethodRouter) HandleFunc(method string, path string, handleFunc HttpHandler) {
  mr.subRouters[method].HandleFunc(path, DecorateMdw(handleFunc))
  return
}

func (mr MethodRouter) ListenAndServe(addr string) {
  log.Fatal(http.ListenAndServe(addr, mr.primaryRouter))
  return
}