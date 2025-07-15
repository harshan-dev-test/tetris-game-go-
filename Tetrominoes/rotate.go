package tetrominoes

import (
	"math/rand"
	"tetris-game/state"
)

// GenerateRandomTetromino selects a random tetromino type, resets its rotation, and updates the game state.
func GenerateRandomTetromino(state *state.GameState) {
	tetroType := rand.Intn(7) // Randomly select one of the 7 tetromino types
	state.CurrentTetroType = tetroType
	state.CurrentRotation = 0
	state.CurrentActiveTetrom = AllTetrominos[tetroType][0]
}
