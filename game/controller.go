package game

import (
	tetrominoes "tetris-game/Tetrominoes"
	"tetris-game/grid"
	state "tetris-game/state"

	"github.com/gdamore/tcell/v2"
)

func ShowPiece(state *state.GameState) {
	piece := state.CurrentActiveTetrom

	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			cellX := state.StartX + j
			cellY := state.StartY + i

			if piece[i][j] == 1 {
				if cellX < 6 || cellX > 25 || cellY < 6 || cellY > 35 {
					continue
				}

				color := tetrominoes.TetrominoColors[state.CurrentTetroType]
				state.Screen.SetContent(cellX, cellY, '█', nil, state.Style.Foreground(color))
			}
		}
	}
}

func ClearPrevPiece(startX, startY int, tetromen *[][]int, s tcell.Screen, style tcell.Style) {
	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece); j++ {
			cellX := startX + j
			cellY := startY + i
			if piece[i][j] == 1 {
				if cellX < 6 || cellX > 25 || cellY < 6 || cellY > 35 {
					continue
				}
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			}
		}
	}
}

func CanMovePiece(newX, newY int, tetromen *[][]int, state *state.GameState) bool {

	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j
				cellY := newY + i

				gridX := cellX - 6
				gridY := cellY - 6

				if gridX < 0 || gridX >= state.Grid.Width || gridY < 0 || gridY >= state.Grid.Height {
					return false
				}

				if state.Grid.Data[gridY][gridX].Filled {
					return false
				}

			}
		}
	}
	return true
}

func RotateTetrom(state *state.GameState) {

	nextRotation := (state.CurrentRotation + 1) % 4
	tempPiece := tetrominoes.AllTetrominos[state.CurrentTetroType][nextRotation]

	if CanMovePiece(state.StartX, state.StartY, &tempPiece, state) {
		ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
		state.CurrentRotation = nextRotation
		state.CurrentActiveTetrom = tempPiece
		ShowPiece(state)

	}
}

func LockGridTetro(newX, newY int, tetromone *[][]int, state *state.GameState) {
	piece := *tetromone
	color := tetrominoes.TetrominoColors[state.CurrentTetroType]
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j - 6
				cellY := newY + i - 6
				if cellY >= 0 && cellY < state.Grid.Height && cellX >= 0 && cellX < state.Grid.Width {
					state.Grid.Data[cellY][cellX] = grid.Cell{Filled: true, Color: color}
				}
			}
		}
	}
}
