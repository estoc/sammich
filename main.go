package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.POST("/sessions", create)
	router.GET("/sessions/:id", get)
	router.POST("/sessions/:id", join)
	router.PUT("/sessions/:id", vote)
	router.POST("/sessions/:id/start", start)
	router.POST("/sessions/:id/end", end)

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("create")
	s, e := sessionCreate()
	sendResult(w, s, e)
}

func get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Printf("get: id=%s\n", id)
	s, e := sessionGet(id)
	sendResult(w, s, e)
}

func join(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	name := getQ(r, "name")
	log.Printf("join: id=%s, nae=%s\n", id, name)
	s, e := sessionJoin(id, name)
	sendResult(w, s, e)
}

func vote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	choice := getQ(r, "choiceid")
	log.Printf("vote: id=%s, voterid=%s\n", id, voterid)
	s, e := sessionVote(id, voterid, choice)
	sendResult(w, s, e)
}

func start(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	log.Printf("start: id=%s, voterid=%s\n", id, voterid)
	s, e := sessionStart(id, voterid)
	sendResult(w, s, e)
}

func end(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	log.Printf("end: id=%s, voterid=%s\n", id, voterid)
	s, e := sessionEnd(id, voterid)
	sendResult(w, s, e)
}

// Get query string parameter
func getQ(r *http.Request, key string) string {
	return r.URL.Query()[key][0]
}

// Send result/error response
func sendResult(w http.ResponseWriter, session Session, err error) {
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		result, _ := json.Marshal(session)
		w.Write(result)
	}
}
