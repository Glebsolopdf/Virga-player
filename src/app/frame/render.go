package frame

import (
	"unicode/utf8"

	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func drawTextMasked(screen tcell.Screen, x, y int, text string, hidden []bool, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for offset, ch := range []rune(text) {
		if offset < len(hidden) && hidden[offset] {
			continue
		}
		screen.SetContent(x+offset, y, ch, nil, style)
	}
}

func drawTextCentered(screen tcell.Screen, y int, text string, fg, bg tcell.Color) {
	w, _ := screen.Size()
	x := (w - utf8.RuneCountInString(text)) / 2
	if x < 0 {
		x = 0
	}
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for i, ch := range []rune(text) {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func (f Frame) Render(dt float64) {
	theme := settings.CurrentTheme()
	f.Screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.MessageText))
	f.Screen.Clear()

	if !f.PlayerEnabled || f.Player == nil {
		if !f.Message.Converted {
			drawTextMasked(f.Screen, f.Message.X, f.Message.Y, f.Message.Text, f.Message.Hidden, theme.MessageText, theme.Background)
		}
	}

	f.ParticleSystem.Update(dt)
	if f.Debug != nil {
		f.Debug.UpdateRuntime(dt, len(f.ParticleSystem.GetParticles()), f.MaxParticles, f.TargetFPS)
	}
	if !f.PlayerEnabled && f.MessageErasable {
		f.ParticleSystem.HitMessage(f.Message.Text, f.Message.X, f.Message.Y, f.Message.Hidden)
	}
	if f.PlayerEnabled && f.Player != nil {
		f.renderPlayerScene()
	} else {
		f.ParticleSystem.Draw(f.Screen)
	}
	if f.Debug != nil {
		f.Debug.DrawOverlay(f.Screen)
	}
	if f.FooterPromptText != "" {
		_, h := f.Screen.Size()
		if h > 0 {
			theme := settings.CurrentTheme()
			drawTextCentered(f.Screen, h-1, f.FooterPromptText, theme.SettingsHint, theme.Background)
		}
	}

	f.Screen.Show()
}

func (f Frame) renderPlayerScene() {
	f.renderStage(settings.RainLayerBehind)
	f.ParticleSystem.DrawBackLayers(f.Screen)
	f.renderStage(settings.RainLayerBetween)
	f.ParticleSystem.DrawFrontLayers(f.Screen)
	f.renderStage(settings.RainLayerFront)
}

func (f Frame) renderStage(mode settings.RainLayerMode) {
	if f.PlayerRainLayer == mode {
		f.renderPlayerInfo()
	}
	if f.LyricsRainLayer == mode {
		f.renderPlayerLyricsOverlay()
	}
}

func (f Frame) renderPlayerInfo() {
	p := f.Player
	if p.Artwork != nil {
		p.Artwork.RenderInfoOnly(f.Screen)
	}
}

func (f Frame) renderPlayerLyricsOverlay() {
	p := f.Player
	if p.Artwork != nil {
		p.Artwork.RenderLyricsOverlay(f.Screen)
	}
}
