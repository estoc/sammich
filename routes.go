package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/dancannon/gorethink"
	router "github.com/julienschmidt/httprouter"
)

func RtePing(w http.ResponseWriter, req *http.Request, p router.Params) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Pong")
}

func RteHelloWorld(w http.ResponseWriter, req *http.Request, p router.Params) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello %s! %s\n", p.ByName("name"), req.URL.Query().Get("test"))
}

func RteGetCategories(w http.ResponseWriter, req *http.Request, p router.Params) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Db("chewcrew").Table("categories").Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []Category
	rows.All(&results)
	j, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(j)
}
