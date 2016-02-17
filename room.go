package main

import (
	"errors"
	"log"
	"sync"
)

type Rooms struct {
	Rooms map[string]*Room
	PlaceAPI
}

type Room struct {
	ID string `json:"id"`

	// The ID of the room creator.
	// Only returned in New() to remain secret.
	HostID string `json:"hostid,omitempty"`

	// List of choices and their total number of votes.
	Choices []string `json:"choices,omitempty"`

	// List of voters.
	Voters []string `json:"voters,omitempty"`

	// List of votes.
	// Seperate from Choices so the number of votes remains secret.
	Votes map[string]int32 `json:"-"`

	// The winning choice.
	Winner string `json:"winner,omitempty"`

	// Extra options for the Place API
	PlaceOptions

	sync.Mutex `json:"-"`
}

var (
	ErrorRoomNotFound = errors.New("Room not found")
	ErrorRoomEnded    = errors.New("Room has ended")
	ErrorUnauthorized = errors.New("Unauthorized host ID")
)

// NewRooms used to initialize the Rooms "API"
func NewRooms() Rooms {
	return Rooms{
		Rooms:    make(map[string]*Room),
		// TODO: Use GooglePlaceAPI if a valid client_id is provided
		PlaceAPI: MockPlaceAPI{},
	}
}

// Get a room!
func (r Rooms) Get(id string) (*Room, error) {
	room, ok := r.Rooms[id]
	if !ok {
		return nil, ErrorRoomNotFound
	}

	// Clear private fields
	room.HostID = ""
	room.Votes = nil
	return room, nil
}

// New creates a new room.
func (r Rooms) New(address string) (*Room, error) {
	log.Printf("NEW address=%s\n", address)

	room := Room{
		ID:           generateID(),
		HostID:       generateID(),
		Choices:      []string{},
		Votes:        make(map[string]int32),
		PlaceOptions: PlaceOptions{},
	}

	// Populate Choices and Votes
	cats := r.PlaceAPI.Categories()
	for _, v := range cats {
		room.Choices = append(room.Choices, string(v))
		room.Votes[string(v)] = 0
	}

	r.Rooms[room.ID] = &room
	return &room, nil
}

// Vote
func (r Rooms) Vote(id string, name string, vote string) error {
	log.Printf("VOTE id=%s name=%s\n", id, name)

	room, ok := r.Rooms[id]
	if !ok {
		return ErrorRoomNotFound
	}
	room.Lock()
	defer room.Unlock()

	// Skip if the room has already finished
	if room.Winner != "" {
		return ErrorRoomEnded
	}

	room.Voters = append(room.Voters, name)
	room.Votes[vote]++
	return nil
}

func (r Rooms) End(id string, hostid string) error {
	log.Printf("END id=%s hostid=%s\n", id, hostid)

	room, ok := r.Rooms[id]
	if !ok {
		return ErrorRoomNotFound
	}
	room.Lock()
	defer room.Unlock()

	// Verify the host ID
	if room.HostID != hostid {
		return ErrorUnauthorized
	}

	// Determine winning category
	max := int32(0)
	var winner Category
	for k, v := range room.Votes {
		if v > max {
			max = v
			winner = Category(k)
		}
	}

	// Find a place to eat!
	place := r.PlaceAPI.Get(room.PlaceOptions, winner)
	room.Winner = string(place)
	return nil
}
