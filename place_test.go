package main

import (
	"testing"
)

// Should return a list of unique places
func TestGetPlaces(t *testing.T) {
	places := getPlaces(8)

	if len(places) <= 1 {
		t.Error("List should have multiple elements")
	}

	for i := range places {
		for j := range places {
			if i != j && places[i] == places[j] {
				t.Error("List should contain unique elements")
			}
		}
	}
}
