package main

import (
	"net/http"

	context "github.com/gorilla/context"
	feeds "github.com/gorilla/feeds"
	router "github.com/julienschmidt/httprouter"
)

// Returns a middleware function that calls the next handler.

// This middleware function is executed for each incoming request and decorates the
// req/res combination with necessary facilities like unique ids and logging facilities.
func DecoratorMdw(next router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, p router.Params) {
		// attach a unique id to the request
		uuidv4 := feeds.NewUUID().String()
		r.Header.Set("X-Request-Id", uuidv4)
		w.Header().Set("X-Request-Id", uuidv4)

		// attach logger to request context
		reqLog, err := ServerLog.Child(uuidv4)
		if err != nil {
			msg := "failed to attach request logger"
			reqLog.Error(NewMaskedErrorWithContext(err, msg))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		context.Set(r, "log", reqLog)

		// log incoming request
		reqLog.Info(r.Header, "[%s] [%s %s]", r.Host, r.Method, r.RequestURI)

		next(w, r, p)
	}
}

// Returns a middleware that cleans up that request context and wraps http.ResponseWriter so that
// we can log outgoing requests.
func CleanupMdw(next router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, p router.Params) {
		log := context.Get(r, "log").(*Logger)

		// wrap our response writer so that we can log when the request leaves the system
		next(&LoggedResponseWriter{w: w, r: r, log: log}, r, p)

		context.Clear(r)
	}
}

//Implements the http.ResponseWriter interface, providing us a way to log outgoing requests.
type LoggedResponseWriter struct {
	w      http.ResponseWriter
	r      *http.Request
	log    *Logger
	status int
}

func (w *LoggedResponseWriter) Write(d []byte) (int, error) {
	// if i never explicitly called WriteHeader, status code will be 200
	if w.status == 0 {
		w.WriteHeader(200)
	}

	// log outgoing request
	w.log.Info(w.w.Header(), "[%s] [%s %s] [%v]", w.r.Host, w.r.Method, w.r.RequestURI, w.status)

	return w.w.Write(d)
}

func (w *LoggedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.w.WriteHeader(status)
}

func (w *LoggedResponseWriter) Header() http.Header {
	return w.w.Header()
}
