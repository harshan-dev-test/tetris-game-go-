package game

import (
	"fmt"
	tetrominoes "tetris-game/Tetrominoes"
	state "tetris-game/state"
	"tetris-game/utils"
	"time"

	"tetris-game/grid"

	"github.com/gdamore/tcell/v2"
)

// Checks if a given row in the grid is completely filled
func IsLineComplete(state *state.GameState, row int) bool {
	for col := 0; col < state.Grid.Width; col++ {
		if !state.Grid.Data[row][col].Filled {
			return false
		}
	}
	return true
}

// Clears a specific line in the grid and shifts all lines above it down by one
func ClearLine(state *state.GameState, lineIndex int) {
	for row := lineIndex; row > 0; row-- {
		copy(state.Grid.Data[row], state.Grid.Data[row-1])
	}

	for col := 0; col < state.Grid.Width; col++ {
		state.Grid.Data[0][col] = grid.Cell{Filled: false, Color: tcell.ColorDefault}
	}
}

// Clears all completed lines in the grid, updates score and level, and returns the number of lines cleared
func ClearCompletedLines(state *state.GameState) int {
	linesCleared := 0

	for row := state.Grid.Height - 1; row >= 0; row-- {
		if IsLineComplete(state, row) {
			ClearLine(state, row)
			linesCleared++
			row++
		}
	}
	if linesCleared > 0 {
		state.Score += CalculateScore(linesCleared, state)

		state.TotalLinesCleared += linesCleared

		UpdateLevel(state)
	}

	return linesCleared
}

// Main game loop for handling the falling tetromino piece
// Handles piece movement, locking, line clearing, and game over logic
func FallingPieceLoop(state *state.GameState) {

	ticker := time.NewTicker(500 * time.Millisecond) // Initial falling speed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !state.GameRunning {
				return // Exit loop if game is not running
			}
			ticker.Reset(GetFallingSpeed(state)) // Adjust falling speed based on level
		}

		state.NewY = state.StartY + 1 // Attempt to move piece down by one

		if CanMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state) {
			// Move piece down if possible
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			ShowPiece(state)
			DisplayGameStats(state)
			state.Screen.Show()
		} else {
			// Lock piece in place and handle line clearing
			LockGridTetro(state.StartX, state.StartY, &state.CurrentActiveTetrom, state)
			linesCleared := ClearCompletedLines(state)
			state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
			DisplayGameStats(state)

			if linesCleared > 0 {
				// Display special message for line clears
				var lineType string
				switch linesCleared {
				case 1:
					lineType = "Single!"
				case 2:
					lineType = "Double!"
				case 3:
					lineType = "Triple!"
				case 4:
					lineType = "TETRIS!"
				}
				linesClearedText := fmt.Sprintf("%s +%d pts", lineType, CalculateScore(linesCleared, state))
				utils.DrawText(state.Screen, 30, 12, linesClearedText, state.Style.Foreground(tcell.ColorRed))
			}

			tetrominoes.GenerateRandomTetromino(state) // Spawn new tetromino
			state.StartX = 7
			state.StartY = 7
			if CanMovePiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state) {
				ShowPiece(state)
				state.Screen.Show()
			} else {
				// Game over condition
				state.GameOver = true
				state.GameRunning = false
				utils.DrawText(state.Screen, 2, 1, "Game Over", state.Style.Foreground(tcell.ColorRed))
				finalScoreText := fmt.Sprintf("Final Score: %d", state.Score)
				utils.DrawText(state.Screen, 2, 2, finalScoreText, state.Style.Foreground(tcell.ColorWhite))
				state.Screen.Show()
				return
			}
		}
	}
}

// Resets the game grid to an empty state
func ResetGame(state *state.GameState) {
	for y := range state.Grid.Data {
		for x := range state.Grid.Data[y] {
			state.Grid.Data[y][x] = grid.Cell{Filled: false, Color: tcell.ColorDefault}
		}
	}
}

// Initializes the game state, grid, stats, and starts the main game loop
func InitializeGame(state *state.GameState) {

	state.GameOver = false
	state.GameRunning = true
	state.StartX = 7
	state.StartY = 7

	ResetGame(state)
	ResetGameStats(state)

	tetrominoes.GenerateRandomTetromino(state)

	state.Screen.Clear()

	state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
	utils.DrawText(state.Screen, 4, 2, "TETRIS - Press Q to quit", state.Style.Foreground(tcell.ColorBlue))

	DisplayGameStats(state)

	state.Screen.Show()

	go FallingPieceLoop(state)
}
