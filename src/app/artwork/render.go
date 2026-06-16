package artwork

import "github.com/gdamore/tcell/v2"

func (a *Artwork) Render(screen tcell.Screen) {
	a.render(screen, true, true)
}

func (a *Artwork) RenderInfoOnly(screen tcell.Screen) {
	a.render(screen, true, false)
}

func (a *Artwork) RenderLyricsOverlay(screen tcell.Screen) {
	a.render(screen, false, true)
}

func (a *Artwork) render(screen tcell.Screen, drawInfo, drawLyrics bool) {
	state := a.snapshot()

	if state.animationEnabled && state.mode == DisplaySixel {
		a.renderTextOnly(screen, state, drawInfo, drawLyrics)
		return
	}

	switch state.mode {
	case DisplaySixel:
		a.renderSixel(screen, state, drawInfo, drawLyrics)
	case DisplayTextOnly:
		a.renderTextOnly(screen, state, drawInfo, drawLyrics)
	}
}
