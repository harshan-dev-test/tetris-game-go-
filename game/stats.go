package game

import (
	"fmt"
	"tetris-game/state"
	"tetris-game/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

// lines cleared at once to the corresponding score.
var lineScores = map[int]int{
	1: 100, // Single
	2: 300, // Double
	3: 500, // Triple
	4: 800, // Tetris
}

// DisplayGameStats renders the current score, level, and total lines cleared on the screen.
func DisplayGameStats(state *state.GameState) {

	scoreText := fmt.Sprintf("Score: %d", state.Score)
	utils.DrawText(state.Screen, 30, 8, scoreText, state.Style.Foreground(tcell.ColorDarkCyan))

	levelText := fmt.Sprintf("Level: %d", state.Level)
	utils.DrawText(state.Screen, 30, 9, levelText, state.Style.Foreground(tcell.ColorGreen))

	linesText := fmt.Sprintf("Lines: %d", state.TotalLinesCleared)
	utils.DrawText(state.Screen, 30, 10, linesText, state.Style.Foreground(tcell.ColorYellow))
}

// CalculateScore returns the score for the number of lines cleared in a single move, based on current level.
func CalculateScore(linesCleared int, state *state.GameState) int {
	if linesCleared == 0 {
		return 0
	}
	baseScore := lineScores[linesCleared]
	return baseScore * state.Level
}

// UpdateLevel (based on totla lines cleared) increases by 1 for every 10 lines cleared.
func UpdateLevel(state *state.GameState) {
	newLevel := (state.TotalLinesCleared / 10) + 1
	if newLevel != state.Level {
		state.Level = newLevel
	}
}

// ResetGameStats resets the score, level, and total lines cleared to their initial values.
func ResetGameStats(state *state.GameState) {
	state.Score = 0
	state.Level = 1
	state.TotalLinesCleared = 0
}

// GetFallingSpeed returns the current falling speed of tetrominoes based on the level.
func GetFallingSpeed(state *state.GameState) time.Duration {
	baseSpeed := 600 - (state.Level-1)*50
	if baseSpeed < 100 {
		baseSpeed = 100
	}
	return time.Duration(baseSpeed) * time.Millisecond
}
