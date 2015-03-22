package main

import (
	"errors"
	"sync"
)

type Sessions map[string]*Session

type Session struct {
	Id string `json:"id"`

	// The ID of the current voter. Only returned when a user first creates/joins a session.
	VoterId string `json:"voterid,omitempty"`

	// List of voters. Populated as users join the session.
	Voters `json:"voters,omitempty"`

	// The current voter. Populated in get().
	CurrentVoter *Voter `json:"currentvoter,omitempty"`

	// List of choices. Populated in ready().
	Choices `json:"choices,omitempty"`

	// The winning choice. Populated in vote().
	Winner *Choice `json:"winner,omitempty"`

	sync.Mutex `json:"-"`
}

var (
	ErrorSessionNotFound   = errors.New("Session not found")
	ErrorSessionInProgress = errors.New("Session is in progress")
	ErrorSessionNotReady   = errors.New("Session is not ready for voting")
	ErrorSessionFinished   = errors.New("Session has ended")
)

// Get a session.
// VoterId is required to get details about the current voter.
// Returns a pointer to the true session (the one persisted in the database).
func (s Sessions) get(id string, voterid string) (*Session, error) {
	session, ok := s[id]
	if !ok {
		return nil, ErrorSessionNotFound
	}

	if voterid != "" {
		voter, err := session.Voters.getById(voterid)
		if err != nil {
			return nil, err
		}
		session.CurrentVoter = voter
	} else {
		session.CurrentVoter = nil
	}

	return session, nil
}

// Creates a new session and adds the first voter as host.
// Returns a session object containing only a SessionID and VoterID.
func (s Sessions) create(votername string) (*Session, error) {
	session := Session{
		Id:      generateId(),
		Voters:  Voters{},
		Choices: Choices{},
		Winner:  nil,
	}

	voter, _ := session.Voters.add(votername)
	s[session.Id] = &session

	// Returning only the essentials to keep the voterID private.
	result := Session{
		Id:      session.Id,
		VoterId: voter.Id,
	}
	return &result, nil
}

// Joins a session.
// Returns a session object containing ONLY a SessionID and VoterID.
func (s Sessions) join(id string, voterName string) (*Session, error) {
	session, ok := s[id]
	if !ok {
		return nil, ErrorSessionNotFound
	}
	session.Lock()
	defer session.Unlock()

	// Can't join a session that is already in progress
	if len(session.Choices) > 0 {
		return nil, ErrorSessionInProgress
	}

	voter, err := session.Voters.add(voterName)
	if err != nil {
		return nil, err
	}

	// Returning only the essentials to keep the voterID private.
	result := Session{
		Id:      session.Id,
		VoterId: voter.Id,
	}
	return &result, nil
}

// Sets a user as Ready to vote.
// Generates choices if everyone is ready.
// Returns only an error if encountered.
func (s Sessions) ready(id string, voterid string) error {
	session, ok := s[id]
	if !ok {
		return ErrorSessionNotFound
	}
	session.Lock()
	defer session.Unlock()

	// Skip if the session has already finished
	if session.Winner != nil {
		return ErrorSessionFinished
	}
	// Skip if session is already in progress
	if len(session.Choices) > 0 {
		return ErrorSessionInProgress
	}

	// Set voter to Ready
	v, err := session.Voters.getById(voterid)
	if err != nil {
		return err
	}
	v.ready()

	// Generate choices once everyone is ready
	// Also clear out the Ready flags
	if session.Voters.everyoneReady() {
		session.Choices.generate(3)
		session.Voters.clearReady()
	}

	return nil
}

// Lock in a vote.
// Tallys votes and determines winner if the last vote.
// Returns only an error if encountered.
func (s Sessions) vote(id string, voterId string, choiceId string) error {
	session, ok := s[id]
	if !ok {
		return ErrorSessionNotFound
	}
	session.Lock()
	defer session.Unlock()

	// Skip if the session has already finished
	if session.Winner != nil {
		return ErrorSessionFinished
	}
	// Skip if the session is not ready for votes
	if len(session.Choices) == 0 {
		return ErrorSessionNotReady
	}

	voter, err := session.Voters.getById(voterId)
	if err != nil {
		return err
	}
	choice, err := session.Choices.getById(choiceId)
	if err != nil {
		return err
	}

	err = voter.vote()
	if err != nil {
		// Error thrown if voter has already voted
		return err
	}
	choice.vote()

	// Tally votes if this was the last vote
	// Also clear out the Choices, and Voted flags
	if session.Voters.everyoneVoted() {
		winner := session.Choices.determineWinner()
		session.Winner = &winner
		session.Choices = nil
		session.Voters.clearVoted()
	}

	return nil
}
