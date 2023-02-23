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
  update_url := url + "/update"
  sync_url := url + "/sync"
  init_url := url + "/init"
  is_host := *is_host_arg == "true"
  fmt.Println("own port: ", *own_port)
  fmt.Println("partner url: ", url)

  board := gameoflife.CreateEmptyBoard(board_width, board_height)

  client := new(gameoflife.Client)
  client.InitUrl   = init_url
  client.UpdateUrl = update_url
  client.SyncUrl   = sync_url

  game := new(gameoflife.Game)
  game.IsHost = is_host
  game.Board = board
  game.Started = false
  game.Changes= make(chan gameoflife.Change, 10)
  game.Inits = make(chan gameoflife.Init, 10)
  game.Syncs = make(chan gameoflife.Sync, 10)
  game.Client = *client


  if is_host {
    gameoflife.InitSeed()
    board.InitializeRandom(0.2)
  }

  server := new(gameoflife.Server)
  server.Inits   = game.Inits
  server.Syncs   = game.Syncs
  server.Changes = game.Changes

  go server.Run(*own_port)

  // ensure server started
  time.Sleep(1 * time.Second)

  // background run game loop
  ticker := time.NewTicker(50 * time.Millisecond)
  quit := make(chan struct{})
  go func() {
    for {
      select {
      case <- ticker.C:
        game.TickCallback()
      case <- quit:
        ticker.Stop()
        return
      }
    }
  }()

  if is_host {
    go client.InitGame(game)
  }

  // i dont know if you do it like that bu GO sounds good
  pixelgl.Run(func() {gameoflife.Run(game, width, height, res)})
}
