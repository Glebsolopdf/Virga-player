package app

import (
	"virga-player/animation"
	"virga-player/app/events"
	"virga-player/app/state"
	"virga-player/audio"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/settings/page"
)

const defaultStatusMessage = "ESC to exit | S - settings"

func (a *App) initComponents() {
	a.width, a.height = a.screen.Size()
	cfg, firstRun, err := settings.LoadOrCreateConfig()
	if err != nil {
		cfg = settings.DefaultConfig()
	}
	theme, _, themeErr := settings.LoadOrCreateTheme()
	if themeErr != nil {
		theme = settings.DefaultTheme()
	}
	settings.SetCurrentTheme(theme)
	aliasesReady := ensureCommandAliases()
	messageText := defaultStatusMessage
	nextMessageText := ""
	if firstRun {
		messageText = "Welcome to Virga!"
		if aliasesReady {
			messageText = "Welcome to Virga! PATH: virga | virgaplayer"
		}
		nextMessageText = defaultStatusMessage
	}
	a.cfg = cfg
	a.particleSystem = rain.NewParticleSystem(a.width, a.height, a.cfg)
	a.setupAudioAnalyzer()
	a.animEngine = animation.NewEngine(a.cfg.FPS)
	a.renderEngine = renderer.NewRenderer(a.screen)
	a.state = state.NewAppState(a.width, a.height, messageText, nextMessageText, a.cfg)
	a.settingsPage = page.NewPage(a.cfg.Clone())
}

func (a *App) initEvents() {
	a.eventChan = events.Start(a.screen)
}

func (a *App) setupAudioAnalyzer() {
	if a.audioAnalyzer != nil {
		a.audioAnalyzer.Stop()
		a.audioAnalyzer = nil
	}

	if !a.cfg.MusicReactive {
		a.particleSystem.ResetSpectrum()
		return
	}

	analyzer, err := audio.NewAnalyzer()
	if err != nil {
		a.particleSystem.ResetSpectrum()
		return
	}
	a.audioAnalyzer = analyzer
}
