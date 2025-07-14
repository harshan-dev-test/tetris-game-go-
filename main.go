package main

import (
	"fmt"
	"log"
	"tetris-game/Tetrominoes"
	game "tetris-game/game"
	"tetris-game/state"

	// "math/rand"
	"time"

	grid "tetris-game/grid"

	"github.com/gdamore/tcell/v2"
)

var lineScores = map[int]int{
	1: 100,  // Single
	2: 300,  // Double  
	3: 500,  // Triple
	4: 800, // Tetris
}

func main() {

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer s.Fini()

	inputChan := make(chan *tcell.EventKey, 3)

	state := game.InitGameState(s)

	InitializeGame(state)

	go func() {
		for {
			ev := s.PollEvent()
			if keyEvent, ok := ev.(*tcell.EventKey); ok {
				inputChan <- keyEvent
			}
		}
	}()

	for {
		select {
		case ev := <- inputChan:
			if state.GameOver {
				fmt.Println("cha", state.RestartChan)
				game.HandleGameOverInput(ev, state)
			} else {
				HandlerGameInput(ev, state)
			}
		case <- state.RestartChan:
			InitializeGame(state)
		}
	}
}

func drawText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	clearLength := 20

	for i := 0; i < clearLength; i++ {
		s.SetContent(x+i, y, ' ', nil, style)
	}

	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}



func ShowGrid(grid *grid.Grid) {
	data := grid.Data

	for i := range data {
		for j := range data[i] {
			fmt.Print(data[i][j])
		}
		fmt.Println()

	}
}

func ClearPrevPiece(startX, startY int, tetromen *[][]int, s tcell.Screen, style tcell.Style) {
	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece); j++ {
			cellX := startX + j
			cellY := startY + i
			if piece[i][j] == 1 {
				if cellX < 6 || cellX > 15 || cellY < 6 || cellY > 25 {
					continue
				}
				s.SetContent(cellX, cellY, ' ', nil, style.Background(tcell.ColorDarkGray))
			}
		}
	}
}
// ShowPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style, state.Grid, state)
func ShowPiece(state *state.GameState) {
	piece := state.CurrentActiveTetrom

	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			cellX := state.StartX + j
			cellY := state.StartY + i

			if piece[i][j] == 1 {
				if cellX < 6 || cellX > 15 || cellY < 6 || cellY > 25 {
					continue
				}

				color := tetrominoes.TetrominoColors[state.CurrentTetroType]
				state.Screen.SetContent(cellX, cellY, '█', nil, state.Style.Foreground(color))
			}
		}
	}
}

func canMovePiece(newX, newY int, tetromen *[][]int, grid *grid.Grid) bool {

	piece := *tetromen
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j
				cellY := newY + i

				gridX := cellX - 6
				gridY := cellY - 6

				if gridX < 0 || gridX >= grid.Width || gridY < 0 || gridY >= grid.Height {
					return false
				}

				if grid.Data[gridY][gridX] == 1 {
					return false
				}

			}
		}
	}
	return true
}

func RotateTetrom(state *state.GameState) {
	
	nextRotation := (state.CurrentRotation+1) %4
	tempPiece := tetrominoes.AllTetrominos[state.CurrentTetroType][nextRotation]

	if canMovePiece(state.StartX, state.StartY, &tempPiece, state.Grid) {
		ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
		state.CurrentRotation = nextRotation
		state.CurrentActiveTetrom = tempPiece
		ShowPiece(state)

	}
}

func LockGridTetro(newX, newY int, tetromone *[][]int, grid *grid.Grid) {
	piece := *tetromone
	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			if piece[i][j] == 1 {
				cellX := newX + j - 6
				cellY := newY + i - 6

				if cellY >= 0 && cellY < grid.Height && cellX >= 0 && cellX < grid.Width {
					grid.Data[cellY][cellX] = 1
				}
			}
		}
	}
}

