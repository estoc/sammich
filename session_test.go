package main

import (
	"testing"
)

func TestHappyPath(t *testing.T) {
	sessions = make(map[string]Session)
	if len(sessions) != 0 {
		t.Error("Should have clean slate")
	}

	// Create
	session, err := sessionCreate()
	if err != nil {
		t.Error("Error creating: " + err.Error())
	}
	if &session == nil {
		t.Error("Session was not returned")
	}

	// Join (first to join is host)
	hostName := "hostname"
	session, err = sessionJoin(session.ID, hostName)
	hostid := session.VoterID
	sessionid := session.ID
	if err != nil {
		t.Error("Error joining as host: " + err.Error())
	}
	if hostid == "" {
		t.Error("VoterID was not returned")
	}
	if len(session.Voters) != 1 {
		t.Error("Host voter was not added")
	}
	if session.HostID != hostid {
		t.Error("HostID should be the first voter to join")
	}
	if session.Voters[hostid].Name != hostName {
		t.Error("Host name does not match")
	}

	// Join again (regular non-host voter)
	voterName := "votername"
	session, err = sessionJoin(sessionid, voterName)
	voterid := session.VoterID
	if err != nil {
		t.Error("Error joining as non-host: " + err.Error())
	}
	if voterid == "" || voterid == hostid {
		t.Error("VoterID was not returned")
	}
	if len(session.Voters) != 2 {
		t.Error("Non-host voter was not added")
	}
	if session.HostID == "" || session.HostID != hostid || session.HostID == voterid {
		t.Error("HostID should still be populated with original hostID")
	}
	if session.Voters[voterid].Name != voterName {
		t.Error("Non-host name does not match")
	}

	// Start session as non-host (should fail)
	session, err = sessionStart(sessionid, voterid)
	if err == nil {
		t.Error("Error should be thrown when non-host tries to start session")
	}
	if len(session.Choices) > 0 {
		t.Error("Choices should not be populated")
	}

	// Start session as host
	session, err = sessionStart(sessionid, hostid)
	if err != nil {
		t.Error("Error starting session: " + err.Error())
	}
	if len(session.Choices) == 0 {
		t.Error("Choices should be populated")
	}

	// Get choice IDs
	choiceIds := []string{}
	for k, _ := range session.Choices {
		choiceIds = append(choiceIds, k)
	}

	// Submit vote for choice 0 as host
	session, err = sessionVote(sessionid, hostid, choiceIds[0])
	if err != nil {
		t.Error("Error voting as host: " + err.Error())
	}
	voter := session.Voters[hostid]
	if voter.Vote != choiceIds[0] || voter.Voted == false {
		t.Error("Host vote should be populated")
	}

	// Submit vote for choice 1 as non-host
	session, err = sessionVote(sessionid, voterid, choiceIds[1])
	if err != nil {
		t.Error("Error voting as non-host: " + err.Error())
	}
	voter = session.Voters[voterid]
	if voter.Vote != choiceIds[1] || voter.Voted == false {
		t.Error("Non-host vote should be populated")
	}

	// Get the session, which should tally the votes
	session, err = sessionGet(sessionid)
	if err != nil {
		t.Error("Error getting finished session: " + err.Error())
	}
	if session.Choices[choiceIds[0]].Votes != 1 || session.Choices[choiceIds[1]].Votes != 1 {
		t.Error("Votes were not populated")
	}
	if session.Choices[choiceIds[2]].Votes != 0 {
		t.Error("Vote was populated where it didn't belong")
	}

	// Changing a vote after votes are tallied should not affect previous tally
	sessionVote(sessionid, hostid, choiceIds[1])
	session, err = sessionGet(sessionid)
	if session.Choices[choiceIds[0]].Votes != 1 || session.Choices[choiceIds[1]].Votes != 1 {
		t.Error("Tallied votes should have not changed")
	}
}

// Ensure votes are tallied when SessionEnd() is ran
func TestSessionEnd(t *testing.T) {
	session, err := sessionCreate()
	sessionid := session.ID

	session, _ = sessionJoin(sessionid, "host")
	hostid := session.VoterID

	session, _ = sessionJoin(sessionid, "dude")
	voterid := session.VoterID

	sessionStart(sessionid, hostid)

	choiceIds := []string{}
	for k, _ := range session.Choices {
		choiceIds = append(choiceIds, k)
	}

	sessionVote(sessionid, hostid, choiceIds[0])

	// End session as non-host (should fail)
	session, err = sessionEnd(sessionid, voterid)
	if err == nil {
		t.Error("Ending session as non-host should fail")
	}

	// End session as host
	session, err = sessionEnd(sessionid, hostid)
	if err != nil {
		t.Error("Error ending session: " + err.Error())
	}
	voters := session.Voters
	if voters[hostid].Voted == false || voters[voterid].Voted == false {
		t.Error("Voted flags should be set to true")
	}
	if voters[hostid].Vote != choiceIds[0] || voters[voterid].Vote != "" {
		t.Error("Votes should be unchanged")
	}

	// Get session which should tally votes
	session, err = sessionGet(sessionid)
	if (session.Choices[choiceIds[0]].Votes != 1) && (session.Choices[choiceIds[1]].Votes != 0) {
		t.Error("Votes are incorrect")
	}
}

// Ensure that multiple sessions can be created
func TestCreateMultipleSessions(t *testing.T) {
	sessions = make(map[string]Session)
	if len(sessions) != 0 {
		t.Error("Should have clean slate")
	}

	sessionCreate()
	sessionCreate()
	sessionCreate()

	if len(sessions) <= 1 {
		t.Error("Should have multiple sessions")
	}
}

// Call functions using non-existent session/voter IDs
func TestUsingBadIds(t *testing.T) {
	sessions = make(map[string]Session)
	if len(sessions) != 0 {
		t.Error("Should have clean slate")
	}

	// Join non-existent session (should fail)
	session, err := sessionJoin("whatever", "whatever")
	if err == nil {
		t.Error("Joining non-existent session should throw error")
	}
	if len(sessions) != 0 {
		t.Error("No session should have been created")
	}

	// Create a session
	session, err = sessionCreate()
	sessionid := session.ID

	_, err = sessionStart(sessionid, "whatever")
	if err == nil {
		t.Error("Starting as non-existent user should throw error")
	}

	_, err = sessionEnd(sessionid, "whatever")
	if err == nil {
		t.Error("Ending as non-existent user should throw error")
	}
}

// Ensure voter names are unique
func TestUniqueVoterNames(t *testing.T) {
	sessions = make(map[string]Session)
	if len(sessions) != 0 {
		t.Error("Should have clean slate")
	}

	session, err := sessionCreate()
	sessionid := session.ID
	sessionJoin(sessionid, "voter1")
	sessionJoin(sessionid, "voter2")
	sessionJoin(sessionid, "voter3")

	session, _ = sessionGet(sessionid)
	if len(session.Voters) != 3 {
		t.Error("Should be 3 voters")
	}

	session, err = sessionJoin(sessionid, "voter1")
	if err == nil || len(session.Voters) != 3 {
		t.Error("Joining with duplicate username should fail")
	}

	session, err = sessionJoin(sessionid, "VOTER2")
	if err == nil || len(session.Voters) != 3 {
		t.Error("Joining with duplicate username should fail")
	}

	session, err = sessionJoin(sessionid, "VoTeR3")
	if err == nil || len(session.Voters) != 3 {
		t.Error("Joining with duplicate username should fail")
	}
}
