package renderer

import (
	"github.com/gdamore/tcell/v2"
)

// Renderer handles terminal rendering
type Renderer struct {
	screen tcell.Screen
}

// NewRenderer creates a new renderer
func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{
		screen: screen,
	}
}

// DrawRune draws a single rune at the given position with colors
func (r *Renderer) DrawRune(x, y int, ch rune, foreground, background tcell.Color) {
	r.screen.SetContent(x, y, ch, nil, tcell.StyleDefault.
		Foreground(foreground).
		Background(background))
}

// DrawText draws text at the given position
func (r *Renderer) DrawText(screen tcell.Screen, x, y int, text string, foreground, background tcell.Color) {
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// DrawTextMasked draws text with hidden characters removed
func (r *Renderer) DrawTextMasked(screen tcell.Screen, x, y int, text string, hidden []bool, foreground, background tcell.Color) {
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	for i, ch := range text {
		if i < len(hidden) && hidden[i] {
			continue
		}
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// DrawTextCentered draws text centered horizontally on a given row
func (r *Renderer) DrawTextCentered(screen tcell.Screen, y int, text string, foreground, background tcell.Color) {
	width, _ := r.screen.Size()
	x := (width - len(text)) / 2
	if x < 0 {
		x = 0
	}
	r.DrawText(screen, x, y, text, foreground, background)
}

// DrawFilledRect draws a filled rectangle
func (r *Renderer) DrawFilledRect(x1, y1, x2, y2 int, ch rune, foreground, background tcell.Color) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			r.DrawRune(x, y, ch, foreground, background)
		}
	}
}

// DrawRect draws a rectangle outline
func (r *Renderer) DrawRect(x1, y1, x2, y2 int, ch rune, foreground, background tcell.Color) {
	// Horizontal lines
	for x := x1; x <= x2; x++ {
		r.DrawRune(x, y1, ch, foreground, background)
		r.DrawRune(x, y2, ch, foreground, background)
	}
	// Vertical lines
	for y := y1; y <= y2; y++ {
		r.DrawRune(x1, y, ch, foreground, background)
		r.DrawRune(x2, y, ch, foreground, background)
	}
}
