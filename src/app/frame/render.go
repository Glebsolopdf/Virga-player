package frame

import (
	"fmt"

	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (f Frame) Render(dt float64) {
	theme := settings.CurrentTheme()
	f.Screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.MessageText))
	f.Screen.Clear()

	if f.PlayerEnabled && f.Player != nil {
		f.renderPlayer()
	} else {
		if !f.Message.Converted {
			f.Renderer.DrawTextMasked(f.Screen, f.Message.X, f.Message.Y, f.Message.Text, f.Message.Hidden, theme.MessageText, theme.Background)
		}
	}

	f.ParticleSystem.Update(dt)
	if f.PlayerEnabled {
		f.hitPlayer()
	} else if f.MessageErasable {
		f.ParticleSystem.HitMessage(f.Message.Text, f.Message.X, f.Message.Y, f.Message.Hidden)
	}
	f.ParticleSystem.Draw(f.Screen)

	f.Screen.Show()
}

func (f Frame) renderPlayer() {
	p := f.Player

	if p.Artwork != nil {
		p.Artwork.Render(f.Screen)
	}

	for _, part := range f.ParticleSystem.GetParticles() {
		for i := 0; i < part.Length; i++ {
			dropX := int(part.X)
			dropY := int(part.Y) + i

			// Hit title
			titleY := p.TextY
			if dropY == titleY && dropX >= p.TextX && dropX < p.TextX+len(p.Title) {
				p.MarkHit("title", dropX, titleY)
			}

			// Hit artist
			artistY := p.TextY + 1
			if dropY == artistY && dropX >= p.TextX && dropX < p.TextX+len(p.Artist) {
				p.MarkHit("artist", dropX, artistY)
			}

			// Hit album
			albumY := p.TextY + 2
			if dropY == albumY && dropX >= p.TextX && dropX < p.TextX+len(p.Album) {
				p.MarkHit("album", dropX, albumY)
			}

			// Hit duration
			timeStr := fmt.Sprintf("%d:%02d / %d:%02d", p.Elapsed/60, p.Elapsed%60, p.Duration/60, p.Duration%60)
			durationY := p.TextY + 3
			if dropY == durationY && dropX >= p.TextX && dropX < p.TextX+len(timeStr) {
				p.MarkHit("duration", dropX, durationY)
			}
		}
	}
}