func ResetGame(state *state.GameState) {
	for y := range state.Grid.Data {
		for x := range state.Grid.Data[y] {
			state.Grid.Data[y][x] = 0
		}
	}
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
			ticker.Reset(getFallingSpeed(state))
		}

		state.NewY = state.StartY + 1
		
		if canMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state.Grid) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			ShowPiece(state)
			displayGameStats(state)
			state.Screen.Show()
		} else {
			LockGridTetro(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Grid)
			linesCleared := ClearCompletedLines(state)
			state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
			displayGameStats(state)

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
				drawText(state.Screen, 20, 12, linesClearedText, state.Style.Foreground(tcell.ColorRed))
			}

			tetrominoes.GenerateRandomTetromino(state)
			state.StartX = 7
			state.StartY = 7
			if canMovePiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Grid) {
				ShowPiece(state)
				state.Screen.Show()
			} else {
				state.GameOver = true
				state.GameRunning = false
				drawText(state.Screen, 2, 1, "Game Over", state.Style.Foreground(tcell.ColorRed))
				finalScoreText := fmt.Sprintf("Final Score: %d", state.Score)
				drawText(state.Screen, 2, 2, finalScoreText, state.Style.Foreground(tcell.ColorWhite))
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
	resetGameStats(state)

	tetrominoes.GenerateRandomTetromino(state)

	state.Screen.Clear()
	
	state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
	drawText(state.Screen, 2, 2, "TETRIS - Press Q to quit", state.Style.Foreground(tcell.ColorBlue))

	displayGameStats(state)

	state.Screen.Show()

	go FallingPieceLoop(state)
}

func HandlerGameInput(ev *tcell.EventKey, state *state.GameState) {

	switch ev.Key() {
	case tcell.KeyEscape:
		drawText(state.Screen, 2, 1, "ESC pressed. Exiting...", state.Style)
		state.Screen.Show()
		state.Screen.Fini()
		return
	case tcell.KeyRight:
		drawText(state.Screen, 2, 1, "Right", state.Style)
		state.NewX = state.StartX + 1

		if canMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state.Grid) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		displayGameStats(state)
		state.Screen.Show()
	case tcell.KeyLeft:
		drawText(state.Screen, 2, 1, "Left", state.Style)
		drawText(state.Screen, 1, 1, fmt.Sprintf("Left",state.GameRunning), state.Style)
		state.NewX = state.StartX - 1

		if canMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state.Grid) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		displayGameStats(state)
		state.Screen.Show()

	case tcell.KeyDown:
		drawText(state.Screen, 2, 1, "Down", state.Style)
		state.NewY = state.StartY + 1

		if canMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state.Grid) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			state.Score += 1
			ShowPiece(state)
		}
		displayGameStats(state)
		state.Screen.Show()

	case tcell.KeyUp:
		drawText(state.Screen, 2, 1, "Up", state.Style)
		state.NewY = state.StartY -1

		if canMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state.Grid) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			ShowPiece(state)
		}
		displayGameStats(state)
		state.Screen.Show()

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'r', 'R':
			drawText(state.Screen, 2, 1, "Rotate", state.Style)
			RotateTetrom(state)
			displayGameStats(state)
			state.Screen.Show()
		}
	}
}

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

		UpdateFallingSpeed()
	}
}

func getFallingSpeed(state *state.GameState) time.Duration {
	baseSpeed := 600 - (state.Level-1)*50
	if baseSpeed < 100 {
		baseSpeed = 100
	}
	return time.Duration(baseSpeed) * time.Millisecond
}

func UpdateFallingSpeed() {
}

func displayGameStats(state *state.GameState) {

	scoreText := fmt.Sprintf("Score: %d", state.Score)
	drawText(state.Screen, 20, 8, scoreText, state.Style.Foreground(tcell.ColorDarkCyan))
	
	levelText := fmt.Sprintf("Level: %d", state.Level)
	drawText(state.Screen, 20, 9, levelText, state.Style.Foreground(tcell.ColorGreen))
	
	linesText := fmt.Sprintf("Lines: %d", state.TotalLinesCleared)
	drawText(state.Screen, 20, 10, linesText, state.Style.Foreground(tcell.ColorYellow))
}

func resetGameStats(state *state.GameState) {
	state.Score = 0
	state.Level = 1
	state.TotalLinesCleared = 0
}

