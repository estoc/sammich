package main

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_=+")
)

// Generates a string of random characters
func generateId() string {
	b := make([]rune, 33)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
