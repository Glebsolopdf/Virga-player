package playerlyrics

import (
	"virga-player/renderer"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func drawLyricsText(screen tcell.Screen, r *renderer.Renderer, x, y int, text string, fg tcell.Color) {
	theme := settings.CurrentTheme()
	r.DrawText(screen, x, y, text, fg, theme.LyricsBackground)
}

func lyricsCellSize(text string) (int, int) {
	return len([]rune(text)), 1
}
