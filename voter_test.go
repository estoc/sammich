package main

import (
	"testing"
)

func TestAddVoters(t *testing.T) {
	var vs Voters
	vs.add("voter1")
	vs.add("voter2")
	vs.add("voter3")

	if len(vs) != 3 {
		t.Error("List should have multiple elements")
	}

	for _, v := range vs {
		if voter, _ := vs.getById(v.Id); voter == nil {
			t.Error("Error finding voter by ID")
		}
		if voter, _ := vs.getByName(v.Name); voter == nil {
			t.Error("Error finding voter by name")
		}
	}
}

func TestEveryoneReady(t *testing.T) {
	var vs Voters
	vs.add("voter1")
	vs.add("voter2")

	if vs.everyoneReady() == true {
		t.Error("EveryoneReady should be false")
	}

	vs[0].ready()
	vs[1].ready()

	if vs.everyoneReady() == false {
		t.Error("EveryoneReady should be true")
	}
}

func TestEveryoneVoted(t *testing.T) {
	var vs Voters
	vs.add("voter1")
	vs.add("voter2")

	if vs.everyoneVoted() == true {
		t.Error("EveryoneVoted should be false")
	}

	vs[0].vote()
	vs[1].vote()

	if vs.everyoneVoted() == false {
		t.Error("EveryoneVoted should be true")
	}
}
