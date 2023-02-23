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
  SyncInterval    uint
}

func (game *Game) Run() {
  // background game update loop
  update_ticker := time.NewTicker(time.Duration(1000/game.UpdateFrequency) * time.Millisecond)
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

  next_gen_ticker := time.NewTicker(time.Duration(1000/game.GenFrequency) * time.Millisecond)
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

  if game.Board.Gen%game.SyncInterval == 0 && game.Board.Gen > 0 {
    game.PerformSync()
  }

  select {
  case chg, ok := <-game.Changes:
    if ok {
      game.Board.SetCell(chg.Alive, chg.X, chg.Y)
      fmt.Println("Setting cell: ", chg)
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

func (game *Game) PerformSync() {
  var sync Sync
  sync.Board = *game.Board
  game.Client.SendSync(sync)
  fmt.Println("Requesting Sync")
  sync, ok := <-game.Syncs
  if ok && !game.IsHost {
    game.Board = &sync.Board
    fmt.Println("Syncing game")
  }
}
