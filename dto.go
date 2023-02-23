package gameoflife

type Change struct {
  X, Y  int
  Alive bool
  Gen   uint
}

type Sync struct {
  Board Board
}

type Init struct {
  Board Board
}
