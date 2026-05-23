package artwork

import "github.com/gdamore/tcell/v2"

func (a *Artwork) Render(screen tcell.Screen) {
	state := a.snapshot()

	if state.animationEnabled && state.mode == DisplaySixel {
		a.renderTextOnly(screen, state)
		return
	}

	switch state.mode {
	case DisplaySixel:
		a.renderSixel(screen, state)
	case DisplayTextOnly:
		a.renderTextOnly(screen, state)
	}
}
