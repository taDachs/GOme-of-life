package gameoflife

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net"
  "net/http"
  "os"
)

type Client struct {
  UpdateUrl string
  InitUrl   string
  IP        string
  UdpPort   int
  UdpSocket *net.UDPConn
}

func (c *Client) SendChanges(changes []Change) {
  buf, err := json.Marshal(changes)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  req, err := http.NewRequest("POST", c.UpdateUrl, bytes.NewBuffer(buf))
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

func (c *Client) SendInit(init Init) {
  buf, err := json.Marshal(init)
  req, err := http.NewRequest("POST", c.InitUrl, bytes.NewBuffer(buf))
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

  fmt.Println("Response status code: ", resp.StatusCode)
}

func (c *Client) SendSync(sync Sync) {
  addr := net.UDPAddr{
    Port: c.UdpPort,
    IP:   net.ParseIP(c.IP),
  }

  buf, err := json.Marshal(sync.Board)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  _, err = c.UdpSocket.WriteToUDP(buf, &addr)
  if err != nil {
    println("Write failed:", err.Error())
    os.Exit(1)
  }
}
