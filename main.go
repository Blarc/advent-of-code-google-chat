package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

//go:embed saved.json
var aocSavedJson string

type Leaderboard struct {
	Event   string            `json:"event"`
	Members map[string]Member `json:"members"`
	OwnerID int               `json:"owner_id"`
}
type Part struct {
	GetStarTs int64 `json:"get_star_ts"`
	StarIndex int   `json:"star_index"`
}

type Member struct {
	GlobalScore int                        `json:"global_score"`
	ID          int                        `json:"id"`
	Name        string                     `json:"name"`
	Days        map[string]map[string]Part `json:"completion_day_level"`
	LocalScore  int                        `json:"local_score"`
	Stars       int                        `json:"stars"`
	LastStarTs  int                        `json:"last_star_ts"`
}

func sendMessageToGoogleChat(webhookUrl string, message string) {

	postBody, _ := json.Marshal(map[string]string{
		"text": message,
	})

	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(
		webhookUrl,
		"application/json; charset=UTF-8",
		requestBody,
	)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}

func getLeaderboard(leaderboardUrl string, sessionCookie string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", leaderboardUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("cookie", "session="+sessionCookie)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func parseLeaderboard(leaderboardJson string) Leaderboard {

	var leaderboard Leaderboard
	err := json.Unmarshal([]byte(leaderboardJson), &leaderboard)
	if err != nil {
		log.Fatal(err)
	}
	return leaderboard
}

func createMessage(leaderboard Leaderboard, changedMembers []Member) string {
	members := make([]Member, 0, len(leaderboard.Members))
	for _, member := range leaderboard.Members {
		members = append(members, member)
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].LocalScore > members[j].LocalScore
	})

	var result string = "🎄Leadboard update🎄\n"
	result += "⭐New stars: "

	for i, member := range changedMembers {
		if i == len(changedMembers)-1 {
			result += member.Name
		} else {
			result += member.Name
			result += ", "
		}
	}

	result += "\n```\n"
	result += "#  Last star     Username        Points Stars\n"

	loc, _ := time.LoadLocation("Europe/Ljubljana")

	// Figure out how many AoC puzzles is currently available
	today := time.Now().In(loc)
	numberOfDays := 25
	if today.Month() == time.December {
		if today.Day() < 26 {
			numberOfDays = today.Day()
		}
	}

	// Find the length of the longest username
	maxNameLength := 0
	for _, member := range members {
		memberNameLenght := len(member.Name)
		if memberNameLenght > maxNameLength {
			maxNameLength = memberNameLenght
		}
	}

	for i, member := range members {

		var stars string
		var lastDayNumber int
		for day := 1; day <= numberOfDays; day++ {

			dayString := strconv.Itoa(day)

			numberOfStars := len(member.Days[dayString])
			if numberOfStars == 2 {
				stars += string('★')
			} else if numberOfStars == 1 {
				stars += string('✮')
			} else {
				stars += string('☆')
			}

			if _, ok := member.Days[dayString]; ok && day > lastDayNumber {
				lastDayNumber = day
			}

		}

		lastDay := member.Days[strconv.Itoa(lastDayNumber)]
		lastStarTimestamp := lastDay[strconv.Itoa(len(lastDay))].GetStarTs
		lastStarDateTime := time.Unix(lastStarTimestamp, 0).
			In(loc).
			Format("(01-02 15:04)")

		result += fmt.Sprintf("%1d) %s %-*s %4d %s\n", i+1, lastStarDateTime, maxNameLength, member.Name, member.LocalScore, stars)

	}

	result += "```"
	return result
}

func compareLeaderboards(a Leaderboard, b Leaderboard) []Member {
	changedMembers := []Member{}
	// Same number of members
	if len(a.Members) == len(b.Members) {
		for memberKeyA := range a.Members {
			if memberB, ok := b.Members[memberKeyA]; ok {
				memberA := a.Members[memberKeyA]
				// Same number of stars
				if memberA.Stars != memberB.Stars {
					changedMembers = append(changedMembers, memberA)
				}
			}
		}
	}
	return changedMembers
}

func main() {

	leaderboardUrl := os.Getenv("leaderboardUrl")
	sessionCookie := os.Getenv("sessionCookie")
	googleChatUrl := os.Getenv("googleChatUrl")

	if len(leaderboardUrl) == 0 {
		panic("Environment variable \"leaderboardUrl\" is not set!")
	}

	if len(sessionCookie) == 0 {
		panic("Environment variable \"sessionCookie\" is not set!")
	}

	if len(googleChatUrl) == 0 {
		panic("Environment variable \"googleChatUrl\" is not set!")
	}

	newLeaderboardJson := getLeaderboard(leaderboardUrl, sessionCookie)
	leaderboard := parseLeaderboard(newLeaderboardJson)
	savedLeaderboard := parseLeaderboard(aocSavedJson)

	changedMembers := compareLeaderboards(savedLeaderboard, leaderboard)
	if len(changedMembers) > 0 {

		err := os.WriteFile("saved.json", []byte(newLeaderboardJson), 0644)
		if err != nil {
			log.Fatal(err)
		}

		message := createMessage(leaderboard, changedMembers)
		sendMessageToGoogleChat(googleChatUrl, message)
		log.Println(message)
	} else {
		log.Println("Leaderboard has not changed.")
	}

}
