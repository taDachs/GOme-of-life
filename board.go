package gameoflife

import (
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

type Cell struct {
    alive bool
    x, y int
}

func (c *Cell) isAliveNextGen(board *Board) bool {
    numNeighbours := 0

    for i := -1; i < 2; i++ {
        nx := (c.x + i + board.dx) % board.dx
        for j := -1; j < 2; j++ {
            ny := (c.y + j + board.dy) % board.dy

            if board.IsAlive(nx, ny) {
                numNeighbours += 1
            }
        }
    }

    if c.alive {
        numNeighbours -= 1
    }


    return (c.alive && numNeighbours == 2) || numNeighbours == 3
}

type Board struct {
    board [][]Cell
    gen uint
    dx, dy int
}

func (b *Board) Board() [][]Cell {
    return b.board
}

func CreateEmptyBoard(dx, dy int) *Board {
    board := make([][]Cell, dy)
    for y := 0; y < dy; y++ {
        board[y] = make([]Cell, dx)
        for x := 0; x < dx; x++ {
            board[y][x] = Cell{false, x, y}
        }
    }

    return &Board{board, 0, dx, dy}
}

func (board *Board) NextGen() {
    newBoard := CreateEmptyBoard(board.dx, board.dy)
    for y, row := range board.board {
        for x, c := range row {
            newBoard.SetCell(c.isAliveNextGen(board), x, y)
        }
    }

    board.board = newBoard.board
    board.gen++
}

func (board *Board) SetCell(alive bool, x, y int) {
    board.board[y][x].alive = alive
}

func (board *Board) IsAlive(x, y int) bool {
    return board.board[y][x].alive
}

func (board *Board) InitializeRandom(aliveFraction float32) {
    for y, row := range board.board {
        for x := range row {
            if int(aliveFraction * 100) > rand.Intn(100) {
                board.SetCell(true, x, y)
            } else {
                board.SetCell(false, x, y)
            }
        }
    }
}

func (board Board) String() string {
    output := ""
    for _, row := range board.board {
        for _, c := range row {
            if c.alive {
                output += "#"
            } else {
                output += "-"
            }
        }
        output += "\n"
    }
    return output
}
