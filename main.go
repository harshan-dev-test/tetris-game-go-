package main

import (
	"fmt"
	"log"
	"math/rand"
	"tetris-game/Tetrominoes"
	game "tetris-game/game"
	"tetris-game/state"

	// "math/rand"
	"time"

	grid "tetris-game/grid"

	"github.com/gdamore/tcell/v2"
)




var currentActiveTetrom = tetrominoes.T[0]
var tempRandomTetrom = tetrominoes.T[0]
var activeTetrom = 0
var currentTetroType = 2
var currentRotation = 0

var startX = 7
var startY = 7

var gameOver = false
var gameRunning = true
var restartChan = make(chan bool, 1)

var score int = 0
var level int = 0
var totalLinesCleared int = 0

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
			if gameOver {
				fmt.Println("cha", state.RestartChan)
				HandleGameOverInput(ev, state)
			} else {
				HandlerGameInput(ev, state)
			}
		case <- restartChan:
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

func ShowPiece(startX, startY int, tetromen *[][]int, s tcell.Screen, style tcell.Style, grid *grid.Grid) {
	piece := *tetromen

	for i := 0; i < len(piece); i++ {
		for j := 0; j < len(piece[i]); j++ {
			cellX := startX + j
			cellY := startY + i

			if piece[i][j] == 1 {
				if cellX < 6 || cellX > 15 || cellY < 6 || cellY > 25 {
					continue
				}

				color := tetrominoes.TetrominoColors[currentTetroType]
				s.SetContent(cellX, cellY, '█', nil, style.Foreground(color))
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

func RandomTetromGenerator(tempRandomTetrom *[][]int, activeTetrom int) {
	activeTetrom = activeTetrom % 4
	*tempRandomTetrom = tetrominoes.T[activeTetrom]
	
}

func RotateTetrom(currentActiveTetrom *[][]int, startX, startY int, s tcell.Screen, style tcell.Style, grid *grid.Grid) {
	
	nextRotation := (currentRotation+1) %4
	tempPiece := tetrominoes.AllTetrominos[currentTetroType][nextRotation]

	if canMovePiece(startX, startY, &tempPiece, grid) {
		ClearPrevPiece(startX, startY, currentActiveTetrom, s, style)
		currentRotation = nextRotation
		*currentActiveTetrom = tempPiece
		ShowPiece(startX, startY, currentActiveTetrom, s, style, grid)

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

func ResetGame(grid *grid.Grid) {
	for y := range grid.Data {
		for x := range grid.Data[y] {
			grid.Data[y][x] = 0
		}
	}
}

func FallingPieceLoop(s tcell.Screen, grid *grid.Grid, style tcell.Style) {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !gameRunning {
				return
			}
			ticker.Reset(getFallingSpeed())
		}

		newY := startY + 1
		
		if canMovePiece(startX, newY, &currentActiveTetrom, grid) {
			ClearPrevPiece(startX, startY, &currentActiveTetrom, s, style)
			startY = newY
			ShowPiece(startX, startY, &currentActiveTetrom, s, style, grid)
			displayGameStats(s, style)
			s.Show()
		} else {
			LockGridTetro(startX, startY, &currentActiveTetrom, grid)
			linesCleared := ClearCompletedLines(grid)
			grid.DrawGrid(s, grid, 5, 5, style.Foreground(tcell.ColorWhite))
			displayGameStats(s, style)

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
				linesClearedText := fmt.Sprintf("%s +%d pts", lineType, CalculateScore(linesCleared))
				drawText(s, 20, 12, linesClearedText, style.Foreground(tcell.ColorRed))
			}

			GenerateRandomTetromino()
			startX = 7
			startY = 7
			if canMovePiece(startX, startY, &currentActiveTetrom, grid) {
				ShowPiece(startX, startY, &currentActiveTetrom, s, style, grid)
				s.Show()
			} else {
				gameOver = true
				gameRunning = false
				drawText(s, 2, 1, "Game Over", style.Foreground(tcell.ColorRed))
				finalScoreText := fmt.Sprintf("Final Score: %d", score)
				drawText(s, 2, 2, finalScoreText, style.Foreground(tcell.ColorWhite))
				s.Show()
				return
			}
		}
	}
}

func InitializeGame(state *state.GameState) {

	gameOver = false
	gameRunning = true
	startX = 7
	startY = 7

	ResetGame(state.Grid)
	resetGameStats()

	GenerateRandomTetromino()

	state.Screen.Clear()
	
	state.Grid.DrawGrid(state.Screen, state.Grid, 5, 5, state.Style.Foreground(tcell.ColorWhite))
	drawText(state.Screen, 2, 2, "TETRIS - Press Q to quit", state.Style.Foreground(tcell.ColorBlue))

	displayGameStats(state.Screen, state.Style)

	state.Screen.Show()

	go FallingPieceLoop(state.Screen, state.Grid, state.Style)
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
		newX := startX + 1

		if canMovePiece(newX, startY, &currentActiveTetrom, state.Grid) {
			ClearPrevPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style)
			startX = newX
			ShowPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style, state.Grid)
		}
		displayGameStats(state.Screen, state.Style)
		state.Screen.Show()
	case tcell.KeyLeft:
		drawText(state.Screen, 2, 1, "Left", state.Style)
		drawText(state.Screen, 1, 1, fmt.Sprintf("Left",state.GameRunning), state.Style)
		newX := startX - 1

		if canMovePiece(newX, startY, &currentActiveTetrom, state.Grid) {
			ClearPrevPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style)
			startX = newX
			ShowPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style, state.Grid)
		}
		displayGameStats(state.Screen, state.Style)
		state.Screen.Show()

	case tcell.KeyDown:
		drawText(state.Screen, 2, 1, "Down", state.Style)
		newY := startY + 1

		if canMovePiece(startX, newY, &currentActiveTetrom, state.Grid) {
			ClearPrevPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style)
			startY = newY
			score += 1
			ShowPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style, state.Grid)
		}
		displayGameStats(state.Screen, state.Style)
		state.Screen.Show()

	case tcell.KeyUp:
		drawText(state.Screen, 2, 1, "Up", state.Style)
		newY := startY - 1

		if canMovePiece(startX, newY, &currentActiveTetrom, state.Grid) {
			ClearPrevPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style)
			startY = newY
			ShowPiece(startX, startY, &currentActiveTetrom, state.Screen, state.Style, state.Grid)
		}
		displayGameStats(state.Screen, state.Style)
		state.Screen.Show()

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'r', 'R':
			drawText(state.Screen, 2, 1, "Rotate", state.Style)
			RotateTetrom(&currentActiveTetrom, startX, startY, state.Screen, state.Style, state.Grid)
			displayGameStats(state.Screen, state.Style)
			state.Screen.Show()
		}
	}
}

