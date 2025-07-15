package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	game "tetris-game/game"
)

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

	game.InitializeGame(state)

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
				game.HandleGameOverInput(ev, state)
			} else {
				game.HandlerGameInput(ev, state)
			}
		case <-state.RestartChan:
			game.InitializeGame(state)
		}
	}
}