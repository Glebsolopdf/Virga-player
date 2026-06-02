package renderer

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

type Renderer struct {
	screen tcell.Screen
}

func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{
		screen: screen,
	}
}

func (r *Renderer) DrawRune(x, y int, ch rune, foreground, background tcell.Color) {
	r.screen.SetContent(x, y, ch, nil, tcell.StyleDefault.
		Foreground(foreground).
		Background(background))
}

func (r *Renderer) DrawText(screen tcell.Screen, x, y int, text string, foreground, background tcell.Color) {
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	for offset, ch := range []rune(text) {
		screen.SetContent(x+offset, y, ch, nil, style)
	}
}

func (r *Renderer) DrawTextMasked(screen tcell.Screen, x, y int, text string, hidden []bool, foreground, background tcell.Color) {
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	for offset, ch := range []rune(text) {
		if offset < len(hidden) && hidden[offset] {
			continue
		}
		screen.SetContent(x+offset, y, ch, nil, style)
	}
}

func (r *Renderer) DrawTextCentered(screen tcell.Screen, y int, text string, foreground, background tcell.Color) {
	width, _ := r.screen.Size()
	x := (width - utf8.RuneCountInString(text)) / 2
	if x < 0 {
		x = 0
	}
	r.DrawText(screen, x, y, text, foreground, background)
}
