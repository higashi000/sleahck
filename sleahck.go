package main

import (
    "github.com/nlopes/slack"
    "fmt"
    "strconv"
    "time"
    "log"
    "os"
//    "reflect"
    )

type Message struct {
  Text string `json:"text"`
  User string `json:"user"`
  Time string `json:"time"`
}

type Messages struct {
  msg [50]Message
}

func main() {
  channelID := ""

  getHistory(channelID)
}

func getHistory(channelID string) *Messages {
  var messages Messages
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  slackLog, err := api.GetChannelHistory(channelID, slack.HistoryParameters{"0", "0", 50, false, false})

  if err != nil {
    log.Println("error")
  }

  for i := 49; i >= 0; i-- {
    messages.msg[i].Text = slackLog.Messages[i].Msg.Text
    messages.msg[i].User = slackLog.Messages[i].Msg.User
    tmpUnixTime, _ := strconv.Atoi(slackLog.Messages[i].Msg.Timestamp)
    messages.msg[i].Time = (time.Unix(int64(tmpUnixTime), 0)).String()
    fmt.Println(messages.msg[i])
  }


  return &messages
}
