package sleahckSlack

import (
  "log"
  "net/http"
  "io/ioutil"
  "os"
  "encoding/json"
)

type UserList struct {
  Ok bool `json:"ok"`
  Members []struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Profile struct {
      Display_name string `json:"display_name"`
    } `json:"profile"`
  } `json:"members"`
}

func GetUserList() *UserList {
  var users *UserList

  reqURL := "https://slack.com/api/users.list"

  req, err := http.NewRequest("GET", reqURL, nil)
  if err != nil {
    log.Println(err)
  }

  req.Header.Add("Authorization", "Bearer " + os.Getenv("SLACK_TOKEN"))
  req.Header.Add("Content-Type", "application/json")

  client := new(http.Client)
  resp, err := client.Do(req)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  err = json.Unmarshal(body, &users)

  if err != nil {
    log.Println(err)
  }

  return users
}
