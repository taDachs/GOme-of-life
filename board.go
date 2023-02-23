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
  Board []bool
  Gen uint
  Width, Height int
}

func (board *Board) isAliveNextGen(x, y int) bool {
  numNeighbours := 0

  for i := -1; i < 2; i++ {
    nx := (x + i + board.Width) % board.Width
    for j := -1; j < 2; j++ {
      if i == 0 && j == 0 {
        continue
      }
      ny := (y + j + board.Height) % board.Height

      if board.IsAlive(nx, ny) {
        numNeighbours += 1
      }
    }
  }

  return (board.IsAlive(x, y) && numNeighbours == 2) || numNeighbours == 3
}

func CreateEmptyBoard(dx, dy int) *Board {
  board := make([]bool, dy * dx)
  for y := 0; y < dy; y++ {
    for x := 0; x < dx; x++ {
      board[y * dx + x] = false
    }
  }

  return &Board{board, 1, dx, dy}
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
  board.Board[y * board.Width + x] = alive
}

func (board *Board) IsAlive(x, y int) bool {
  return board.Board[y * board.Width + x]
}

func (board *Board) InitializeRandom(aliveFraction float32) {
  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      if int(aliveFraction * 100) > rand.Intn(100) {
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
      if board.Board[y * board.Width + x] {
        output += "#"
      } else {
        output += "-"
      }
    }
    output += "\n"
  }
  return output
}
