package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
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

// Logs incoming requests and outgoing responses
func MdwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		go log.WithFields(log.Fields{
			"Id":         req.Header.Get(ID_HEADER),
			"Host":       req.Host,
			"Method":     req.Method,
			"RequestURI": req.RequestURI}).Info("Request")

		start := time.Now()
		next.ServeHTTP(w, req)
		finish := time.Since(start)

		go log.WithFields(log.Fields{
			"Id":       req.Header.Get(ID_HEADER),
			"Duration": finish}).Info("Response")
	})
}
