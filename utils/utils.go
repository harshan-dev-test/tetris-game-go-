package utils

import (
	"github.com/gdamore/tcell/v2"
)

func DrawText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	clearLength := 30

	for i := 0; i < clearLength; i++ {
		s.SetContent(x+i, y, ' ', nil, style)
	}

	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
