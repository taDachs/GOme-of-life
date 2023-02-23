package gameoflife

type Change struct {
  X, Y  int
  Alive bool
}

type Sync struct {
  Board Board
}

type Init struct {
  Board Board
}
