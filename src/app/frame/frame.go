package frame

import (
	"fmt"
	"virga-player/app/message"
	"virga-player/app/player"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

type Frame struct {
	Screen          tcell.Screen
	Renderer        *renderer.Renderer
	ParticleSystem  *rain.ParticleSystem
	Message         *message.Message
	Player          *player.Player
	PlayerEnabled   bool
	MessageErasable bool
}

func NewFrame(screen tcell.Screen, renderer *renderer.Renderer, particles *rain.ParticleSystem, msg *message.Message, p *player.Player, playerEnabled bool, messageErasable bool) Frame {
	return Frame{
		Screen:          screen,
		Renderer:        renderer,
		ParticleSystem:  particles,
		Message:         msg,
		Player:          p,
		PlayerEnabled:   playerEnabled,
		MessageErasable: messageErasable,
	}
}

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

func (f Frame) hitPlayer() {
	if f.Player == nil || !f.PlayerEnabled {
		return
	}

	for _, part := range f.ParticleSystem.GetParticles() {
		for i := 0; i < part.Length; i++ {
			dropX := int(part.X)
			dropY := int(part.Y) + i

			// Hit album cover
			if dropX >= f.Player.CoverX && dropX < f.Player.CoverX+f.Player.CoverW &&
				dropY >= f.Player.CoverY && dropY < f.Player.CoverY+f.Player.CoverH {
				f.Player.MarkHit("cover", dropX, dropY)
			}

			// Hit title
			titleY := f.Player.TextY
			if dropY == titleY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Title) {
				f.Player.MarkHit("title", dropX, titleY)
			}

			// Hit artist
			artistY := f.Player.TextY + 1
			if dropY == artistY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Artist) {
				f.Player.MarkHit("artist", dropX, artistY)
			}

			// Hit album
			albumY := f.Player.TextY + 2
			if dropY == albumY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Album) {
				f.Player.MarkHit("album", dropX, albumY)
			}

			// Hit duration
			durationY := f.Player.TextY + 3
			progress := fmt.Sprintf("%d:%02d / %d:%02d", f.Player.Elapsed/60, f.Player.Elapsed%60, f.Player.Duration/60, f.Player.Duration%60)
			if dropY == durationY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(progress) {
				f.Player.MarkHit("duration", dropX, durationY)
			}
		}
	}
}
