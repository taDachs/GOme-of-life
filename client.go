package gameoflife

import (
  "encoding/json"
  "bytes"
  "fmt"
  "net/http"
)

func SendChanges(changes []Change, url string) {
  buf, err := json.Marshal(changes)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }
  defer resp.Body.Close()

  fmt.Println("Response status code:", resp.StatusCode)
}


func InitGame(url string) {
  req, err := http.NewRequest("Get", url, nil)
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }
  var init Init
  err = json.NewDecoder(resp.Body).Decode(&init)

  if err != nil {
    fmt.Println("Error decoding response: ", err)
    return
  }

  InitChannel <- init

  defer resp.Body.Close()

  fmt.Println("Response status code: ", resp.StatusCode)
}

func SyncGame(game *Game, url string) {
  var sync Sync
  sync.Board = *game.Board
  buf, err := json.Marshal(sync)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }
  defer resp.Body.Close()

  fmt.Println("Response status code:", resp.StatusCode)
}
