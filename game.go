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

  GenFrequency    float64 // in hertz
  UpdateFrequency float64
  SyncFrequency   float64
}

func (game *Game) Run() {
  quit := make(chan struct{})
  next_gen_ticker := time.NewTicker(time.Duration(1000/game.GenFrequency) * time.Millisecond)
  sync_ticker := time.NewTicker(time.Duration(1000/game.SyncFrequency) * time.Millisecond)
  go func() {
    for {
      game.UpdateTickCallback()
      select {
      case <-next_gen_ticker.C:
        game.NextGenTickCallback()
      case <-sync_ticker.C:
        game.PerformSync()
      case <-quit:
        next_gen_ticker.Stop()
        return
      }
    }
  }()
}

func (game *Game) UpdateTickCallback() {
  game.Mutex.Lock()

  select {
  case chg, ok := <-game.Changes:
    // only host should draw cells directly, the client gets the updated board
    if ok && game.IsHost {
      if chg.Gen > game.Board.Gen {
        fmt.Println("Change from the future, delaying: ", chg)
        game.Changes <- chg
      } else if chg.Gen < game.Board.Gen {
        fmt.Println("Change from the past, applying now: ", chg)
        game.Board.SetCell(chg.Alive, chg.X, chg.Y)
      } else {
        game.Board.SetCell(chg.Alive, chg.X, chg.Y)
      }
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
  if game.IsHost && game.Started {
    game.Board.NextGen()
  }

  defer game.Mutex.Unlock()
}

func (game *Game) PerformSync() {
  game.Mutex.Lock()
  defer game.Mutex.Unlock()

  var sync Sync
  sync.Board = *game.Board
  if game.IsHost {
    game.Client.SendSync(sync)
  }
  if !game.IsHost {
    fmt.Println("Requesting Sync")
    sync, ok := <-game.Syncs
    if ok {
      game.Board = &sync.Board
      fmt.Println("Syncing game")
    }
  }
}
