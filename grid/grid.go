package grid

import (

	"github.com/gdamore/tcell/v2"
)

type Grid struct {
	Width  int
	Height int
	Data   [][]int
}

func NewGrid(width, height int) *Grid {
	data := make([][]int, height)
	for i := range data {
		data[i] = make([]int, width)
	}

	return &Grid{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

func (g *Grid)DrawGrid(s tcell.Screen, grid *Grid, startX, startY int, style tcell.Style) {
	width := grid.Width
	height := grid.Height
	data := grid.Data

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			cellX := startX + j + 1
			cellY := startY + i + 1

			if data[i][j] == 0 {
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			} else {
				s.SetContent(cellX, cellY, '█', nil, style.Foreground(tcell.ColorGreen))
			}
		}
	}
}