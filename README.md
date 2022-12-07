# Advent of Code Leaderboard
This script sends Advent of Code leaderboard information to a Google Chat space.

## How it works
The script saves the current leaderboard in a `saved.json` file on GitHub. Every 15 minutes, the script uses the Advent of Code API to retrieve the updated leaderboard. 
If the leaderboard has changed, the script saves the updated leaderboard to GitHub and sends a message with the updated leaderboard to the Google Chat space.

## Getting started
To use this script, click the `Use this template` button and set up the required GitHub environment variables. 
You can set up the environment variables by going to the  `⚙️ Settings > Secrets` section of your project repository.

## Running the script
Before running the script, you need to set the following environment variables:

```bash
export sessionCookie="YOUR_SESSION_COOKIE"
export googleChatUrl="YOUR_GOOGLE_CHAT_URL"
export leaderboardUrl="YOUR_LEADERBOARD_URL"
```

Once you have set the environment variables, you can run the script using the following command:

```bash
go run aoc.go
```
