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
  Board          ByteBoard
  PlayerOneBoard ByteBoard
  PlayerTwoBoard ByteBoard

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

func (b *Board) SetCell(alive bool, x, y int, player Player) {
  if alive {
    b.Board.SetCell(true, x, y)
    b.PlayerOneBoard.SetCell(player == PLAYER_ONE, x, y)
    b.PlayerTwoBoard.SetCell(player == PLAYER_TWO, x, y)
  } else {
    b.Board.SetCell(false, x, y)
    b.PlayerOneBoard.SetCell(false, x, y)
    b.PlayerTwoBoard.SetCell(false, x, y)
  }
}

func (b *Board) IsPlayerOne(x, y int) bool {
  return b.PlayerOneBoard.IsAlive(x, y)
}

func (b *Board) IsPlayerTwo(x, y int) bool {
  return b.PlayerTwoBoard.IsAlive(x, y)
}

func CreateEmptyBoard(dx, dy int) *Board {
  board := CreateEmptyByteBoard(dx, dy)
  player_one_board := CreateEmptyByteBoard(dx, dy)
  player_two_board := CreateEmptyByteBoard(dx, dy)

  return &Board{*board, *player_one_board, *player_two_board, 0, dx, dy, false, false}
}

func (b *Board) NextGen() {
  newBoard := CreateEmptyBoard(b.Width, b.Height)
  for y := 0; y < b.Height; y++ {
    for x := 0; x < b.Width; x++ {
      is_alive := b.isAliveNextGen(x, y)
      new_born := is_alive && !b.IsAlive(x, y)

      current_player := b.getCurrentPlayer(x, y)
      dominant_player := b.getDominantPlayer(x, y)

      if new_born {
        newBoard.SetCell(is_alive, x, y, dominant_player)
      } else {
        newBoard.SetCell(is_alive, x, y, current_player)
      }
    }
  }

  b.Board = newBoard.Board
  b.PlayerOneBoard = newBoard.PlayerOneBoard
  b.PlayerTwoBoard = newBoard.PlayerTwoBoard
  b.Gen += 1
}

func (b *Board) getCurrentPlayer(x, y int) Player {
  if b.IsPlayerOne(x, y) {
    return PLAYER_ONE
  } else if b.IsPlayerTwo(x, y) {
    return PLAYER_TWO
  } else {
    return NONE
  }
}

func (b *Board) getDominantPlayer(x, y int) Player {
  p1 := 0
  p2 := 0

  for i := -1; i < 2; i++ {
    nx := x + i
    if nx >= b.Width || nx < 0 {
      if b.HWrap {
        nx = (x + b.Width) % b.Width
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
      if ny >= b.Height || ny < 0 {
        if b.VWrap {
          ny = (y + b.Height) % b.Height
        } else {
          continue
        }
      }

      if b.PlayerOneBoard.IsAlive(nx, ny) {
        p1 += 1
      }
      if b.PlayerTwoBoard.IsAlive(nx, ny) {
        p2 += 1
      }
    }
  }

  if p1 > p2 {
    return PLAYER_ONE
  } else if p2 > p1 {
    return PLAYER_TWO
  } else {
    return NONE
  }
}

func (board *Board) InitializeRandom(aliveFraction float32) {
  deadzone := int(float64(board.Height) * 0.1)
  for y := deadzone; y < board.Height-deadzone; y++ {
    for x := 0; x < board.Width; x++ {
      if int(aliveFraction*100) > rand.Intn(100) {
        board.Board.SetCell(true, x, y)
        if int(50) > rand.Intn(100) {
          board.SetCell(true, x, y, PLAYER_ONE)
        } else {
          board.SetCell(true, x, y, PLAYER_TWO)
        }
      } else {
        board.SetCell(false, x, y, NONE)
      }
    }
  }
}

func (board Board) String() string {
  output := "Board:\n"
  output += board.Board.String()
  return output
}
