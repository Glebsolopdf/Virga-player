package overlay

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func DrawText(screen tcell.Screen, x, y int, text string, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func DrawLines(screen tcell.Screen, x, y, width int, lines []Line, bg tcell.Color) {
	for i, line := range lines {
		DrawText(screen, x, y+i, padRightRunes(line.Text, width), line.Color, bg)
	}
}

func padRightRunes(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return string(r[:n])
	}
	return s + strings.Repeat(" ", n-len(r))
}
