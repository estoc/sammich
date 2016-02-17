package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var (
	// Program arguments
	port   = flag.String("port", "8080", "Port to run on")
	apikey = flag.String("apikey", "", "Places API Key")

	// The "database"
	rooms = NewRooms()
)

func main() {
	// Parse program arguments
	flag.Parse()

	// Initialize PlaceAPI
	if *apikey != "" {
		rooms.PlaceAPI = GooglePlaceAPI{*apikey}
	}

	// Set api endpoints
	http.HandleFunc("/room", Get)
	http.HandleFunc("/room/new", New)
	http.HandleFunc("/room/vote", Vote)
	http.HandleFunc("/room/end", End)

	// Serve static files for default endpoint
	http.Handle("/", http.FileServer(http.Dir("/web")))

	// Start server
	log.Println("Server running on port " + *port)
	log.Fatal(http.ListenAndServe(":"+(*port), nil))
}

// Get Session Handler
func Get(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	id := qp.Get("id")
	room, err := rooms.Get(id)
	sendResult(w, room, err)
}

// New Session Handler
func New(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := r.URL.Query()
	address := qp.Get("address")

	room, err := rooms.New(address)
	sendResult(w, room, err)
}

// Vote Session Handler
func Vote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := r.URL.Query()
	id := qp.Get("id")
	name := qp.Get("name")
	vote := qp.Get("vote")

	err := rooms.Vote(id, name, vote)
	sendResult(w, nil, err)
}

// End Session Handler
func End(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := r.URL.Query()
	id := qp.Get("id")
	hostid := qp.Get("hostid")

	err := rooms.End(id, hostid)
	sendResult(w, nil, err)
}

// Send result/error response back to client
func sendResult(w http.ResponseWriter, room *Room, err error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/json")
	if err != nil {
		// Send Error JSON result
		e := map[string]string{"error": err.Error()}
		result, _ := json.Marshal(e)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(result)
	} else if room != nil {
		// Send Room result
		result, _ := json.Marshal(*room)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		// Send blank result
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}
}
