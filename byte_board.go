package gameoflife

import (
  "math"
)

type ByteBoard struct {
  Board         []byte
  Width, Height int
}

func CreateEmptyByteBoard(dx, dy int) *ByteBoard {
  byte_board := make([]byte, int(math.Ceil(float64(dy)*float64(dx)/8.0))) // 8 bits in a byte
  for y := 0; y < dy; y++ {
    for x := 0; x < dx; x++ {
      i := int(math.Trunc(float64(y*dx+x) / 8.0))
      byte_board[i] = 0
    }
  }

  board := new(ByteBoard)

  board.Board = byte_board
  board.Width = dx
  board.Height = dy

  return board
}

func (board *ByteBoard) SetCell(alive bool, x, y int) {
  i := int(math.Trunc(float64(y*board.Width+x) / 8.0))
  offset := (y*board.Width + x) % 8
  if alive {
    board.Board[i] |= (1 << offset)
  } else {
    board.Board[i] &= ^(1 << offset)
  }
}

func (board *ByteBoard) IsAlive(x, y int) bool {
  // return board.Board[y * board.Width + x]
  i := int(math.Trunc(float64(y*board.Width+x) / 8.0))
  offset := (y*board.Width + x) % 8
  return (board.Board[i] & byte(1<<offset)) > 0
}

func (board ByteBoard) String() string {
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
