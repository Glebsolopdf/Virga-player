package frame

import (
	"virga-player/app/message"
	"virga-player/app/player"
	"virga-player/app/playerlyrics"
	"virga-player/audio"
	debugmgr "virga-player/debug/manager"
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
	PlayerPosition    settings.PositionPreset
	LyricsPosition    settings.PositionPreset
	LyricsDisplayMode settings.LyricsDisplayMode
	Lyrics            *playerlyrics.State
	AudioBands        audio.Bands
	FooterPromptText  string
}

func NewFrame(screen tcell.Screen, renderer *renderer.Renderer, particles *rain.ParticleSystem, msg *message.Message, p *player.Player, playerEnabled bool, messageErasable bool, dbg *debugmgr.Manager, maxParticles int, targetFPS int, playerRainLayer, lyricsRainLayer settings.RainLayerMode, playerPosition, lyricsPosition settings.PositionPreset, lyricsMode settings.LyricsDisplayMode, lyrics *playerlyrics.State, bands audio.Bands, footerPromptText string) Frame {
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
		PlayerPosition:    playerPosition,
		LyricsPosition:    lyricsPosition,
		LyricsDisplayMode: lyricsMode,
		Lyrics:            lyrics,
		AudioBands:        bands,
		FooterPromptText:  footerPromptText,
	}
}
