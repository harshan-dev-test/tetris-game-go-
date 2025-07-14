package game

import (
	"fmt"
	"tetris-game/Tetrominoes"
	state "tetris-game/state"
	"tetris-game/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

func IsLineComplete(state *state.GameState, row int) bool {
	for col := 0; col < state.Grid.Width; col++ {
		if state.Grid.Data[row][col] == 0{
			return false
		}
	}
	return true
}

func ClearLine(state *state.GameState, lineIndex int) {
	for row := lineIndex; row > 0; row-- {
		copy(state.Grid.Data[row], state.Grid.Data[row-1])
	}

	for col := 0; col < state.Grid.Width; col++ {
		state.Grid.Data[0][col]=0
	}
}

func ClearCompletedLines(state *state.GameState) int {
	linesCleared := 0

	for row := state.Grid.Height-1 ; row >= 0; row -- {
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

func FallingPieceLoop(state *state.GameState) {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !state.GameRunning {
				return
			}
			ticker.Reset(GetFallingSpeed(state))
		}

		state.NewY = state.StartY + 1
		
		if CanMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			ShowPiece(state)
			DisplayGameStats(state)
			state.Screen.Show()
		} else {
			LockGridTetro(state.StartX, state.StartY, &state.CurrentActiveTetrom, state)
			linesCleared := ClearCompletedLines(state)
			state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
			DisplayGameStats(state)

			if linesCleared > 0 {
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
				utils.DrawText(state.Screen, 20, 12, linesClearedText, state.Style.Foreground(tcell.ColorRed))
			}

			tetrominoes.GenerateRandomTetromino(state)
			state.StartX = 7
			state.StartY = 7
			if CanMovePiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state) {
				ShowPiece(state)
				state.Screen.Show()
			} else {
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
	utils.DrawText(state.Screen, 2, 2, "TETRIS - Press Q to quit", state.Style.Foreground(tcell.ColorBlue))

	DisplayGameStats(state)

	state.Screen.Show()

	go FallingPieceLoop(state)
}