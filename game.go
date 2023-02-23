package gameoflife

import (
  "fmt"
  "sync"
)

type Game struct {
  Board   *Board
  Started bool
  IsHost  bool
  Mutex   sync.Mutex
  Changes chan Change
  Syncs   chan Sync
  Inits   chan Init
  Client  Client
}

func (game *Game) TickCallback() {
  if game.IsHost && game.Board.Gen%SYNC_INTERVAL == 0 && game.Board.Gen > 0 {
    go game.Client.SyncGame(game)
    fmt.Println("Requesting Sync")
  }

  select {
  case chg, ok := <-game.Changes:
    if ok {
      game.Board.SetCell(chg.Alive, chg.X, chg.Y)
      fmt.Println("Setting cell: ", chg)
    }
  case sync, ok := <-game.Syncs:
    if ok && !game.IsHost {
      game.Board = &sync.Board
      fmt.Println("Syncing game")
    }
  case init, ok := <-game.Inits:
    if ok && !game.Started {
      if !game.IsHost {
        game.Board = &init.Board
      }
      game.Started = true
      fmt.Println("Init board")
    }
  default:
  }
  if game.Started {
    game.Board.NextGen()
  }
}
