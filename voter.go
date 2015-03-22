package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorNameInUse         = errors.New("Name already used")
	ErrorVoterNotFound     = errors.New("Voter not found")
	ErrorVoterAlreadyVoted = errors.New("Voter already voted")
)

type Voters []*Voter

type Voter struct {
	Name  string `json:"name"`
	Voted bool   `json:"voted,omitempty"`
	Ready bool   `json:"ready,omitempty"`

	// Omit to prevent user spoofing
	// Users get their VoterID through the Session object, upon creating/joining the session
	Id string `json:"-"`
}

func (vs *Voters) add(name string) (*Voter, error) {
	// Ensure name is unique
	if vs.voterNameUsed(name) {
		return nil, ErrorNameInUse
	}

	voter := Voter{
		Id:    generateId(),
		Name:  name,
		Ready: false,
		Voted: false,
	}
	*vs = append(*vs, &voter)
	return &voter, nil
}

func (vs *Voters) voterNameUsed(name string) bool {
	for _, v := range *vs {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func (vs *Voters) getByName(name string) (*Voter, error) {
	for _, v := range *vs {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return v, nil
		}
	}
	return nil, ErrorVoterNotFound
}

func (vs *Voters) getById(id string) (*Voter, error) {
	for _, v := range *vs {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, ErrorVoterNotFound
}

func (vs *Voters) everyoneVoted() bool {
	for _, v := range *vs {
		if v.Voted == false {
			return false
		}
	}
	return true
}

func (vs *Voters) everyoneReady() bool {
	for _, v := range *vs {
		if v.Ready == false {
			return false
		}
	}
	return true
}

// Clears flag to save some bytes
func (vs *Voters) clearVoted() {
	for _, v := range *vs {
		v.Voted = false
	}
}

// Clears flag to save some bytes
func (vs *Voters) clearReady() {
	for _, v := range *vs {
		v.Ready = false
	}
}

func (v *Voter) vote() error {
	if v.Voted == true {
		return ErrorVoterAlreadyVoted
	}
	v.Voted = true
	return nil
}

func (v *Voter) ready() {
	v.Ready = true
}

func (v Voter) String() string {
	return fmt.Sprintf("Voter: Id=%s, Name=%s, Ready=%s, Voted=%s", v.Id, v.Name, v.Ready, v.Voted)
}
