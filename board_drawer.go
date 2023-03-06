package gameoflife

import (
  // "fmt"
  "github.com/faiface/pixel"
  "github.com/faiface/pixel/imdraw"
  "github.com/faiface/pixel/pixelgl"
  "golang.org/x/image/colornames"
  "time"
)

const FRAME_RATE = 30

type BoardDrawer struct {
  Game               *Game
  Width, Height, Res float64
}

func (bd *BoardDrawer) Run() {
  cfg := pixelgl.WindowConfig{
    Title:  "Game of life (in GO)",
    Bounds: pixel.R(0, 0, bd.Width, bd.Height),
    VSync:  true,
  }
  win, err := pixelgl.NewWindow(cfg)
  if err != nil {
    panic(err)
  }

  for !win.Closed() {
    bd.handleClick(win)

    win.Clear(colornames.White)

    imd := imdraw.New(nil)
    deadzone := float64(bd.Height) * 0.1

    imd.Color = colornames.Turquoise
    imd.Push(pixel.V(0, 0), pixel.V(bd.Width, deadzone))
    imd.Rectangle(0)

    imd.Color = colornames.Orange
    imd.Push(pixel.V(0, bd.Height), pixel.V(bd.Width, bd.Height-deadzone))
    imd.Rectangle(0)
    imd.Draw(win)

    bd.drawBoard(win)
    win.Update()
    time.Sleep(1000 / FRAME_RATE * time.Millisecond)
  }
}

func (bd *BoardDrawer) handleClick(win *pixelgl.Window) {
  if win.JustPressed(pixelgl.MouseButtonLeft) {
    mouse_pos := win.MousePosition()

    var change Change
    change.X, change.Y = bd.screenToBoard(mouse_pos.X, mouse_pos.Y)
    change.Alive = !bd.Game.Board.IsAlive(change.X, change.Y)
    change.Gen = bd.Game.Board.Gen

    change.Player = bd.Game.Player

    bd.Game.Changes <- change

    changes := make([]Change, 1)
    changes[0] = change
    go bd.Game.Client.SendChanges(changes)
  }
}

func (bd *BoardDrawer) screenToBoard(x, y float64) (int, int) {
  if bd.Game.Player == PLAYER_ONE {
    return int(x / bd.Res), int(y / bd.Res)
  } else {
    return int(x / bd.Res), int((float64(bd.Height) - y) / bd.Res)
  }
}

func (bd *BoardDrawer) boardToScreen(x, y int) (float64, float64) {
  if bd.Game.Player == PLAYER_ONE {
    return float64(x) * bd.Res, float64(y) * bd.Res
  } else {
    return float64(x) * bd.Res, bd.Height - float64(y)*bd.Res
  }
}

func (bd *BoardDrawer) drawBoard(win *pixelgl.Window) {
  imd := imdraw.New(nil)

  for y := 0; y < bd.Game.Board.Height; y++ {
    for x := 0; x < bd.Game.Board.Width; x++ {
      if bd.Game.Board.IsAlive(x, y) {
        imd.Color = colornames.Black
        if bd.Game.Board.IsPlayerOne(x, y) {
          imd.Color = colornames.Blue
        }
        if bd.Game.Board.IsPlayerTwo(x, y) {
          imd.Color = colornames.Red
        }
        board_x1, board_y1 := bd.boardToScreen(x, y)
        board_x2, board_y2 := bd.boardToScreen(x+1, y+1)
        imd.Push(pixel.V(board_x1, board_y1), pixel.V(board_x2, board_y2))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}
