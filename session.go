package main

import (
	"errors"
	"strings"
)

type Session struct {
	ID string `json:"id,omitempty"`

	// The ID of the current voter
	VoterID string `json:"voterid,omitempty"`

	// [string] keys are IDs
	Voters  map[string]Voter  `json:"voters,omitempty"`
	Choices map[string]Choice `json:"choices,omitempty"`

	// The Host is the one who Creates/Starts the session
	// Omit to prevent host spoofing
	HostID string `json:"-"`
}

type Voter struct {
	Name  string `json:"name"`
	Voted bool   `json:"voted,omitempty"`

	// Omit to maintain vote anonymity
	Vote string `json:"-"`
}

type Choice struct {
	Place string `json:"place"`
	Votes int    `json:"votes,omitempty"`
}

var (
	sessions map[string]Session = make(map[string]Session)

	ErrorSessionNotFound = errors.New("Session not found")
	ErrorNameInUse       = errors.New("Name already used")
	ErrorVoterNotFound   = errors.New("Voter not found")
	ErrorUnauthorized    = errors.New("Must be host to perform this action")
)

// Create a session
func sessionCreate() (Session, error) {
	session := Session{
		ID:      generateId(),
		Voters:  make(map[string]Voter),
		Choices: make(map[string]Choice),
	}
	sessions[session.ID] = session
	return session, nil
}

// Retrieve a session
func sessionGet(id string) (Session, error) {
	session, ok := sessions[id]
	if !ok {
		return Session{}, ErrorSessionNotFound
	}

	if everyoneVoted(session) && votesAreUncounted(session) {
		sessions[id] = tallyVotes(session)
	}
	return session, nil
}

func everyoneVoted(session Session) bool {
	for _, v := range session.Voters {
		if v.Voted == false {
			return false
		}
	}
	return true
}

func votesAreUncounted(session Session) bool {
	for _, v := range session.Choices {
		if v.Votes > 0 {
			return false
		}
	}
	return true
}

func tallyVotes(session Session) Session {
	choices := session.Choices
	for k, v := range session.Voters {
		choice := choices[k]
		result := Choice{
			Place: choice.Place,
			Votes: choice.Votes + 1,
		}
		choices[v.Vote] = result
	}
	return session
}

// Joins a session and adds the new user
func sessionJoin(id string, name string) (Session, error) {
	session, ok := sessions[id]
	if !ok {
		return Session{}, ErrorSessionNotFound
	}

	// Ensure name is unique
	for _, v := range session.Voters {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return session, ErrorNameInUse
		}
	}

	voterid := generateId()
	voter := Voter{
		Name:  name,
		Voted: false,
	}
	session.Voters[voterid] = voter

	// First to join is host
	if session.HostID == "" {
		session.HostID = voterid
	}

	sessions[id] = session
	// Do this after updating the session object because this should be sent to the client, but not saved (since each user will have a different voterID)
	session.VoterID = voterid

	return session, nil
}

// Vote in a session
func sessionVote(id string, voterid string, choiceid string) (Session, error) {
	session, ok := sessions[id]
	if !ok {
		return Session{}, ErrorSessionNotFound
	}

	voter, ok := session.Voters[voterid]
	if !ok {
		return session, ErrorVoterNotFound
	}

	voter.Vote = choiceid
	voter.Voted = true
	session.Voters[voterid] = voter

	sessions[id] = session
	return session, nil
}

// Starts a session by populating its Choices
func sessionStart(id string, voterid string) (Session, error) {
	session, ok := sessions[id]
	if !ok {
		return Session{}, ErrorSessionNotFound
	}

	if voterid != session.HostID {
		return Session{}, ErrorUnauthorized
	}

	for _, place := range getPlaces(5) {
		choice := Choice{
			Place: place,
			Votes: 0,
		}
		session.Choices[generateId()] = choice
	}

	sessions[id] = session
	return session, nil
}

// Ends a session by forcing all voters to show as voted
// This will kick off TallyVotes the next time sessionGet() is called
func sessionEnd(id string, voterid string) (Session, error) {
	session, ok := sessions[id]
	if !ok {
		return Session{}, ErrorSessionNotFound
	}

	if voterid != session.HostID {
		return Session{}, ErrorUnauthorized
	}

	for k, v := range session.Voters {
		v.Voted = true
		session.Voters[k] = v
	}

	sessions[id] = session
	return session, nil
}
