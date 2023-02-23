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
    win.Clear(colornames.Skyblue)
    bd.drawBoard(win)
    win.Update()
    time.Sleep(1000 / FRAME_RATE * time.Millisecond)
  }
}

func (bd *BoardDrawer) handleClick(win *pixelgl.Window) {
  if win.JustPressed(pixelgl.MouseButtonLeft) {
    mouse_pos := win.MousePosition()
    change := new(Change)
    change.X = int(mouse_pos.X / bd.Res)
    change.Y = int(mouse_pos.Y / bd.Res)
    change.Alive = !bd.Game.Board.IsAlive(change.X, change.Y)

    bd.Game.Changes <- *change

    changes := make([]Change, 1)
    changes[0] = *change
    go bd.Game.Client.SendChanges(changes)
  }
}

func (bd *BoardDrawer) drawBoard(win *pixelgl.Window) {
  imd := imdraw.New(nil)
  imd.Color = colornames.Black

  for y := 0; y < bd.Game.Board.Height; y++ {
    for x := 0; x < bd.Game.Board.Width; x++ {
      if bd.Game.Board.IsAlive(x, y) {
        imd.Push(pixel.V(float64(x)*bd.Res, float64(y)*bd.Res), pixel.V(float64(x+1)*bd.Res, float64(y+1)*bd.Res))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}
