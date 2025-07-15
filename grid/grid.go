package grid

import (
	"github.com/gdamore/tcell/v2"
)

// Cell represents a single cell in the Tetris grid.
type Cell struct {
	Filled bool        // Indicates if the cell is filled by a tetromino
	Color  tcell.Color // The color of the filled cell
}

// Grid represents the entire Tetris playfield.
type Grid struct {
	Width  int
	Height int
	Data   [][]Cell
}

// NewGrid creates and returns a new Grid with the specified width and height.
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

// DrawGrid renders the grid onto the tcell screen at the specified starting coordinates.
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
				// Draw empty cell with dark gray background
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			} else {
				// Draw filled cell with its color
				s.SetContent(cellX, cellY, '█', nil, style.Foreground(cell.Color))
			}
		}
	}
}
