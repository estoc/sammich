package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// A common type used to express what middleware and route handler function signatures look like
type HttpHandler func(w http.ResponseWriter, r *http.Request)

/*
The MethodRouter type provides a convenience wrapper around mux.Router that provides easy-to-use
http method-based routing.
*/
type MethodRouter struct {
	methods       []string
	primaryRouter *mux.Router
	subRouters    map[string]*mux.Router
}

// Get a new MethodRouter
func NewMethodRouter() *MethodRouter {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	router := &MethodRouter{methods, mux.NewRouter(), make(map[string]*mux.Router)}

	for _, method := range router.methods {
		// all api routes start with /api/
		router.subRouters[method] = router.primaryRouter.PathPrefix("/api/").Methods(method).Subrouter()
	}

	return router
}

/*
Register a route. See mux.Router documentation for regexp path rules, etc.

This method also wraps the provided handler with a decorator middleware and a middleware that
clears the request context.
*/
func (mr MethodRouter) HandleFunc(method string, path string, handleFunc HttpHandler) {
	serverLog.Info("Registering api route [%s /api%s]", method, path)
	mr.subRouters[method].HandleFunc(path, DecoratorMdw(CleanupMdw(handleFunc)))
	return
}

// Serve static content from an absolute path on the fs
func (mr MethodRouter) ServeStatic(dirPath string) {
	serverLog.Info("Serving statics assets from \"%s\"", dirPath)
	mr.primaryRouter.PathPrefix("/").Handler(http.FileServer(http.Dir(dirPath)))
}

// Start server.
func (mr MethodRouter) ListenAndServe(addr string) {
	log.Fatal(http.ListenAndServe(addr, mr.primaryRouter))
	return
}
