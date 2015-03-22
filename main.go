package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

var (
	port = flag.String("port", "8080", "Port to run on")

	// This is the "database"
	sessions = make(Sessions)

	// Used to hold the reference api doc
	api []byte
)

func main() {
	flag.Parse()

	router := httprouter.New()
	router.GET("/api", getApiRef)
	router.GET("/sessions/:id", get)
	router.POST("/sessions", create)
	router.POST("/sessions/:id/join", join)
	router.POST("/sessions/:id/vote", vote)
	router.POST("/sessions/:id/ready", ready)

	log.Println("Server running on port " + *port)
	log.Fatal(http.ListenAndServe(":"+(*port), router))
}

func getApiRef(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if api == nil {
		api = generateApiRef()
	}
	w.Write(api)
}

func get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	log.Printf("get: id=%s, voterid=%s\n", id, voterid)
	s, e := sessions.get(id, voterid)
	sendResult(w, s, e)
}

func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := getQ(r, "name")
	log.Printf("create: name=%s", name)
	s, e := sessions.create(name)
	sendResult(w, s, e)
}

func join(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	name := getQ(r, "name")
	log.Printf("join: id=%s, name=%s\n", id, name)
	s, e := sessions.join(id, name)
	sendResult(w, s, e)
}

func vote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("vote(): r=" + r.RequestURI)
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	choiceid := getQ(r, "choiceid")
	log.Printf("vote: id=%s, voterid=%s, choiceid=%s\n", id, voterid, choiceid)
	e := sessions.vote(id, voterid, choiceid)
	sendResult(w, nil, e)
}

func ready(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	voterid := getQ(r, "voterid")
	log.Printf("ready: id=%s, voterid=%s\n", id, voterid)
	e := sessions.ready(id, voterid)
	sendResult(w, nil, e)
}

// Get query string parameter
// Returns "" if not found
func getQ(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// Send result/error response back to client
func sendResult(w http.ResponseWriter, session *Session, err error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		e := map[string]string{"error": err.Error()}
		result, _ := json.Marshal(e)
		w.Write(result)
	} else if session != nil {
		result, _ := json.Marshal(*session)
		w.Write(result)
	} else {
		result, _ := json.Marshal(map[string]string{})
		w.Write(result)
	}
}

// Generates the API quick reference.
// Also acts as an integration test.
func generateApiRef() []byte {
	log.Println("Generating API Reference")
	var site = "http://localhost:" + *port + "/sessions"

	var result bytes.Buffer
	result.WriteString("# ChewCrew API Reference\n")
	result.WriteString("Get returns a full session object.\n")
	result.WriteString("Create/Join return only sessionId and voterId.\n")
	result.WriteString("Ready/Vote return an empty JSON object\n")
	result.WriteString("All return an error if encountered.\n")
	result.WriteString("\n\n")

	result.WriteString("## Creating and joining a session\n\n")

	// Create
	url := site + "?name=Voter1"
	session := executeApi(&result, url, "POST", "Start/Create as Voter1. Notice that VoterID is returned.")
	sessionId := session.Id
	voterid1 := session.VoterId

	// Get
	url = site + "/" + sessionId + "?voterid=" + voterid1
	executeApi(&result, url, "GET", "Get. Notice that VoterID is NOT returned.")

	// Join
	url = site + "/" + sessionId + "/join?name=Voter2"
	session = executeApi(&result, url, "POST", "Join as Voter2. Notice that VoterID is returned.")
	voterid2 := session.VoterId

	// Get
	url = site + "/" + sessionId + "?voterid=" + voterid2
	executeApi(&result, url, "GET", "Get.")

	// Get (without optional param voterId)
	url = site + "/" + sessionId
	executeApi(&result, url, "GET", "Get. Notice that a VoterID was not passed, so a CurrentVoter is not returned.")

	result.WriteString("## Readying and generating Choices\n\n")

	// Ready Voter1
	url = site + "/" + sessionId + "/ready?voterid=" + voterid1
	executeApi(&result, url, "POST", "Ready Voter1.")

	// Get
	url = site + "/" + sessionId + "?voterid=" + voterid1
	executeApi(&result, url, "GET", "Get. Notice the Ready flag set.")

	// Ready Voter2
	url = site + "/" + sessionId + "/ready?voterid=" + voterid2
	executeApi(&result, url, "POST", "Ready Voter2. Choices are generated after all voters are Ready.")

	// Get (now with Choices populated)
	url = site + "/" + sessionId
	session = executeApi(&result, url, "GET", "Get. Notice that Choices are now populated, and Ready flags are cleared.")
	choiceId := session.Choices[0].Id

	result.WriteString("## Voting and generating Winner\n\n")

	// Vote Voter1
	url = site + "/" + sessionId + "/vote?voterid=" + voterid1 + "&choiceid=" + choiceId
	executeApi(&result, url, "POST", "Vote Voter1.")

	// Get
	url = site + "/" + sessionId + "?voterid=" + voterid1
	executeApi(&result, url, "GET", "Get. Notice the Voted flag set.")

	// Vote Voter2
	url = site + "/" + sessionId + "/vote?voterid=" + voterid2 + "&choiceid=" + choiceId
	executeApi(&result, url, "POST", "Vote Voter2. Winner is generated after all votes are in.")

	// Get (now with Winner populated)
	url = site + "/" + sessionId
	executeApi(&result, url, "GET", "Final Get. Notice that Winner is now populated, and Voted flags are cleared.")

	result.WriteString("## Error Examples\n\n")

	// Get a non-existent session (should throw error)
	url = site + "/ugh"
	executeApi(&result, url, "GET", "[ERROR] Try to get a non-existent session.")

	// Ready Voter1 again (should throw error)
	url = site + "/" + sessionId + "/ready?voterid=" + voterid1
	executeApi(&result, url, "POST", "[ERROR] Try to ready in a finished session.")

	// Vote Voter2 (should throw error)
	url = site + "/" + sessionId + "/vote?voterid=" + voterid2 + "&choiceid=" + choiceId
	executeApi(&result, url, "POST", "[ERROR] Try to vote in a finished session.")

	return result.Bytes()
}

// Make an http request to the running server.
// Return the unmarshalled json Session object.
func executeApi(b *bytes.Buffer, url string, method string, comment string) Session {
	b.WriteString("[" + method + "] " + comment + "\n")
	b.WriteString(url + "\n")
	req, _ := http.NewRequest(method, url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	respBytes := readerToBytes(resp.Body)
	b.Write(respBytes)
	b.WriteString("\n\n\n")
	session := Session{}
	json.Unmarshal(respBytes, &session)
	return session
}

func readerToBytes(reader io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.Bytes()
}
