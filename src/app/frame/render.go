package frame

import (
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (f Frame) Render(dt float64) {
	theme := settings.CurrentTheme()
	f.Screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.MessageText))
	f.Screen.Clear()

	if !f.PlayerEnabled || f.Player == nil {
		if !f.Message.Converted {
			f.Renderer.DrawTextMasked(f.Screen, f.Message.X, f.Message.Y, f.Message.Text, f.Message.Hidden, theme.MessageText, theme.Background)
		}
	}

	f.ParticleSystem.Update(dt)
	if f.Debug != nil {
		f.Debug.UpdateRuntime(dt, len(f.ParticleSystem.GetParticles()), f.MaxParticles, f.TargetFPS)
	}
	if f.PlayerEnabled {
		f.hitPlayer()
	} else if f.MessageErasable {
		f.ParticleSystem.HitMessage(f.Message.Text, f.Message.X, f.Message.Y, f.Message.Hidden)
	}
	if f.PlayerEnabled && f.Player != nil && !f.RainInFront {
		f.ParticleSystem.Draw(f.Screen)
		f.renderPlayer()
	} else {
		if f.PlayerEnabled && f.Player != nil {
			f.renderPlayer()
		}
		f.ParticleSystem.Draw(f.Screen)
	}
	if f.Debug != nil {
		f.Debug.DrawOverlay(f.Screen)
	}

	f.Screen.Show()
}

func (f Frame) renderPlayer() {
	p := f.Player

	if p.Artwork != nil {
		p.Artwork.Render(f.Screen)
	}
}
