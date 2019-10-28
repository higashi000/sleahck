package main

import (
    "github.com/nlopes/slack"
    "github.com/gin-gonic/gin"
//    "fmt"
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

type History struct {
  ChannelName string `json:"channeName"`
  ChannelID string `json:"channelID"`
  Msg []Message `json:"msg"`
}

type Histories struct {
  history []History `json:"history"`
}

// }}}

func main() {
  r := gin.Default()
  channel := getChannels()
  users := getUserData()
  var histories Histories
  histories = update(*channel, *users)

  sendChannelList(r, *channel)
  sendHistory(r, histories)

 go callUpdate(&histories, *channel, *users)
 r.Run()
}

func sendChannelList(r *gin.Engine, channel Channels) {
  r.GET("/sleahck/channelList", func(c *gin.Context) {
      c.JSON(200, channel.Channels)
      })
}

func sendHistory(r *gin.Engine, history Histories) {
  r.GET("/sleahck/GetHistory/:cN", func(c *gin.Context) {

      channel := c.Params.ByName("cN")
      var sendIndex int
      for i, e := range history.history {
        if e.ChannelName == channel {
          sendIndex = i
          break
        }
      }
      c.JSON(200, history.history[sendIndex])
      })
}

// Get Channels Data --- {{{
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

  return &channels
}

// }}}

// Get Users Data --- {{{
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
//}}}

// Channel History {{{
func getHistory(channelID string, userData Users) History {
  var messages History
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  slackLog, err := api.GetChannelHistory(channelID, slack.HistoryParameters{"0", "0", 50, false, false})
  messages.ChannelID = channelID

  if err != nil {
    log.Println("error")
  }

  for i := 0; i < len(slackLog.Messages); i++ {
    tmpUnixTime, _ := strconv.Atoi(slackLog.Messages[i].Msg.Timestamp)
    messages.Msg = append(messages.Msg, Message{slackLog.Messages[i].Msg.Text, "", (time.Unix(int64(tmpUnixTime), 0)).String()})

    for _, e := range userData.UserData {
      if slackLog.Messages[i].Msg.User == e.UserID {
        if e.Profile.DisplayName != "" {
          messages.Msg[i].User = e.Profile.DisplayName
        } else {
          messages.Msg[i].User = e.Profile.RealName
        }
      }
    }
  }

  return messages
}
// }}}

// update --- {{{
func update(channel Channels, users Users) Histories {
  var histories Histories
  for i, e := range channel.Channels {
    tmp := getHistory(e.Id, users)
    histories.history = append(histories.history, tmp)
    histories.history[i].ChannelName = e.Name
  }

  return histories
}
// }}}

func callUpdate(histories *Histories, channel Channels, users Users) {
  for {
    t := time.Now()
    if t.Second() == 0 {
      tmp := update(channel, users)

      for i, _ := range channel.Channels {
        for j, f := range tmp.history[i].Msg {
          (*histories).history[i].Msg[j].Text = f.Text
          (*histories).history[i].Msg[j].Time = f.Time
          (*histories).history[i].Msg[j].User = f.User
        }

        histories.history[i].ChannelName = channel.Channels[i].Name
        histories.history[i].ChannelID = channel.Channels[i].Id
      }
    }
  }
}
