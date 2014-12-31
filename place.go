package main

import (
	"math/rand"
	"time"
)

var (
	places = [...]string{
		"Boracho",
		"Chick-fil-A",
		"Union Park",
		"Woolsworth",
		"Wing Bucket",
		"Which Wich",
		"Jason's Deli",
		"Moe's",
		"Noodle Nexus",
		"Iron Cactus",
	}
)

func getPlaces(num int) []string {
	rand.Seed(time.Now().UnixNano())
	result := []string{}
	for num > 0 {
		place := places[rand.Intn(len(places))]

		// Ensure no duplicates
		for _, v := range result {
			if v == place {
				place = ""
				break
			}
		}

		// Duplicate found so try again
		if place == "" {
			continue
		}

		result = append(result, place)
		num--
	}
	return result
}
