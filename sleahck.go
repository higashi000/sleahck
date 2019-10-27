package main

import (
    "github.com/nlopes/slack"
    "fmt"
    "strconv"
    "time"
    "log"
    "os"
    "net/http"
    "io/ioutil"
    "encoding/json"
//    "reflect"
    )


// Data Struct --- {{{

type Channels struct {
  Ok bool `json:"ok"`
  Channels []struct {
    Id string `json:"id"`
    Name string `json:"name"`
  } `json:"channels"`
}

type UsersList struct {
  Ok bool `json:"ok"`
  Members []struct {
    Id string `json:"id"`
    Name string `json:"name"`
    DisplayName string `json:""`
  } `json:"members"`
}

type User struct {
  UserID string `json:"id"`
  Profile struct {
    DisplayName string `json:"display_name"`
    RealName string `json:"real_name"`
  } `json:"profile"`
}

type Users struct {
  Ok bool `json:"ok"`
  UserData []User `json:"members"`
}

type Message struct {
  Text string `json:"text"`
  User string `json:"user"`
  Time string `json:"time"`
}

type Messages struct {
  Msg [50]Message `json:"msg"`
}

// }}}

func main() {
}

func getChannels() *Channels {
  var channels Channels

  reqURL := "https://slack.com/api/users.conversations"

  req, err := http.NewRequest("GET", reqURL, nil)
  if err != nil {
    log.Println("error")
  }

  req.Header.Add("Authorization", "Bearer " + os.Getenv("SLACK_TOKEN"))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  client := &http.Client{}
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    log.Println("error")
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println("error")
  }

  err = json.Unmarshal(body, &channels)
  if err != nil {
    log.Println("error")
  }

  for _, e := range channels.Channels {
    fmt.Println(e)
  }

  return &channels
}

func getUserData() *Users {
  var users Users

  reqURL := "https://slack.com/api/users.list"

  req, err := http.NewRequest("GET", reqURL, nil)
  if err != nil {
    log.Println("error")
  }

  req.Header.Add("Authorization", "Bearer " + os.Getenv("SLACK_TOKEN"))
  req.Header.Add("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  err = json.Unmarshal(body, &users)

  if err != nil {
    log.Println("error")
  }

  return &users
}

func getHistory(channelID string, userData Users) *Messages {
  var messages Messages
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  slackLog, err := api.GetChannelHistory(channelID, slack.HistoryParameters{"0", "0", 50, false, false})

  if err != nil {
    log.Println("error")
  }

  for i := 49; i >= 0; i-- {
    messages.Msg[i].Text = slackLog.Messages[i].Msg.Text

    for _, e := range userData.UserData {
      if slackLog.Messages[i].Msg.User == e.UserID {
        if e.Profile.DisplayName != "" {
          messages.Msg[i].User = e.Profile.DisplayName
        } else {
          messages.Msg[i].User = e.Profile.RealName
        }
      }
    }
    tmpUnixTime, _ := strconv.Atoi(slackLog.Messages[i].Msg.Timestamp)
    messages.Msg[i].Time = (time.Unix(int64(tmpUnixTime), 0)).String()
    fmt.Println(messages.Msg[i])
  }

  return &messages
}
