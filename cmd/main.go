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

  game := new(gameoflife.Game)
  game.IsHost = is_host
  game.Board = board
  game.Started = false

  if is_host {
    gameoflife.InitSeed()
    board.InitializeRandom(0.2)
  } else {
    go func() {
      init := <- gameoflife.InitChannel
      game.Board = &init.Board
      game.Started = true
      fmt.Println("Updated board")
    }()
  }

  go gameoflife.RunServer(*own_port, game)

  // ensure server started
  time.Sleep(1 * time.Second)

  if !is_host {
    gameoflife.InitGame(init_url)
  }

  // i dont know if you do it like that bu GO sounds good
  pixelgl.Run(func() {gameoflife.Run(update_url, sync_url, game, width, height, res)})
}
