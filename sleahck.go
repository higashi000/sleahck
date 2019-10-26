package main

import (
    "github.com/nlopes/slack"
    "fmt"
    "os"
    )

func main() {
  api := slack.New(os.Getenv("SLACK_TOKEN"))
  channelID := "channelID"
  slackLog, err := api.GetChannelHistory(channelID, slack.HistoryParameters{"0", "0", 2, false, false})

  if err != nil {
    fmt.Println("error")
    return
  }

  fmt.Println(slackLog)
}
