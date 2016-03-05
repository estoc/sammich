package main

import (
	"testing"
)

func TestNew(t *testing.T) {
	var api = NewAPI(MockPlaceAPI{})
	room, err := api.New("address")
	if err != nil {
		t.Error(err)
	}
	if room.ID == "" {
		t.Error("ID should be populated")
	}
	if room.HostID == "" {
		t.Error("HostID should be populated")
	}
	if len(room.Choices) == 0 {
		t.Error("Choices should be populated")
	}
}

func TestGet(t *testing.T) {
	var api = NewAPI(MockPlaceAPI{})
	room, _ := api.New("address")
	room, err := api.Get(room.ID)
	if err != nil {
		t.Error(err)
	}
	if len(room.Choices) == 0 {
		t.Error("Choices should be populated")
	}
	if room.HostID != "" {
		t.Error("HostID should be empty")
	}
	if room.Votes != nil {
		t.Error("Votes should be empty")
	}
}

func TestGetNonexistent(t *testing.T) {
	var api = NewAPI(MockPlaceAPI{})
	_, err := api.Get("ID That Doesnt Exist")
	if err != ErrorRoomNotFound {
		t.Error("Should have gotten not found error")
	}
}

func TestVote(t *testing.T) {
	var api = NewAPI(MockPlaceAPI{})
	room, _ := api.New("address")
	voterName := "Name123"
	err := api.Vote(room.ID, voterName, room.Choices[0])
	if err != nil {
		t.Error(err)
	}

	room, err = api.Get(room.ID)
	if err != nil {
		t.Error(err)
	}
	if room.Voters[0] != voterName {
		t.Error("Voter list should contain name")
	}
}

func TestEnd(t *testing.T) {
	var api = NewAPI(MockPlaceAPI{})
	room, _ := api.New("address")
	hostID := room.HostID
	choice := room.Choices[0]
	api.Vote(room.ID, "votername", choice)

	err := api.End(room.ID, "Not Host ID")
	if err != ErrorUnauthorized {
		t.Error("Should get unauthorized error with wrong hostID")
	}

	err = api.End(room.ID, hostID)
	if err != nil {
		t.Error(err)
	}

	room, err = api.Get(room.ID)
	// MockAPI returns choice string as winner
	if room.Winner != choice {
		t.Error("Should have winner after end")
	}

	err = api.Vote(room.ID, "votername", "choice")
	if err != ErrorRoomEnded {
		t.Error("Should have gotten room ended error")
	}
}
