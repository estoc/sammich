package chewcrew

import (
	"flag"
)

func main() {
	// parse command line args
	assetsDir := flag.String("assetsDir", ".", "the absolute path to static assets. defaults to working directory.")
	logLevel := flag.Int("logLevel", 5, "index of server logging level. levels: [CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG]. defaults to 5 (DEBUG).")
	port := flag.String("port", "8080", "desired port to listen on. defaults to 8080.")
	flag.Parse()

	// start server logger
	level := *logLevel
	serverLog = NewServerLogger(level, logFormat)
	serverLog.Debug("Server logging configured.")

	// initialize router
	router := NewMethodRouter()

	// provide static assets
	// served from root path
	router.ServeStatic(*assetsDir)

	// register api routes
	// served from "/api" path
	router.HandleFunc("GET", "/", HandleHelloWorld) // GET /api/

	// TODO: move to http.Server instantiation if we need TLS
	p := ":" + *port
	serverLog.Info("Starting server on %s", p)
	router.ListenAndServe(p)

	return
}
