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

	for x := 0; x <= width+1; x++ {
		s.SetContent(startX+x, startY, '─', nil, style)
		s.SetContent(startX+x, startY+height+1, '─', nil, style)
	}

	for y := 0; y <= height+1; y++ {
		s.SetContent(startX, startY+y, '|', nil, style)
		s.SetContent(startX+width+1, startY+y, '|', nil, style)
	}

	s.SetContent(startX, startY, '┌', nil, style)
	s.SetContent(startX+width+1, startY, '┐', nil, style)
	s.SetContent(startX, startY+height+1, '└', nil, style)
	s.SetContent(startX+width+1, startY+height+1, '┘', nil, style)

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