package utils

import (
	"github.com/gdamore/tcell/v2"
)

// DrawText draws a string of text at the specified (x, y) position on the given tcell.Screen with the provided style.
// It first clears a fixed length of the line to ensure no leftover characters remain from previous draws.
//
// Parameters:
//
//	s     - the tcell.Screen to draw on
//	x, y  - the starting coordinates for the text
//	text  - the string to display
//	style - the tcell.Style to use for the text
func DrawText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	clearLength := 30 // Number of characters to clear before drawing new text

	// Clear the line by overwriting with spaces
	for i := 0; i < clearLength; i++ {
		s.SetContent(x+i, y, ' ', nil, style)
	}

	// Draw the actual text
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
