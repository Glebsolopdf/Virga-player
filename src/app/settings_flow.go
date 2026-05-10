package app

import (
	"virga-player/animation"
	"virga-player/app/player"
	"virga-player/music"
	"virga-player/settings"
	"virga-player/settings/page"
)

func (a *App) openSettings() {
	a.settingsOpen = true
	a.settingsPage = page.NewPage(a.cfg.Clone())
}

func (a *App) closeSettings(save bool, deleteVirga bool) {
	if deleteVirga {
		_ = removeVirgaInstallation()
		a.cfg = settings.DefaultConfig()
		a.applyConfig()
		a.settingsOpen = false
		return
	}

	if save {
		a.cfg = a.settingsPage.Config
		_ = settings.SaveConfig(a.cfg)
		a.applyConfig()
	}
	a.settingsOpen = false
}

func (a *App) applyConfig() {
	a.animEngine.Stop()
	a.animEngine = animation.NewEngine(a.cfg.FPS)
	a.particleSystem.ApplyConfig(a.cfg)
	a.setupAudioAnalyzer()

	if a.cfg.Player && !a.state.PlayerEnabled {
		a.state.Player = player.New(a.width/2, a.height/2)
		a.state.Player.UpdatePosition(a.width, a.height)
		a.state.PlayerEnabled = true

		track := music.GetTrackInfo()
		artworkPath := track.GetArtworkPath()
		a.state.Player.SetTrackInfoWithArtwork(
			track.Title,
			track.Artist,
			track.Album,
			track.DurationFormatted(),
			track.ElapsedFormatted(),
			track.Duration,
			track.Elapsed,
			artworkPath,
		)
	} else if !a.cfg.Player && a.state.PlayerEnabled {
		a.state.PlayerEnabled = false
		a.state.Player = &player.Player{}
	}
}
