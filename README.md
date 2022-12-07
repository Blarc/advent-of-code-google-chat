# Advent of Code Leaderboard
A script that sends leaderboard information to Google Chat space.

## Checking for leaderboard updates
Current leaderboard is saved in github (`saved.json`). Every 15 minutes, a request is sent to Advent of Code API, via Github actions, to retrieve the new leaderboard. If the leaderboard has changed, the leaderboard is saved to github and message with the new leaderboard is sent to Google Chat space.

## Running
Before running, you need to set these environment variables:
```bash
export sessionCookie="YOUR_SESSION_COOKIE"
export googleChatUrl="YOUR_GOOGLE_CHAT_URL"
export leaderboardUrl="YOUR_LEADERBOARD_URL"
```

Then you can run the script:
```bash
go run aoc.go
```