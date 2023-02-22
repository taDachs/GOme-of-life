package gameoflife

import (
  "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const width = 1920
const height = 1080

const res = 8
const board_width = width / res
const board_height = height / res

func Run() {
    board := CreateEmptyBoard(board_width, board_height)
    board.InitializeRandom(0.2)

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
        win.Clear(colornames.Skyblue)
        drawBoard(board, win)
		win.Update()
        board.NextGen()
        time.Sleep(50 * time.Millisecond)
	}
}

func drawBoard(board *Board, win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = colornames.Black

    for y, row := range board.Board() {
        for x := range row {
            if board.IsAlive(x, y) {
                imd.Push(pixel.V(float64(x * res), float64(y * res)), pixel.V(float64((x + 1) * res), float64((y + 1) * res)))
                imd.Rectangle(0)
            }
        }
    }

    imd.Draw(win)
}
