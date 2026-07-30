package frame

import (
	"virga-player/app/message"
	"virga-player/app/player"
	debugmgr "virga-player/debug/manager"
	"virga-player/rain"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

type Frame struct {
	Screen            tcell.Screen
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
	FooterPromptText  string
}


