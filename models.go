package main

type Category struct {
	Id   string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}

type Session struct {
	VoterId string
	Voters  []string
	Voted   int
	Choices []string
	Result  string
}
