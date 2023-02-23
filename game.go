package gameoflife

import (
  "fmt"
  "sync"
  "time"
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

const GEN_FREQUENCY = 20

func (game *Game) Run() {
  // background game update loop
  update_ticker := time.NewTicker(50 * time.Millisecond)
  quit := make(chan struct{})
  go func() {
    for {
      select {
      case <-update_ticker.C:
        game.UpdateTickCallback()
      case <-quit:
        update_ticker.Stop()
        return
      }
    }
  }()

  next_gen_ticker := time.NewTicker(1000 / GEN_FREQUENCY * time.Millisecond)
  go func() {
    for {
      select {
      case <-next_gen_ticker.C:
        game.NextGenTickCallback()
      case <-quit:
        next_gen_ticker.Stop()
        return
      }
    }
  }()

}

func (game *Game) UpdateTickCallback() {
  game.Mutex.Lock()

  if game.IsHost && game.Board.Gen%SYNC_INTERVAL == 0 && game.Board.Gen > 0 {
    game.Client.SyncGame(game)
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

  defer game.Mutex.Unlock()
}

func (game *Game) NextGenTickCallback() {
  game.Mutex.Lock()
  if game.Started {
    game.Board.NextGen()
  }
  defer game.Mutex.Unlock()
}