func HandleGameOverInput(ev *tcell.EventKey, state *state.GameState) {
	switch ev.Key() {
	case tcell.KeyEscape:
		state.Screen.Fini()
		return
	case tcell.KeyRune:
		switch ev.Rune() {
		case 's', 'S':
			drawText(state.Screen, 2, 1, "Restarting game...", state.Style.Foreground(tcell.ColorGreen))
			state.Screen.Show()
			time.Sleep(500 * time.Millisecond)
			restartChan <- true
		case 'q', 'Q':
			drawText(state.Screen, 2, 1, "Q pressed. Exiting.", state.Style.Foreground(tcell.ColorRed))
			state.Screen.Show()
			state.Screen.Fini()
			return
		}
	}
}

func GenerateRandomTetromino() {
	tetroType := rand.Intn(7)
	currentTetroType = tetroType
	currentRotation = 0
	currentActiveTetrom = tetrominoes.AllTetrominos[tetroType][0]
}

func IsLineComplete(grid *grid.Grid, row int) bool {
	for col := 0; col < grid.Width; col++ {
		if grid.Data[row][col] == 0{
			return false
		}
	}
	return true
}

func ClearLine(grid *grid.Grid, lineIndex int) {
	for row := lineIndex; row > 0; row-- {
		copy(grid.Data[row], grid.Data[row-1])
	}

	for col := 0; col < grid.Width; col++ {
		grid.Data[0][col]=0
	}
}

func ClearCompletedLines(grid *grid.Grid) int {
	linesCleared := 0

	for row := grid.Height -1 ; row >= 0; row -- {
		if IsLineComplete(grid, row) {
			ClearLine(grid, row)
			linesCleared++
			row++
		}
	}
	if linesCleared > 0 {
		score += CalculateScore(linesCleared)
	
		totalLinesCleared += linesCleared
		
		UpdateLevel()
	}
	
	return linesCleared
}

func CalculateScore(linesCleared int) int{
	if linesCleared == 0 {
		return 0
	}
	baseScore := lineScores[linesCleared]
	return baseScore * level
}

func UpdateLevel() {
	newLevel := (totalLinesCleared / 10) +1
	if newLevel != level {
		level = newLevel

		UpdateFallingSpeed()
	}
}

func getFallingSpeed() time.Duration {
	baseSpeed := 600 - (level-1)*50
	if baseSpeed < 100 {
		baseSpeed = 100
	}
	return time.Duration(baseSpeed) * time.Millisecond
}

func UpdateFallingSpeed() {
}

func displayGameStats(s tcell.Screen, style tcell.Style) {

	scoreText := fmt.Sprintf("Score: %d", score)
	drawText(s, 20, 8, scoreText, style.Foreground(tcell.ColorDarkCyan))
	
	levelText := fmt.Sprintf("Level: %d", level)
	drawText(s, 20, 9, levelText, style.Foreground(tcell.ColorGreen))
	
	linesText := fmt.Sprintf("Lines: %d", totalLinesCleared)
	drawText(s, 20, 10, linesText, style.Foreground(tcell.ColorYellow))
}

func resetGameStats() {
	score = 0
	level = 1
	totalLinesCleared = 0
}

