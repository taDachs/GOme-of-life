package gameoflife

import (
  "math"
  "math/rand"
  "time"
)

func InitSeed() {
  rand.Seed(time.Now().UTC().UnixNano())
  // rand.Seed(1)
}

type Board struct {
  Board         []byte
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
      if ny >= board.Width || ny < 0 {
        if board.VWrap {
          ny = (y + board.Height) % board.Height
        } else {
          continue
        }
      }

      if board.IsAlive(nx, ny) {
        numNeighbours += 1
      }
    }
  }

  return (board.IsAlive(x, y) && numNeighbours == 2) || numNeighbours == 3
}

func CreateEmptyBoard(dx, dy int) *Board {
  board := make([]byte, int(math.Ceil(float64(dy)*float64(dx)/8.0))) // 8 bits in a byte
  for y := 0; y < dy; y++ {
    for x := 0; x < dx; x++ {
      i := int(math.Trunc(float64(y*dx+x) / 8.0))
      board[i] = 0
    }
  }

  return &Board{board, 1, dx, dy, false, false}
}

func (board *Board) NextGen() {
  newBoard := CreateEmptyBoard(board.Width, board.Height)
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      newBoard.SetCell(board.isAliveNextGen(x, y), x, y)
    }
  }

  board.Board = newBoard.Board
  board.Gen++
}

func (board *Board) SetCell(alive bool, x, y int) {
  i := int(math.Trunc(float64(y*board.Width+x) / 8.0))
  offset := (y*board.Width + x) % 8
  if alive {
    board.Board[i] |= (1 << offset)
  } else {
    board.Board[i] &= ^(1 << offset)
  }
}

func (board *Board) IsAlive(x, y int) bool {
  // return board.Board[y * board.Width + x]
  i := int(math.Trunc(float64(y*board.Width+x) / 8.0))
  offset := (y*board.Width + x) % 8
  return (board.Board[i] & byte(1<<offset)) > 0
}

func (board *Board) InitializeRandom(aliveFraction float32) {
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      if int(aliveFraction*100) > rand.Intn(100) {
        board.SetCell(true, x, y)
      } else {
        board.SetCell(false, x, y)
      }
    }
  }
}

func (board Board) String() string {
  output := "Board:\n"
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      i := int(math.Trunc(float64(y*board.Width+x) / 8.0))
      offset := (y*board.Width + x) % 8
      if board.Board[i]&byte(1<<offset) > 0 {
        output += "#"
      } else {
        output += "-"
      }
    }
    output += "\n"
  }
  return output
}
