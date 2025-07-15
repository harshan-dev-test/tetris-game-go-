package main

import (
	"log"
	game "tetris-game/game"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Create a new tcell screen & initialize the screen for rendering the game
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	// Ensure the screen is finalized when the program exits
	defer s.Fini()

	// Create a buffered channel to handle keyboard input events
	inputChan := make(chan *tcell.EventKey, 3)

	// Initialize the game state & game
	state := game.InitGameState(s)
	game.InitializeGame(state)

	// Start a goroutine to poll for keyboard events and send them to inputChan
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
		case ev := <-inputChan:
			if state.GameOver {
				// Handle game over input event (restart or quit)
				game.HandleGameOverInput(ev, state)
			} else {
				// Handle normal game input event (move, rotate, drop)
				game.HandlerGameInput(ev, state)
			}
		// Handle game restart event
		case <-state.RestartChan:
			// Re-initialize the game state to start a new game
			game.InitializeGame(state)
		}
	}
}
