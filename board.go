package gameoflife

import (
  "math/rand"
  "time"
)

func InitSeed() {
  rand.Seed(time.Now().UTC().UnixNano())
  // rand.Seed(1)
}

type Board struct {
  Board         ByteBoard
  Gen           uint
  Width, Height int
  HWrap         bool
  VWrap         bool
}

func (board *Board) isAliveNextGen(x, y int) bool {
  numNeighbours := 0

  for i := -1; i < 2; i++ {
    nx := x + i
    if nx >= board.Width || nx < 0 {
      if board.HWrap {
        nx = (x + board.Width) % board.Width
      } else {
        continue
      }
    }
    for j := -1; j < 2; j++ {
      // skip itself
      if i == 0 && j == 0 {
        continue
      }
      ny := y + j
      if ny >= board.Height || ny < 0 {
        if board.VWrap {
          ny = (y + board.Height) % board.Height
        } else {
          continue
        }
      }

      if board.Board.IsAlive(nx, ny) {
        numNeighbours += 1
      }
    }
  }

  return (board.Board.IsAlive(x, y) && numNeighbours == 2) || numNeighbours == 3
}

func (b *Board) IsAlive(x, y int) bool {
  return b.Board.IsAlive(x, y)
}

func (b *Board) SetCell(alive bool, x, y int) {
  b.Board.SetCell(alive, x, y)
}

func CreateEmptyBoard(dx, dy int) *Board {
  board := CreateEmptyByteBoard(dx, dy)

  return &Board{*board, 0, dx, dy, false, false}
}

func (board *Board) NextGen() {
  newBoard := CreateEmptyBoard(board.Width, board.Height)
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      newBoard.Board.SetCell(board.isAliveNextGen(x, y), x, y)
    }
  }

  board.Board = newBoard.Board
  board.Gen++
}

func (board *Board) InitializeRandom(aliveFraction float32) {
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      if int(aliveFraction*100) > rand.Intn(100) {
        board.Board.SetCell(true, x, y)
      } else {
        board.Board.SetCell(false, x, y)
      }
    }
  }
}

func (board Board) String() string {
  output := "Board:\n"
  output += board.Board.String()
  return output
}
