package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	db "github.com/dancannon/gorethink"
	router "github.com/julienschmidt/httprouter"
	chain "github.com/yageek/shttap"
)

var session *db.Session

func main() {
	// parse command line args
	logLevel := flag.String("logLevel", "debug", "index of server logging level. levels: [debug, info, warning, error, fatal, panic]. defaults to debug.")
	port := flag.String("port", "8080", "desired port to listen on. defaults to 8080.")
	flag.Parse()

	configureLogger(*logLevel)
	session = connectToDatabase()
	server := buildServer()

	// start server
	log.Info("Starting server on port ", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server))
}

func configureLogger(level string) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	lvl, _ := log.ParseLevel(level)
	log.SetLevel(lvl)
}

func connectToDatabase() *db.Session {
	log.Info("Connecting to database...")
	session, err := db.Connect(db.ConnectOpts{
		Address:     "chewcrew.cc:28015",
		Database:    "chewcrew",
		MaxActive:   20,
		MaxIdle:     5,
		IdleTimeout: time.Minute * 5,
	})
	if err != nil {
		log.Fatal(NewMaskedError(err))
	}
	log.Info("Connected to database ", session)
	return session
}

func buildServer() *chain.Stack {
	// build router
	router := router.New()
	router.GET("/ping", RtePing)
	router.GET("/hello/:name", RteHelloWorld)
	router.GET("/preferences/categories", RteGetCategories)
	router.POST("/sessions", RteNotImplemented)
	router.GET("/sessions/:id", RteNotImplemented)
	router.PUT("/sessions/:id", RteNotImplemented)
	router.POST("/sessions/:id", RteNotImplemented)
	router.POST("/sessions/:id/start", RteNotImplemented)
	router.POST("/sessions/:id/end", RteNotImplemented)

	// chain middleware
	stack := chain.NewStack()
	stack.Use(MdwId, MdwLog)
	stack.UseRouter(router)
	return stack
}
