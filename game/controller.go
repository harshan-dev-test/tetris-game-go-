package game

import (
	tetrominoes "tetris-game/Tetrominoes"
	"tetris-game/grid"
	state "tetris-game/state"

	"github.com/gdamore/tcell/v2"
)

// ShowPiece draws the current active tetromino piece on the screen at its current position.
func ShowPiece(state *state.GameState) {
	piece := state.CurrentActiveTetrom

	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			cellX := state.StartX + j
			cellY := state.StartY + i

			if piece[i][j] == 1 {
				// Only draw within the visible play area
				if cellX < 6 || cellX > 25 || cellY < 6 || cellY > 35 {
					continue
				}

				color := tetrominoes.TetrominoColors[state.CurrentTetroType]
				state.Screen.SetContent(cellX, cellY, '█', nil, state.Style.Foreground(color))
			}
		}
	}
}

// ClearPrevPiece erases the previous position of a tetromino piece from the screen.
// It fills the previous cells with the background color.
func ClearPrevPiece(startX, startY int, tetromen *[][]int, s tcell.Screen, style tcell.Style) {
	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece); j++ {
			cellX := startX + j
			cellY := startY + i
			if piece[i][j] == 1 {
				// Only clear within the visible play area
				if cellX < 6 || cellX > 25 || cellY < 6 || cellY > 35 {
					continue
				}
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			}
		}
	}
}

// CanMovePiece checks if a tetromino piece can be moved to a new position without colliding or going out of bounds.
func CanMovePiece(newX, newY int, tetromen *[][]int, state *state.GameState) bool {
	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j
				cellY := newY + i

				gridX := cellX - 6
				gridY := cellY - 6

				// Check if out of grid bounds
				if gridX < 0 || gridX >= state.Grid.Width || gridY < 0 || gridY >= state.Grid.Height {
					return false
				}

				// Check for collision with filled cells
				if state.Grid.Data[gridY][gridX].Filled {
					return false
				}
			}
		}
	}
	return true
}

// RotateTetrom attempts to rotate the current tetromino piece clockwise.
// If the rotation is valid (no collision or out-of-bounds), it updates the piece and redraws it.
func RotateTetrom(state *state.GameState) {
	// Calculate the next rotation index (0-3)-every tetro have 4 types(4 rotations)
	nextRotation := (state.CurrentRotation + 1) % 4
	tempPiece := tetrominoes.AllTetrominos[state.CurrentTetroType][nextRotation]

	if CanMovePiece(state.StartX, state.StartY, &tempPiece, state) {
		ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
		state.CurrentRotation = nextRotation
		state.CurrentActiveTetrom = tempPiece
		ShowPiece(state)
	}
}

// LockGridTetro locks a tetromino piece into the grid at its current position, marking its cells as filled.
func LockGridTetro(newX, newY int, tetromone *[][]int, state *state.GameState) {
	piece := *tetromone
	color := tetrominoes.TetrominoColors[state.CurrentTetroType]
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j - 6
				cellY := newY + i - 6
				// Only lock within the grid bounds
				if cellY >= 0 && cellY < state.Grid.Height && cellX >= 0 && cellX < state.Grid.Width {
					state.Grid.Data[cellY][cellX] = grid.Cell{Filled: true, Color: color}
				}
			}
		}
	}
}
