package main

import (
	"strconv"
	"testing"
)

var voterCount = 0 // used when generating voternames

func TestCreateSession(t *testing.T) {
	var s = make(Sessions)
	if len(s) != 0 {
		t.Error("Should have clean slate")
	}

	session, err := s.create(genVotername())
	if err != nil {
		t.Error("Error creating session: " + err.Error())
	}
	if session.Id == "" {
		t.Error("Session ID was not populated")
	}
	if len(s) != 1 {
		t.Error("Session was not created")
	}
}

// Ensure that multiple sessions can be created
func TestCreateMultipleSessions(t *testing.T) {
	var s = make(Sessions)
	s.create(genVotername())
	s.create(genVotername())
	s.create(genVotername())

	if len(s) != 3 {
		t.Error("Should have multiple sessions")
	}
}

func TestJoin(t *testing.T) {
	var s = make(Sessions)
	session, _ := s.create(genVotername())
	_, err := s.join(session.Id, genVotername())
	if err != nil {
		t.Error("Error joining: " + err.Error())
	}
	if session.VoterId == "" {
		t.Error("VoterID was not returned")
	}

	if len(s) != 1 {
		t.Error("Session should have been saved")
	}
}

func TestAllReadyShouldStartSession(t *testing.T) {
	var s = make(Sessions)
	// session should contain the voterID of the session creator
	session, _ := s.create(genVotername())
	voter2, _ := s.join(session.Id, genVotername())
	err := s.ready(session.Id, session.VoterId)
	if err != nil {
		t.Error("Error readying: " + err.Error())
	}
	err = s.ready(session.Id, voter2.VoterId)
	if err != nil {
		t.Error("Error readying: " + err.Error())
	}

	session, _ = s.get(session.Id, "")
	if len(session.Choices) == 0 {
		t.Error("Choices should be populated")
	}
}

func TestVoting(t *testing.T) {
	var s = make(Sessions)
	// session should contain the voterID of the session creator
	session, _ := s.create(genVotername())
	voter1Id := session.VoterId
	voter2, _ := s.join(session.Id, genVotername())
	s.ready(session.Id, voter1Id)
	s.ready(session.Id, voter2.VoterId)

	// Submit votes for choice[1]
	session, _ = s.get(session.Id, "")
	choice := session.Choices[1]
	if err := s.vote(session.Id, voter1Id, choice.Id); err != nil {
		t.Error("Error voting: " + err.Error())
	}
	if session.Choices[1].Votes != 1 {
		t.Error("Votes should have incremented")
	}
	if err := s.vote(session.Id, voter2.VoterId, choice.Id); err != nil {
		t.Error("Error voting: " + err.Error())
	}

	// Winner should be populated when all votes are in
	if session.Winner.Id != choice.Id {
		t.Error("Choice with most votes should have won")
	}
}

func TestSessionNotFound(t *testing.T) {
	var s = make(Sessions)
	if _, err := s.get("DoesntExist", "whatever"); err != ErrorSessionNotFound {
		t.Error("Should throw error when getting non-existant session")
	}
}

func TestGetSession(t *testing.T) {
	var s = make(Sessions)
	expected, _ := s.create(genVotername())
	result, err := s.get(expected.Id, expected.VoterId)
	if err != nil {
		t.Error("Error getting session: ", err)
	}
	if result == nil || result.Id != expected.Id {
		t.Error("Session should have been found")
	}
	_, err = result.Voters.getById(expected.VoterId)
	if err != nil {
		t.Error("Error retrieving the initial voter: ", err)
	}
}

func TestUniqueVoterNames(t *testing.T) {
	var s = make(Sessions)
	session, _ := s.create("john smith")
	session, _ = s.get(session.Id, "")

	if len(session.Voters) != 1 {
		t.Error("Should be 1 voter")
	}

	_, err := s.join(session.Id, "John Smith")
	if err != ErrorNameInUse {
		t.Error("Joining with duplicate username should fail")
	}

	_, err = s.join(session.Id, "JOHN SMITH")
	if err != ErrorNameInUse {
		t.Error("Joining with duplicate username should fail")
	}

	if len(session.Voters) != 1 {
		t.Error("Should still be 1 voter")
	}
}

// Helper function to generate unique voter names
func genVotername() string {
	voterCount++
	return "votername" + strconv.Itoa(voterCount)
}
