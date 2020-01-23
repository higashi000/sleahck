package sleahckSlack

import (
  "log"
  "net/http"
  "io/ioutil"
  "os"
  "encoding/json"
)

type Channels struct {
  Ok bool `json:"ok"`
  Channel []struct {
    Id string `json:"id"`
    Name string `json:"name"`
  } `json:"Channels"`
}

func GetChannels() *Channels {
  var channels *Channels

  reqURL := "https://slack.com/api/users.conversations"

  req, err := http.NewRequest("GET", reqURL, nil)
  if err != nil {
    log.Println(err)
  }

  req.Header.Add("Authorization", "Bearer " + os.Getenv("SLACK_TOKEN"))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  client := new(http.Client)

  resp, err := client.Do(req)
  if err != nil {
    log.Println(err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println(err)
  }

  err = json.Unmarshal(body, &channels)
  if err != nil {
    log.Println(err)
  }

  return channels
}
