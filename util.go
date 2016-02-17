package main

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
)

// Generates a string of random characters
func generateID() string {
	b := make([]rune, 6)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
