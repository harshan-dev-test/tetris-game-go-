package game

import (
	tetrominoes "tetris-game/Tetrominoes"
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

	NewX int
	NewY int

	GameOver bool
	GameRunning bool
	RestartChan chan bool

	Score int
	Level int
	TotalLinesCleared int

}

func InitGameState(s tcell.Screen) *GameState{

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	return &GameState{
		Screen: s,
		Grid: grid.NewGrid(10,20),
		Style: defStyle,
		CurrentActiveTetrom: tetrominoes.T[0],
		TempRandomTetrom: tetrominoes.T[0],
		ActiveTetrom: 0,
		CurrentTetroType: 2,
		CurrentRotation: 0,
		StartX: 7,
		StartY: 7,
		NewX: 0,
		NewY: 0,
		GameOver: false,
		GameRunning: true,
		RestartChan: make(chan bool, 1),
		Score: 0,
		Level: 0,
		TotalLinesCleared: 0,
	}

}
