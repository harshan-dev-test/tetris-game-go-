package tetrominoes

import (
	"math/rand"
	"tetris-game/state"

)

func GenerateRandomTetromino(state *state.GameState) {
	tetroType := rand.Intn(7)
	state.CurrentTetroType = tetroType
	state.CurrentRotation = 0
	state.CurrentActiveTetrom = AllTetrominos[tetroType][0]
}