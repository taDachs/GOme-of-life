package gameoflife

import (
  "math/rand"
  "time"
)

func InitSeed() {
  rand.Seed(time.Now().UTC().UnixNano())
  // rand.Seed(1)
}

const NUMBER_OBJECTIVES int = 4

type Point struct {
  X, Y int
}

type Board struct {
  Board               ByteBoard
  PlayerOneBoard      ByteBoard
  PlayerTwoBoard      ByteBoard
  PlayerOneObjectives []Point
  PlayerTwoObjectives []Point

  Gen           uint
  Width, Height int
  HWrap         bool
  VWrap         bool
  Deadzone      int
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
  board := new(Board)
  board.Board = *CreateEmptyByteBoard(dx, dy)
  board.PlayerOneBoard = *CreateEmptyByteBoard(dx, dy)
  board.PlayerTwoBoard = *CreateEmptyByteBoard(dx, dy)

  board.Gen = 0
  board.Width = dx
  board.Height = dy
  board.VWrap = false
  board.HWrap = false
  return board
}

func (b *Board) SetupPlayerAreas(deadzone float64) {
  b.Deadzone = int(float64(b.Height) * deadzone)

  h_spacing := b.Width / (NUMBER_OBJECTIVES + 1)
  v_spacing := b.Deadzone / 3

  for i := 1; i <= NUMBER_OBJECTIVES; i++ {
    for xi := 0; xi <= 1; xi++ {
      for yi := 0; yi <= 1; yi++ {
        x := i*h_spacing - xi
        y1 := v_spacing - yi

        y2 := b.Height - y1 - 1
        b.SetCell(true, x, y1, PLAYER_ONE)
        b.SetCell(true, x, y2, PLAYER_TWO)

        b.PlayerOneObjectives = append(b.PlayerOneObjectives, Point{x, y1})
        b.PlayerTwoObjectives = append(b.PlayerTwoObjectives, Point{x, y2})
      }
    }
  }
}

func (b *Board) NextGen() {
  newBoard := CreateEmptyBoard(b.Width, b.Height)
  for y := 0; y < b.Height; y++ {
    for x := 0; x < b.Width; x++ {
      is_alive := b.isAliveNextGen(x, y)
      new_born := is_alive && !b.IsAlive(x, y)
      died := !is_alive && b.IsAlive(x, y)

      current_player := b.getCurrentPlayer(x, y)
      dominant_player := b.getDominantPlayer(x, y)

      if new_born {
        newBoard.SetCell(is_alive, x, y, dominant_player)
      } else {
        newBoard.SetCell(is_alive, x, y, current_player)
      }

      if died {
        b.updateObjectives(x, y)
      }
    }
  }

  b.Board = newBoard.Board
  b.PlayerOneBoard = newBoard.PlayerOneBoard
  b.PlayerTwoBoard = newBoard.PlayerTwoBoard
  b.Gen += 1
}

func (b *Board) updateObjectives(x, y int) {
  if b.IsPlayerOneObjective(x, y) {
    var new_objectives []Point
    for _, p := range b.PlayerOneObjectives {
      if !(p.X == x && p.Y == y) {
        new_objectives = append(new_objectives, p)
      }
    }
    b.PlayerOneObjectives = new_objectives
  }
  if b.IsPlayerTwoObjective(x, y) {
    var new_objectives []Point
    for _, p := range b.PlayerTwoObjectives {
      if !(p.X == x && p.Y == y) {
        new_objectives = append(new_objectives, p)
      }
    }
    b.PlayerTwoObjectives = new_objectives
  }
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

      if b.IsPlayerOne(nx, ny) {
        p1 += 1
      }
      if b.IsPlayerTwo(nx, ny) {
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

func (b *Board) IsPlayerOneAlive() bool {
  for _, p := range b.PlayerOneObjectives {
    if b.IsPlayerOne(p.X, p.Y) {
      return true
    }
  }
  return false
}

func (b *Board) IsPlayerTwoAlive() bool {
  for _, p := range b.PlayerTwoObjectives {
    if b.IsPlayerTwo(p.X, p.Y) {
      return true
    }
  }
  return false
}

func (board *Board) InitializeRandom(aliveFraction float32) {
  for y := board.Deadzone; y < board.Height-board.Deadzone; y++ {
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

func (b *Board) IsPlayerOneObjective(x, y int) bool {
  for _, p := range b.PlayerOneObjectives {
    if p.X == x && p.Y == y {
      return true
    }
  }
  return false
}

func (b *Board) IsPlayerTwoObjective(x, y int) bool {
  for _, p := range b.PlayerTwoObjectives {
    if p.X == x && p.Y == y {
      return true
    }
  }
  return false
}

func (b *Board) IsObjective(x, y int) bool {
  return b.IsPlayerOneObjective(x, y) || b.IsPlayerTwoObjective(x, y)
}

func (b *Board) IsClickAllowed(x, y int, player Player) bool {
  if player == PLAYER_ONE {
    return y < b.Deadzone
  }

  if player == PLAYER_TWO {
    return y >= b.Height-b.Deadzone
  }

  return true
}

func (board Board) String() string {
  output := "Board:\n"
  output += board.Board.String()
  return output
}
