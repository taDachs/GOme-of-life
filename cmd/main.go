package main

import (
  "flag"
  "fmt"
  "gameoflife"
  "time"

  "github.com/faiface/pixel/pixelgl"
)

const width = 800
const height = 800

const res = 8
const board_width = width / res
const board_height = height / res

func main() {
  arg := flag.String("ip", "0.0.0.0", "IP address to connect to")
  partner_port := flag.String("partner-port", "8081", "IP address to connect to")
  own_port := flag.String("own-port", "8080", "IP address to connect to")
  is_host_arg := flag.String("is-host", "false", "If host")
  flag.Parse()

  url := "http://" + *arg + ":" + *partner_port + "/api"
  fmt.Println("own port: ", *own_port)
  fmt.Println("partner url: ", url)

  board := gameoflife.CreateEmptyBoard(board_width, board_height)

  client := new(gameoflife.Client)
  client.InitUrl   = url + "/init"
  client.UpdateUrl = url + "/update"
  client.SyncUrl   = url + "/sync"

  game := new(gameoflife.Game)
  game.IsHost = *is_host_arg == "true"
  game.Board = board
  game.Started = false
  game.Changes= make(chan gameoflife.Change, 10)
  game.Inits = make(chan gameoflife.Init, 10)
  game.Syncs = make(chan gameoflife.Sync, 10)
  game.Client = *client


  if game.IsHost {
    gameoflife.InitSeed()
    board.InitializeRandom(0.2)
  }

  server := new(gameoflife.Server)
  server.Inits   = game.Inits
  server.Syncs   = game.Syncs
  server.Changes = game.Changes


  go server.Run(*own_port)
  go game.Run()

  // ensure server started
  time.Sleep(1 * time.Second)

  if game.IsHost {
    go client.InitGame(game)
  }

  pixelgl.Run(func() {gameoflife.Run(game, width, height, res)})
}
