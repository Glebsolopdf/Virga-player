package frame

import (
	"virga-player/app/message"
	"virga-player/app/player"
	debugmgr "virga-player/debug/manager"
	"virga-player/notification"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

type Frame struct {
	Screen            tcell.Screen
	Renderer          *renderer.Renderer
	ParticleSystem    *rain.ParticleSystem
	Message           *message.Message
	Player            *player.Player
	PlayerEnabled     bool
	MessageErasable   bool
	Debug             *debugmgr.Manager
	MaxParticles      int
	TargetFPS         int
	PlayerRainLayer   settings.RainLayerMode
	LyricsRainLayer   settings.RainLayerMode
	NotificationToast *notification.Toast
	FooterPromptText  string
}

func NewFrame(screen tcell.Screen, renderer *renderer.Renderer, particles *rain.ParticleSystem, msg *message.Message, p *player.Player, playerEnabled bool, messageErasable bool, dbg *debugmgr.Manager, maxParticles int, targetFPS int, playerRainLayer, lyricsRainLayer settings.RainLayerMode, notificationToast *notification.Toast, footerPromptText string) Frame {
	return Frame{
		Screen:            screen,
		Renderer:          renderer,
		ParticleSystem:    particles,
		Message:           msg,
		Player:            p,
		PlayerEnabled:     playerEnabled,
		MessageErasable:   messageErasable,
		Debug:             dbg,
		MaxParticles:      maxParticles,
		TargetFPS:         targetFPS,
		PlayerRainLayer:   playerRainLayer,
		LyricsRainLayer:   lyricsRainLayer,
		NotificationToast: notificationToast,
		FooterPromptText:  footerPromptText,
	}
}
