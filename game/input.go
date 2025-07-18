package game

import (
	"os"
	"tetris-game/state"
	"tetris-game/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Handle game over input events
// 's' or 'S' restart the game
// 'q' or 'Q' quit the game
func HandleGameOverInput(ev *tcell.EventKey, state *state.GameState) {
	switch ev.Key() {
	case tcell.KeyEscape:
		state.Screen.Fini()
		return
	case tcell.KeyRune:
		switch ev.Rune() {
		case 's', 'S':
			utils.DrawText(state.Screen, 4, 4, "Restarting game...", state.Style.Foreground(tcell.ColorGreen))
			state.Screen.Show()
			time.Sleep(500 * time.Millisecond)
			state.RestartChan <- true
		case 'q', 'Q':
			utils.DrawText(state.Screen, 2, 1, "Q pressed. Exiting.", state.Style.Foreground(tcell.ColorRed))
			state.Screen.Show()
			state.Screen.Fini()
			os.Exit(0)
			return
		}
	}
}

// Handle normal game input events

// Up Arrow: Rotate the tetromino
// Down Arrow: Move the tetromino down faster
// Left Arrow: Move the tetromino left
// Right Arrow: Move the tetromino right
func HandlerGameInput(ev *tcell.EventKey, state *state.GameState) {

	switch ev.Key() {
	case tcell.KeyRight:
		state.NewX = state.StartX + 1

		if CanMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()

	case tcell.KeyLeft:
		state.NewX = state.StartX - 1

		if CanMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()

	case tcell.KeyDown:
		state.NewY = state.StartY + 1

		if CanMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			state.Score += 1
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()

	case tcell.KeyUp:
			RotateTetrom(state)
			DisplayGameStats(state)
			state.Screen.Show()

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q', 'Q':
			utils.DrawText(state.Screen, 2, 1, "Q pressed. Exiting.", state.Style.Foreground(tcell.ColorRed))
			state.Screen.Show()
			state.Screen.Fini()
			os.Exit(0)
			return
		}
	}
}