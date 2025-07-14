package game

import (
	// "math/rand"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)


func HandleGameOverInput(ev *tcell.EventKey, state *GameState) {
	fmt.Println("chann", state.RestartChan)
	switch ev.Key() {
	case tcell.KeyEscape:
		state.Screen.Fini()
		return
	case tcell.KeyRune:
		switch ev.Rune() {
		case 's', 'S':
			//drawText(state.Screen, 2, 1, "Restarting game...", state.Style.Foreground(tcell.ColorGreen))
			state.Screen.Show()
			time.Sleep(500 * time.Millisecond)
			state.RestartChan <- true
		case 'q', 'Q':
			//drawText(state.Screen, 2, 1, "Q pressed. Exiting.", state.Style.Foreground(tcell.ColorRed))
			state.Screen.Show()
			state.Screen.Fini()
			return
		}
	}
}