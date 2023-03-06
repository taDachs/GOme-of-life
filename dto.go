package gameoflife

const (
  PLAYER_ONE int = iota
  PLAYER_TWO
)

type Change struct {
  X, Y   int
  Alive  bool
  Player int
  Gen    uint
}

type Sync struct {
  Board Board
}

type Init struct {
  Board Board
}
