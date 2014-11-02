package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	uuid "github.com/satori/go.uuid"
)

const ID_HEADER = "X-Request-Id"

// Decorates request/response with unique ID
func MdwId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := uuid.NewV4().String()
		req.Header.Set(ID_HEADER, id)
		next.ServeHTTP(w, req)
	})
}

// Logs incoming requests and outgoing responses. Also provides a child logger to each request's context.
func MdwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := log.WithFields(log.Fields{"Id": req.Header.Get(ID_HEADER)})
		context.Set(req, l, "log")

		go l.WithFields(log.Fields{
			"Host":       req.Host,
			"Method":     req.Method,
			"RequestURI": req.RequestURI}).Info("Request")

		start := time.Now()
		next.ServeHTTP(w, req)
		finish := time.Since(start)

		go l.WithFields(log.Fields{"Duration": finish}).Info("Response")
	})
}
