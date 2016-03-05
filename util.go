// util.go contains self contained utility functions
package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

var (
	// chars used to generate a random ID
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	letterslen = len(letters)
)

// Generates a string of random characters
func generateID() string {
	b := make([]rune, 7)
	for i := range b {
		b[i] = letters[rand.Intn(letterslen)]
	}
	return string(b)
}

// Send JSON result/error response back to client
func sendJSON(w http.ResponseWriter, room *Room, err error) {
	// Enable CORS
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
