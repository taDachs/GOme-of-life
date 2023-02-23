package gameoflife

import (
  "sync"
)

type Game struct {
  Board   *Board
  Started bool
  IsHost  bool
  Mutex   sync.Mutex
}
