package gameoflife

import (
  "testing"
  //"fmt"
)

func TestCreateEmptyBoard(t *testing.T) {
  board := CreateEmptyBoard(8, 8)
  if len(board.Board.Board) != 8 {
    t.Errorf("invalid size")
  } else {
    for i := 0; i < 8; i++ {
      for j := 0; j < 8; j++ {
        if board.IsAlive(i, j) {
          t.Errorf("all cells should be dead")
        }
      }
    }
  }
}

func TestNextGen(t *testing.T) {
  board := CreateEmptyBoard(8, 8)
  board.SetCell(true, 0, 0, NONE) // breaks other tests somehow
  board.NextGen()
  if board.IsAlive(0, 0) {
    t.Errorf("cell should be dead")
  }
  board = nil
}

func TestInitializeRandom(t *testing.T) {
  board := CreateEmptyBoard(8, 8)
  board.InitializeRandom(0)
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      if board.IsAlive(i, j) {
        t.Errorf("all cells should be dead")
      }
    }
  }
  board.InitializeRandom(1)
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      if !board.IsAlive(i, j) {
        t.Errorf("all cells should be alive")
      }
    }
  }
}

// this test is just for that one weird bug that keep popping up: for some reason the cell at (0,0) has alive=true. Removing the board.SetCell(true, 0, 0) in the above test fixes this somehow
func TestZombieCell(t *testing.T) {
  board := CreateEmptyBoard(8, 8)
  if board.IsAlive(0, 0) {
    t.Error("this cell should be dead")
  }
  board.SetCell(true, 6, 6, NONE)
  board.SetCell(true, 6, 7, NONE)
  board.SetCell(true, 7, 6, NONE)
  board.SetCell(true, 7, 7, NONE)

  for i := 0; i < 20; i++ {
    board.NextGen()
    if board.IsAlive(0, 0) {
      t.Error("this cell should be dead")
    }
  }
}

func TestStillLife(t *testing.T) {
  board := CreateEmptyBoard(8, 8)
  board.SetCell(true, 6, 6, NONE)
  board.SetCell(true, 6, 7, NONE)
  board.SetCell(true, 7, 6, NONE)
  board.SetCell(true, 7, 7, NONE)

  for i := 0; i < 20; i++ {
    board.NextGen()

    if !board.IsAlive(6, 6) || !board.IsAlive(6, 7) || !board.IsAlive(7, 6) || !board.IsAlive(7, 7) {
      t.Error("cube is not stable")
    }
  }
}
