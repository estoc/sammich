package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// Parse program arguments
	port := flag.String("port", "8080", "Port to run on.")
	apikey := flag.String("apikey", "", "Places API Key. Mock API will be used if none provided.")
	flag.Parse()

	// Seed RNG for generateID()
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialize PlaceAPI, use Mock by default
	placeAPI := PlaceAPI(MockPlaceAPI{})
	if *apikey != "" {
		// Credentials provided, use GooglePlaceAPI
		placeAPI = GooglePlaceAPI{*apikey}
	}

	// Initialize the ChewCrew API
	api := NewAPI(placeAPI)

	// Set api endpoints
	http.HandleFunc("/room", api.GetHandler)
	http.HandleFunc("/room/new", api.NewHandler)
	http.HandleFunc("/room/vote", api.VoteHandler)
	http.HandleFunc("/room/end", api.EndHandler)

	// Serve static files for default endpoint
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Start server
	log.Println("Server running on port " + *port)
	log.Fatal(http.ListenAndServe(":"+(*port), nil))
}
