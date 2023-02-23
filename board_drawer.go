package gameoflife

import (
  "time"
  "github.com/faiface/pixel"
  "github.com/faiface/pixel/pixelgl"
  "github.com/faiface/pixel/imdraw"
  "golang.org/x/image/colornames"
  "fmt"
)

const SYNC_INTERVAL = 10

func Run(update_url string, sync_url string, game *Game, width, height, res float64) {
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
    if !game.Started {
      continue
    }

    if game.IsHost && game.Board.Gen % SYNC_INTERVAL == 0 {
      go SyncGame(game, sync_url)
    }

    handleClick(game.Board, win, res, update_url)

    select {
    case chg, ok := <-ChangeChannel:
      if ok {
        game.Board.SetCell(chg.Alive, chg.X, chg.Y)
        fmt.Println("Setting cell: ", chg)
      }
    case sync, ok := <-SyncChannel:
      if ok && !game.IsHost {
        game.Board = &sync.Board
        fmt.Println("Syncing game")
      }
    default:
    }
    win.Clear(colornames.Skyblue)
    drawBoard(game.Board, win, int(res))
    win.Update()
    game.Board.NextGen()
    time.Sleep(50 * time.Millisecond)
  }
}

func handleClick(board *Board, win *pixelgl.Window, res float64, url string) {
  if win.JustPressed(pixelgl.MouseButtonLeft) {
    mouse_pos := win.MousePosition()
    change := new(Change)
    change.X = int(mouse_pos.X / res)
    change.Y = int(mouse_pos.Y / res)
    change.Alive = !board.IsAlive(change.X, change.Y)
    board.SetCell(change.Alive, change.X, change.Y)
    changes := make([]Change, 1)
    changes[0] = *change
    SendChanges(changes, url)
  }
}

func drawBoard(board *Board, win *pixelgl.Window, res int) {
  imd := imdraw.New(nil)
  imd.Color = colornames.Black

  for y := 0; y < board.Height; y++ {
    for x := 0; x < board.Width; x++ {
      if board.IsAlive(x, y) {
        imd.Push(pixel.V(float64(x * res), float64(y * res)), pixel.V(float64((x + 1) * res), float64((y + 1) * res)))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}
