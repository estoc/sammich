// api.go contains the main ChewCrew API functionality
package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

type API struct {
	// Rooms is the list of current rooms/sessions
	// Key = Room ID
	Rooms map[string]*Room
	PlaceAPI
}

// Room is the room/session returned to the client
type Room struct {
	// ID of the Room
	ID string `json:"id"`

	// The ID of the room creatoapi.
	// Only returned in New() to remain secret.
	HostID string `json:"hostid,omitempty"`

	// List of voters.
	Voters []string `json:"voters,omitempty"`

	// List of choices and their total number of votes.
	Choices []string `json:"choices,omitempty"`

	// List of votes.
	// Seperate from Choices so the number of votes remains secret.
	Votes map[string]int32 `json:"-"`

	// The winning choice.
	Winner string `json:"winner,omitempty"`

	// Options for the Place API.
	// TODO: Set in New().
	PlaceOptions `json:"-"`

	// Mutex used to ensure syncronization
	sync.Mutex `json:"-"`
}

var (
	ErrorRoomNotFound = errors.New("Room not found")
	ErrorRoomEnded    = errors.New("Room has ended")
	ErrorUnauthorized = errors.New("Unauthorized host ID")
)

// NewAPI used to initialize a new API
func NewAPI(placeAPI PlaceAPI) *API {
	return &API{
		Rooms:    make(map[string]*Room),
		PlaceAPI: placeAPI,
	}
}

// Get Session Handler
func (api *API) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	room, err := api.Get(id)
	sendJSON(res, room, err)
}

// New Session Handler
func (api *API) NewHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := req.URL.Query()
	address := qp.Get("address")

	room, err := api.New(address)
	sendJSON(res, room, err)
}

// Vote Session Handler
func (api *API) VoteHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := req.URL.Query()
	id := qp.Get("id")
	name := qp.Get("name")
	vote := qp.Get("vote")

	err := api.Vote(id, name, vote)
	sendJSON(res, nil, err)
}

// End Session Handler
func (api *API) EndHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Invalid request method", http.StatusBadRequest)
		return
	}

	qp := req.URL.Query()
	id := qp.Get("id")
	hostid := qp.Get("hostid")

	err := api.End(id, hostid)
	sendJSON(res, nil, err)
}

// Get a room!
func (api *API) Get(id string) (*Room, error) {
	room, ok := api.Rooms[id]
	if !ok {
		return nil, ErrorRoomNotFound
	}

	// Clear private fields
	room.HostID = ""
	room.Votes = nil
	return room, nil
}

// New creates a new room.
// The only method that returns HostID to keep it secret.
func (api *API) New(address string) (*Room, error) {
	log.Printf("NEW address=%s\n", address)

	// Create new rom
	room := Room{
		ID:           generateID(),
		HostID:       generateID(),
		Choices:      []string{},
		Votes:        make(map[string]int32),
		PlaceOptions: PlaceOptions{},
	}

	// Populate Choices and Votes
	cats := api.PlaceAPI.Categories()
	for _, v := range cats {
		room.Choices = append(room.Choices, string(v))
		room.Votes[string(v)] = 0
	}

	// Add room
	api.Rooms[room.ID] = &room
	return &room, nil
}

// Vote adds a new voter name and their vote to a room
func (api *API) Vote(id string, name string, vote string) error {
	log.Printf("VOTE id=%s name=%s\n", id, name)

	// Get room
	room, ok := api.Rooms[id]
	if !ok {
		return ErrorRoomNotFound
	}
	room.Lock()
	defer room.Unlock()

	// Skip if room has already ended
	if room.Winner != "" {
		return ErrorRoomEnded
	}

	// Add voter and vote
	room.Voters = append(room.Voters, name)
	room.Votes[vote]++
	return nil
}

// End ends a voting session
// Tally votes and deterimine winning place
// Can only be used by the Host
func (api *API) End(id string, hostid string) error {
	log.Printf("END id=%s hostid=%s\n", id, hostid)

	// Get room
	room, ok := api.Rooms[id]
	if !ok {
		return ErrorRoomNotFound
	}
	room.Lock()
	defer room.Unlock()

	// Verify host ID
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
	// TODO: If a place is not found, try 2nd, 3rd, etc. winning category
	place, _ := api.PlaceAPI.Get(room.PlaceOptions, winner)
	room.Winner = string(place)
	return nil
}
