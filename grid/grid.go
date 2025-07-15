package grid

import (
	"github.com/gdamore/tcell/v2"
)

type Cell struct {
	Filled bool
	Color  tcell.Color
}

type Grid struct {
	Width  int
	Height int
	Data   [][]Cell
}

func NewGrid(width, height int) *Grid {
	data := make([][]Cell, height)
	for i := range data {
		data[i] = make([]Cell, width)
	}

	return &Grid{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

func (g *Grid) DrawGrid(s tcell.Screen, grid *Grid, startX, startY int, style tcell.Style) {
	width := grid.Width
	height := grid.Height
	data := grid.Data

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			cellX := startX + j + 1
			cellY := startY + i + 1
			cell := data[i][j]
			if !cell.Filled {
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			} else {
				s.SetContent(cellX, cellY, '█', nil, style.Foreground(cell.Color))
			}
		}
	}
}