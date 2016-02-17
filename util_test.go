package main

import (
	"testing"
)

// Generate IDs and ensure they are unique
func TestGenereateId(t *testing.T) {
	// key value = number of occurences
	var m = make(map[string]int)
	i := 1000
	for i > 0 {
		id := generateID()
		m[id] = m[id] + 1
		i--
	}

	for _, v := range m {
		if v > 1 {
			t.Error("Generated a duplicate ID")
		}
	}
}
