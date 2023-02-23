package gameoflife

var ChangeChannel = make(chan Change, 10)
var InitChannel = make(chan Init, 10)
var SyncChannel = make(chan Sync, 10)

type Change struct {
  X, Y int
  Alive bool
}

type Sync struct {
  Board Board
}

type Init struct {
  Board Board
}

