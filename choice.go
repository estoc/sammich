package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	mockPlaces = [...]string{
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
		"Enchillada",
		"Original Italian",
	}

	ErrorChoiceNotFound = errors.New("Choice not found")
)

type Choices []*Choice

type Choice struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	// Omit to maintain vote anonymity
	Votes int `json:"-"`
}

func (c *Choices) getById(id string) (*Choice, error) {
	for _, v := range *c {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, ErrorChoiceNotFound
}

func (c *Choices) determineWinner() Choice {
	winner := (*c)[0]
	for _, v := range *c {
		if v.Votes > winner.Votes {
			winner = v
		}
	}
	return *winner
}

// Generates and adds choices
func (c *Choices) generate(num int) {
	rand.Seed(time.Now().UnixNano())
	for num > 0 {
		name := mockPlaces[rand.Intn(len(mockPlaces))]
		if c.placeAlreadyUsed(name) {
			continue
		}

		place := Choice{
			Id:    generateId(),
			Name:  name,
			Votes: 0,
		}
		*c = append(*c, &place)
		num--
	}
}

func (c *Choices) placeAlreadyUsed(newplace string) bool {
	for _, v := range *c {
		if v.Name == newplace {
			return true
		}
	}
	return false
}

func (c *Choice) vote() {
	c.Votes++
}

func (c Choice) String() string {
	return fmt.Sprintf("Choice: Id=%s, Name=%s, Votes=%s", c.Id, c.Name, c.Votes)
}
