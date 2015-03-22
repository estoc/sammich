package main

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	var c Choices
	c.generate(3)

	if len(c) != 3 {
		t.Error("List should have multiple elements")
	}

	for i := range c {
		for j := range c {
			if i != j && c[i] == c[j] {
				t.Error("List should contain unique elements")
			}
		}
	}
}

func TestVote(t *testing.T) {
	var c Choices
	c.generate(2)

	if c[0].Votes != 0 {
		t.Error("Vote should start at 0")
	}

	c[0].vote()
	c[0].vote()

	if c[0].Votes != 2 {
		t.Error("Vote should have incremented")
	}
	if c[1].Votes != 0 {
		t.Error("Vote should have stayed 0")
	}
}

func TestDetermineWinner(t *testing.T) {
	var c Choices
	c.generate(3)
	c[1].vote()
	winner := c.determineWinner()

	if winner != *c[1] {
		t.Error("Winner should be one with most votes")
	}
}
