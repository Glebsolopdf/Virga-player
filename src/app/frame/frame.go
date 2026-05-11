package frame

import (
	"virga-player/app/message"
	"virga-player/app/player"
	"virga-player/rain"
	"virga-player/renderer"

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
