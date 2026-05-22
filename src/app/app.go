package app

import (
	"time"

	"virga-player/animation"
	"virga-player/app/state"
	"virga-player/audio"
	"virga-player/debug"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/settings/page"

	"github.com/gdamore/tcell/v2"
)

type App struct {
	screen         tcell.Screen
	particleSystem *rain.ParticleSystem
	animEngine     *animation.Engine
	renderEngine   *renderer.Renderer
	state          *state.AppState
	settingsPage   *page.Page
	cfg            *settings.Config
	eventChan      <-chan tcell.Event

	width        int
	height       int
	lastTick     time.Time
	settingsOpen bool
	exitAt       time.Time

	audioAnalyzer *audio.Analyzer
	debug         *debug.Manager
	debugForced   bool
}

type Options struct {
	Debug bool
}
