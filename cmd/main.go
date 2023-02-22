package main

import (
  "gameoflife"
  "github.com/faiface/pixel/pixelgl"
)

func main() {
    pixelgl.Run(gameoflife.Run)
}
