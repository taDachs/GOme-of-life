package gameoflife

type Player int

const (
  PLAYER_ONE Player = iota
  PLAYER_TWO
  NONE
)

type Change struct {
  X, Y   int
  Alive  bool
  Player Player
  Gen    uint
}

type Sync struct {
  Board Board
}

type Init struct {
  Board Board
}
