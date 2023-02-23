package gameoflife

import (
  // "fmt"
  "github.com/faiface/pixel"
  "github.com/faiface/pixel/imdraw"
  "github.com/faiface/pixel/pixelgl"
  "golang.org/x/image/colornames"
  "time"
)

const SYNC_INTERVAL = 5
const FRAME_RATE = 30

func Run(game *Game, width, height, res float64) {
  cfg := pixelgl.WindowConfig{
    Title:  "Game of life (in GO)",
    Bounds: pixel.R(0, 0, width, height),
    VSync:  true,
  }
  win, err := pixelgl.NewWindow(cfg)
  if err != nil {
    panic(err)
  }

  for !win.Closed() {
    handleClick(game, win, res)
    win.Clear(colornames.Skyblue)
    drawBoard(game.Board, win, int(res))
    win.Update()
    time.Sleep(1000 / FRAME_RATE * time.Millisecond)
  }
}

func handleClick(game *Game, win *pixelgl.Window, res float64) {
  if win.JustPressed(pixelgl.MouseButtonLeft) {
    mouse_pos := win.MousePosition()
    change := new(Change)
    change.X = int(mouse_pos.X / res)
    change.Y = int(mouse_pos.Y / res)
    change.Alive = !game.Board.IsAlive(change.X, change.Y)

    game.Changes <- *change

    changes := make([]Change, 1)
    changes[0] = *change
    go game.Client.SendChanges(changes)
  }
}

func drawBoard(board *Board, win *pixelgl.Window, res int) {
  imd := imdraw.New(nil)
  imd.Color = colornames.Black

  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      if board.IsAlive(x, y) {
        imd.Push(pixel.V(float64(x*res), float64(y*res)), pixel.V(float64((x+1)*res), float64((y+1)*res)))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}
