package game

import (
	grid "tetris-game/grid"

	"github.com/gdamore/tcell/v2"
)

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

	GameOver bool
	GameRunning bool
	RestartChan chan bool

	Score int
	Level int
	TotalLinesCleared int

}
