package state

import (
	"tetris-game/grid"

	"github.com/gdamore/tcell/v2"
)

// GameState holds all the information about the current state of the Tetris game.
type GameState struct {
	Screen tcell.Screen
	Grid *grid.Grid
	Style tcell.Style

	CurrentActiveTetrom [][]int
	TempRandomTetrom [][]int
	ActiveTetrom int
	CurrentTetroType int
	CurrentRotation int

	StartX int
	StartY int

	NewX int
	NewY int

	GameOver bool
	GameRunning bool
	RestartChan chan bool

	Score int
	Level int
	TotalLinesCleared int
}
