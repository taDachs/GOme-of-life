package main

import (
  "flag"
  "fmt"
  "gameoflife"
  "os"
  "strconv"
  "time"
  "net"

  "github.com/faiface/pixel/pixelgl"
)

const width = 800
const height = 800

const res = 8
const board_width = width / res
const board_height = height / res

func main() {
  partner_ip := flag.String("ip", "127.0.0.1", "IP address to connect to")
  partner_port := flag.String("partner-port", "8081", "IP address to connect to")
  own_port := flag.String("own-port", "8080", "IP address to connect to")
  partner_udp_port_arg := flag.String("partner-udp-port", "1234", "IP address to connect to")
  own_udp_port_arg := flag.String("own-udp-port", "4321", "IP address to connect to")
  is_host_arg := flag.String("is-host", "false", "If host")
  flag.Parse()

  url := "http://" + *partner_ip + ":" + *partner_port + "/api"
  fmt.Println("own port: ", *own_port)
  fmt.Println("partner url: ", url)

  board := gameoflife.CreateEmptyBoard(board_width, board_height)

  own_udp_port, err := strconv.Atoi(*own_udp_port_arg)
  if err != nil {
    fmt.Println(*own_udp_port_arg, " is not a number")
    os.Exit(1)
  }
  addr := net.UDPAddr{
      Port: own_udp_port,
      IP: net.ParseIP("127.0.0.1"),
  }
  conn, err := net.ListenUDP("udp", &addr)
  if err != nil {
    println("ResolveUDPAddr failed:", err.Error())
    os.Exit(1)
  }


  partner_udp_port, err := strconv.Atoi(*partner_udp_port_arg)
  if err != nil {
    fmt.Println(*partner_udp_port_arg, " is not a number")
    os.Exit(1)
  }
  client := new(gameoflife.Client)
  client.InitUrl   = url + "/init"
  client.UpdateUrl = url + "/update"
  client.UdpPort   = partner_udp_port
  client.IP = *partner_ip
  client.UdpSocket = conn

  game := new(gameoflife.Game)
  game.IsHost    = *is_host_arg == "true"
  game.Board     = board
  game.Started   = false
  game.Changes   = make(chan gameoflife.Change, 10)
  game.Inits     = make(chan gameoflife.Init, 10)
  game.Syncs     = make(chan gameoflife.Sync, 10)
  game.Client    = *client

  game.GenFrequency    = 10
  game.UpdateFrequency = 100
  game.SyncFrequency   = 10


  if game.IsHost {
    gameoflife.InitSeed()
    board.InitializeRandom(0.2)
  }

  server := new(gameoflife.Server)
  server.Inits   = game.Inits
  server.Syncs   = game.Syncs
  server.Changes = game.Changes
  server.Port    = *own_port
  server.UdpSocket = conn

  go server.Run()
  go game.Run()

  // ensure server started
  time.Sleep(1 * time.Second)

  var init gameoflife.Init
  init.Board = *game.Board
  if game.IsHost {
    go client.SendInit(init)
    game.Inits <- init
  }

  var drawer gameoflife.BoardDrawer
  drawer.Game = game
  drawer.Width = width
  drawer.Height = height
  drawer.Res = res

  // run gui
  pixelgl.Run(drawer.Run)
}
