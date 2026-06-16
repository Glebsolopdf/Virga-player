package app

import (
	"errors"

	"virga-player/animation"
	"virga-player/app/events"
	"virga-player/app/install"
	"virga-player/app/state"
	"virga-player/audio"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/settings/page"
	"virga-player/welcome"
)

const defaultStatusMessage = "ESC to exit | S - settings"

func (a *App) prepareConfig() error {
	cfg, firstRun, err := settings.LoadConfig()
	if err != nil {
		if a.debug != nil {
			a.debug.Errorf("config load failed, using defaults: %v", err)
		}
		cfg = settings.DefaultConfig()
	}
	if firstRun {
		cfg, err = welcome.Run(welcome.Options{
			DefaultConfig: settings.DefaultConfig(),
			PlayerConfig:  settings.PlayerConfig(),
		})
		if err != nil {
			if errors.Is(err, welcome.ErrCanceled) {
				return err
			}
			return err
		}
		if err := settings.SaveConfig(cfg); err != nil {
			return err
		}
	}
	if a.debugForced {
		cfg.Debug = true
	}
	a.cfg = cfg
	return nil
}

func (a *App) initComponents() {
	a.width, a.height = a.screen.Size()
	theme, _, themeErr := settings.LoadOrCreateTheme()
	if themeErr != nil {
		if a.debug != nil {
			a.debug.Errorf("theme load failed, using default: %v", themeErr)
		}
		theme = settings.DefaultTheme()
	}
	if a.debug != nil {
		a.debug.SetEnabled(a.cfg.Debug)
		a.debug.Infof("config loaded: fps=%d particles=%d debug=%v lyrics_visible=%v lyrics_mode=%s lyrics_auto_save_after=%ds", a.cfg.FPS, a.cfg.MaxParticles, a.cfg.Debug, a.cfg.LyricsVisible, a.cfg.LyricsMode, a.cfg.LyricsAutoSaveAfterSec)
	}
	settings.SetCurrentTheme(theme)
	_ = install.EnsureCommandAliases()
	messageText := defaultStatusMessage
	a.resetLyricsManager()
	a.particleSystem = rain.NewParticleSystem(a.width, a.height, a.cfg)
	a.setupAudioAnalyzer()
	a.animEngine = animation.NewEngine(a.cfg.FPS)
	a.renderEngine = renderer.NewRenderer(a.screen)
	a.state = state.NewAppState(a.width, a.height, messageText, "", a.cfg)
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

	if !a.cfg.MusicReactive && !a.cfg.RainVisualizer {
		a.particleSystem.ResetSpectrum()
		if a.debug != nil {
			a.debug.Debugf("audio analyzer disabled: music_reactive=false rain_visualizer=false")
		}
		return
	}

	analyzer, err := audio.NewAnalyzer()
	if err != nil {
		a.particleSystem.ResetSpectrum()
		if a.debug != nil {
			a.debug.Warnf("audio analyzer unavailable: %v", err)
		}
		return
	}
	a.audioAnalyzer = analyzer
	if a.debug != nil {
		a.debug.Infof("audio analyzer connected")
	}
}
