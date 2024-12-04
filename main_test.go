package main

import (
	"os"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	buffer, err := os.ReadFile("saved.json")
	if err != nil {
		panic(err)
	}

	leaderboard := parseLeaderboard(string(buffer))
	message := createMessage(leaderboard, []Member{})
	println(message)

}
