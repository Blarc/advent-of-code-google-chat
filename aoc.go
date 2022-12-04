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

const aocUrl = "https://adventofcode.com/2022/leaderboard/private/view/427349.json"

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

func googleChat(key string, token string, message string) {

	webhookUrl := fmt.Sprintf("https://chat.googleapis.com/v1/spaces/AAAAoYQLvw0/messages?key=%s&token=%s", key, token)

	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"text": message,
	})

	requestBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(
		webhookUrl,
		"application/json; charset=UTF-8",
		requestBody,
	)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Println(sb)
}

func aoc(sessionCookie string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", aocUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("cookie", "session="+sessionCookie)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
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

func createMessage(leaderboard Leaderboard) string {
	members := make([]Member, 0, len(leaderboard.Members))
	for _, member := range leaderboard.Members {
		members = append(members, member)
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].LocalScore > members[j].LocalScore
	})

	var result string = "🎄Leadboard update🎄\n```\n"
	for i, member := range members {

		var stars string
		for i := 0; i < len(member.Days); i++ {
			numberOfStars := len(member.Days[strconv.Itoa(i+1)])
			if numberOfStars == 2 {
				stars += string('★')
			} else {
				stars += string('☆')
			}
		}

		lastDay := member.Days[strconv.Itoa(len(member.Days))]

		lastStarTimestamp := lastDay[strconv.Itoa(len(lastDay))].GetStarTs
		loc, _ := time.LoadLocation("Europe/Ljubljana")
		lastStarDateTime := time.Unix(lastStarTimestamp, 0).
			In(loc).
			Format("(01-02 15:04)")

		result += fmt.Sprintf("%1d) %3d %-25s %13s %s\n", i+1, member.LocalScore, stars, lastStarDateTime, member.Name)

	}

	result += "```"
	return result
}

func compareLeaderboards(a Leaderboard, b Leaderboard) bool {
	// Same number of members
	if len(a.Members) == len(b.Members) {
		for memberKeyA := range a.Members {
			if memberB, ok := b.Members[memberKeyA]; ok {
				memberA := a.Members[memberKeyA]
				// Same number of stars
				if memberA.Stars != memberB.Stars {
					return false
				}
			}
		}
	} else {
		return false
	}
	return true
}

func main() {

	key := os.Getenv("googleChatKey")
	token := os.Getenv("googleChatToken")
	sessionCookie := os.Getenv("sessionCookie")

	newLeaderboardJson := aoc(sessionCookie)
	leaderboard := parseLeaderboard(newLeaderboardJson)
	savedLeaderboard := parseLeaderboard(aocSavedJson)

	if !compareLeaderboards(savedLeaderboard, leaderboard) {

		err := os.WriteFile("saved.json", []byte(newLeaderboardJson), 0644)
		if err != nil {
			log.Fatal(err)
		}

		message := createMessage(leaderboard)
		googleChat(key, token, message)
		log.Println(message)
	} else {
		log.Println("Leaderboard has not changed.")
	}

}
