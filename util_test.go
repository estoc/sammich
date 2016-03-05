package main

import (
	"testing"
)

// Generate IDs and ensure they are unique
func TestGenereateId(t *testing.T) {
	// Set of generated IDs
	var m = make(map[string]bool)

	for i := 10000; i > 0; i-- {
		id := generateID()

		// Ensure unique ID
		if _, found := m[id]; found == true {
			t.Error("Generated a duplicate ID")
			return
		}

		// Add to ID set
		m[id] = true
	}
}
