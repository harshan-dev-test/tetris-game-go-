package game

import (
	"fmt"
	"tetris-game/state"
	"tetris-game/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

var lineScores = map[int]int{
	1: 100,  // Single
	2: 300,  // Double  
	3: 500,  // Triple
	4: 800, // Tetris
}

func DisplayGameStats(state *state.GameState) {

	scoreText := fmt.Sprintf("Score: %d", state.Score)
	utils.DrawText(state.Screen, 30, 8, scoreText, state.Style.Foreground(tcell.ColorDarkCyan))
	
	levelText := fmt.Sprintf("Level: %d", state.Level)
	utils.DrawText(state.Screen, 30, 9, levelText, state.Style.Foreground(tcell.ColorGreen))
	
	linesText := fmt.Sprintf("Lines: %d", state.TotalLinesCleared)
	utils.DrawText(state.Screen, 30, 10, linesText, state.Style.Foreground(tcell.ColorYellow))
}

func CalculateScore(linesCleared int, state *state.GameState) int{
	if linesCleared == 0 {
		return 0
	}
	baseScore := lineScores[linesCleared]
	return baseScore * state.Level
}

func UpdateLevel(state *state.GameState) {
	newLevel := ( state.TotalLinesCleared/ 10) +1
	if newLevel != state.Level {
		state.Level = newLevel
	}
}

func ResetGameStats(state *state.GameState) {
	state.Score = 0
	state.Level = 1
	state.TotalLinesCleared = 0
}

func GetFallingSpeed(state *state.GameState) time.Duration {
	baseSpeed := 600 - (state.Level-1)*50
	if baseSpeed < 100 {
		baseSpeed = 100
	}
	return time.Duration(baseSpeed) * time.Millisecond
}



