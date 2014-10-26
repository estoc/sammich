package chewcrew

import (
	"github.com/gorilla/context"
	"io"
	"net/http"
)

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	log := context.Get(r, "log").(*Logger)
	log.Debug("Hello, World!")

	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Hello, World!")
}
