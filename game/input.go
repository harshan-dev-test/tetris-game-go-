package game

import (
	// "math/rand"
	"fmt"
	"os"
	"tetris-game/state"
	"tetris-game/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

func HandleGameOverInput(ev *tcell.EventKey, state *state.GameState) {
	switch ev.Key() {
	case tcell.KeyEscape:
		state.Screen.Fini()
		return
	case tcell.KeyRune:
		switch ev.Rune() {
		case 's', 'S':
			utils.DrawText(state.Screen, 2, 1, "Restarting game...", state.Style.Foreground(tcell.ColorGreen))
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

func HandlerGameInput(ev *tcell.EventKey, state *state.GameState) {

	switch ev.Key() {
	case tcell.KeyEscape:
		utils.DrawText(state.Screen, 2, 1, "ESC pressed. Exiting...", state.Style)
		state.Screen.Show()
		state.Screen.Fini()
		return
	case tcell.KeyRight:
		utils.DrawText(state.Screen, 2, 1, "Right", state.Style)
		state.NewX = state.StartX + 1

		if CanMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()
	case tcell.KeyLeft:
		utils.DrawText(state.Screen, 2, 1, "Left", state.Style)
		utils.DrawText(state.Screen, 1, 1, fmt.Sprintf("Left",state.GameRunning), state.Style)
		state.NewX = state.StartX - 1

		if CanMovePiece(state.NewX, state.StartY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartX = state.NewX
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()

	case tcell.KeyDown:
		utils.DrawText(state.Screen, 2, 1, "Down", state.Style)
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
		utils.DrawText(state.Screen, 2, 1, "Up", state.Style)
		state.NewY = state.StartY -1

		if CanMovePiece(state.StartX, state.NewY, &state.CurrentActiveTetrom, state) {
			ClearPrevPiece(state.StartX, state.StartY, &state.CurrentActiveTetrom, state.Screen, state.Style)
			state.StartY = state.NewY
			ShowPiece(state)
		}
		DisplayGameStats(state)
		state.Screen.Show()

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'r', 'R':
			utils.DrawText(state.Screen, 2, 1, "Rotate", state.Style)
			RotateTetrom(state)
			DisplayGameStats(state)
			state.Screen.Show()
		case 'q', 'Q':
			utils.DrawText(state.Screen, 2, 1, "Q pressed. Exiting.", state.Style.Foreground(tcell.ColorRed))
			state.Screen.Show()
			state.Screen.Fini()
			os.Exit(0)
			return
		}
	}
}