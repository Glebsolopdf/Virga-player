package app

import (
	"time"

	"virga-player/animation"
	"virga-player/app/events"
	"virga-player/app/install"
	"virga-player/app/state"
	"virga-player/audio"
	"virga-player/notification"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/settings/page"
	"virga-player/version"
)

const defaultStatusMessage = "ESC to exit | S - settings"

func (a *App) initComponents() {
	a.width, a.height = a.screen.Size()
	cfg, firstRun, err := settings.LoadOrCreateConfig()
	if err != nil {
		if a.debug != nil {
			a.debug.Errorf("config load failed, using defaults: %v", err)
		}
		cfg = settings.DefaultConfig()
	}
	if a.debugForced {
		cfg.Debug = true
	}
	theme, _, themeErr := settings.LoadOrCreateTheme()
	if themeErr != nil {
		if a.debug != nil {
			a.debug.Errorf("theme load failed, using default: %v", themeErr)
		}
		theme = settings.DefaultTheme()
	}
	if a.debug != nil {
		a.debug.SetEnabled(cfg.Debug)
		a.debug.Infof("config loaded: fps=%d particles=%d debug=%v lyrics_visible=%v lyrics_mode=%s lyrics_auto_save_after=%ds", cfg.FPS, cfg.MaxParticles, cfg.Debug, cfg.LyricsVisible, cfg.LyricsMode, cfg.LyricsAutoSaveAfterSec)
	}
	a.cfg = cfg
	settings.SetCurrentTheme(theme)
	a.setupNotifications(firstRun)
	aliasesReady := install.EnsureCommandAliases()
	messageText := defaultStatusMessage
	nextMessageText := ""
	if firstRun {
		messageText = "Welcome to Virga!"
		if aliasesReady {
			messageText = "Welcome to Virga! PATH: virga | virgaplayer"
		}
		nextMessageText = defaultStatusMessage
	}
	a.resetLyricsManager()
	a.particleSystem = rain.NewParticleSystem(a.width, a.height, a.cfg)
	a.setupAudioAnalyzer()
	a.animEngine = animation.NewEngine(a.cfg.FPS)
	a.renderEngine = renderer.NewRenderer(a.screen)
	a.state = state.NewAppState(a.width, a.height, messageText, nextMessageText, a.cfg)
	a.settingsPage = page.NewPage(a.cfg.Clone())
	a.settingsPage.SetNotificationsSupported(a.notificationsSupported)
	a.settingsPage.SetNotifications(a.notificationItems(), a.openNotificationsPage)
}

func (a *App) setupNotifications(firstRun bool) {
	if a.uninstallInProgress.Load() {
		return
	}

	service, err := notification.Load()
	if err != nil {
		if a.debug != nil {
			a.debug.Warnf("notification state load failed, using empty state: %v", err)
		}
		service = notification.NewForPath(notification.StatePath())
	}
	a.notifications = service

	if supported, known, stateErr := a.notifications.RemoteSupportState(); stateErr != nil {
		if a.debug != nil {
			a.debug.Warnf("notification support-state load failed: %v", stateErr)
		}
	} else if known {
		a.notificationsSupported = supported
		if !supported && a.cfg != nil {
			a.cfg.NotificationsEnabled = false
		}
	}

	a.startRemoteNotificationSync()

	if a.cfg != nil && !a.cfg.NotificationsEnabled && a.notificationsSupported {
		if err := a.notifications.ClearAll(); err != nil {
			if a.debug != nil {
				a.debug.Warnf("notifications clear failed: %v", err)
			}
		}
		if a.notificationToast != nil {
			a.notificationToast.Hide()
		}
		a.refreshNotificationsPageBindings()
		return
	}

	if removed, err := a.notifications.RemoveUpdateNotices(); err != nil {
		if a.debug != nil {
			a.debug.Warnf("legacy update notifications cleanup failed: %v", err)
		}
	} else if removed && a.debug != nil {
		a.debug.Infof("legacy update notifications removed")
	}

	if _, added, err := a.notifications.EnsureNotificationsIntro(); err != nil {
		if a.debug != nil {
			a.debug.Warnf("notifications intro init failed: %v", err)
		}
	} else if added && a.debug != nil {
		a.debug.Infof("notifications intro added")
	}

	if _, added, err := a.notifications.EnsureWelcome(firstRun, version.AppVersion); err != nil {
		if a.debug != nil {
			a.debug.Warnf("welcome notification init failed: %v", err)
		}
	} else if added && a.debug != nil {
		a.debug.Infof("welcome notification added")
	}

	if firstRun {
		go func() {
			time.Sleep(7 * time.Second)
			a.armNotificationToast()
		}()
	} else {
		a.armNotificationToast()
	}

	a.refreshNotificationsPageBindings()
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
