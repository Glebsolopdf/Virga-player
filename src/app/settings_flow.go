package app

import (
	"time"

	"virga-player/animation"
	"virga-player/app/install"
	"virga-player/app/player"
	"virga-player/music"
	"virga-player/settings"
	"virga-player/settings/page"
)

func (a *App) openSettings() {
	a.settingsOpen = true
	a.settingsPage = page.NewPage(a.cfg.Clone())
	a.settingsPage.SetNotificationsSupported(a.notificationsSupported)
	a.settingsPage.SetNotifications(a.notificationItems(), a.openNotificationsPage)
}

func (a *App) closeSettings(save bool, deleteVirga bool) bool {
	if deleteVirga {
		a.uninstallInProgress.Store(true)
		a.notifications = nil
		if a.notificationToast != nil {
			a.notificationToast.Hide()
		}
		_ = install.RemoveVirgaInstallation()
		a.state.Message.SetText("Virga removed. Restart your shell or run 'hash -r' to refresh command lookup.", a.width, a.height)
		a.state.Message.Persistent = true
		a.settingsOpen = false
		a.exitAt = time.Now().Add(7 * time.Second)
		return false
	}

	if save {
		a.cfg = a.settingsPage.Config
		if a.debugForced {
			a.cfg.Debug = true
		}
		_ = settings.SaveConfig(a.cfg)
		a.applyConfig()
		if a.debug != nil {
			a.debug.Infof("settings saved")
		}
	}
	a.settingsOpen = false
	return false
}

func (a *App) applyConfig() {
	a.animEngine.Stop()
	a.animEngine = animation.NewEngine(a.cfg.FPS)
	a.particleSystem.ApplyConfig(a.cfg)
	if a.debug != nil {
		a.debug.SetEnabled(a.cfg.Debug)
		a.debug.Debugf("config applied: fps=%d speed=%d debug=%v lyrics_visible=%v lyrics_mode=%s lyrics_auto_save_after=%ds", a.cfg.FPS, a.cfg.RainSpeed, a.cfg.Debug, a.cfg.LyricsVisible, a.cfg.LyricsMode, a.cfg.LyricsAutoSaveAfterSec)
	}
	a.resetLyricsManager()
	a.setupNotifications(false)
	a.setupAudioAnalyzer()
	if !a.cfg.LyricsVisible {
		a.clearLyricsState("hidden")
	}

	if a.cfg.Player && !a.state.PlayerEnabled {
		a.state.Player = player.New(a.width/2, a.height/2)
		a.state.Player.UpdatePosition(a.width, a.height)
		a.state.PlayerEnabled = true

		track := music.GetTrackInfo()
		artworkPath := track.ArtworkPath
		if artworkPath == "" {
			artworkPath = track.GetArtworkPath()
		}
		a.state.Player.SetTrackInfoWithArtwork(
			track.Title,
			track.Artist,
			track.Album,
			track.Duration,
			track.Elapsed,
			artworkPath,
		)
		a.state.Player.SetSyncedLyrics(a.currentLyrics)
		a.state.Player.ArtworkURL = track.ArtworkURL
	} else if !a.cfg.Player && a.state.PlayerEnabled {
		a.state.PlayerEnabled = false
		a.state.Player = &player.Player{}
	}
}
