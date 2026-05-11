package app

import (
	"time"

	"virga-player/app/frame"
	"virga-player/music"
)

func (a *App) onTick(dt float64) {
	coverPulse := 0.0
	if a.audioAnalyzer != nil {
		bands, ok := a.audioAnalyzer.Bands()
		if ok {
			a.particleSystem.ApplySpectrum(bands.Low, bands.Mid, bands.High, bands.Envelope)
			coverPulse = bands.Envelope
		} else {
			a.particleSystem.ResetSpectrum()
		}
	} else {
		a.particleSystem.ResetSpectrum()
	}

	if a.state.PlayerEnabled && a.state.Player != nil {
		track := music.GetTrackInfo()
		artworkPath := a.state.Player.ArtworkPath
		if track.ArtworkURL != a.state.Player.ArtworkURL {
			resolvedPath := track.GetArtworkPath()
			artworkPath = resolvedPath
			if resolvedPath == "" && a.state.Player.ArtworkPath != "" {
				artworkPath = a.state.Player.ArtworkPath
			}
			a.state.Player.ArtworkURL = track.ArtworkURL
		}

		if track.Title != a.state.Player.Title || track.Artist != a.state.Player.Artist || track.Album != a.state.Player.Album || artworkPath != a.state.Player.ArtworkPath {
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
		} else {
			a.state.Player.SetTrackInfo(
				track.Title,
				track.Artist,
				track.Album,
				track.DurationFormatted(),
				track.ElapsedFormatted(),
				track.Duration,
				track.Elapsed,
			)
		}

		if a.state.Player.Artwork != nil {
			a.state.Player.Artwork.SetAnimationEnabled(a.cfg.CoverAnimation)
			a.state.Player.Artwork.UpdateAnimation(dt, coverPulse)
		}
	}

	if a.state.AdvanceIntroIfNeeded() {
		a.state.TriggerWashNow()
	}

	a.particleSystem.SetSpawnPaused(a.state.ShouldPauseRain())

	if !a.state.IsIntroActive() && !a.state.WashTriggered && a.state.SinceStart() >= 3*time.Second {
		a.state.TriggerWash()
	}

	if a.state.ShouldSpawnMessageDrops() {
		a.particleSystem.SpawnMessageDrops(a.state.Message.X, a.state.Message.Y, a.state.Message.Text, a.state.Message.Hidden)
		if !a.state.AdvanceMessage() {
			a.state.Message.Converted = true
		}
	}

	if a.settingsOpen {
		a.particleSystem.Update(dt)
		if !a.state.IsMessageProtected() {
			a.particleSystem.HitMessage(a.state.Message.Text, a.state.Message.X, a.state.Message.Y, a.state.Message.Hidden)
		}
		a.particleSystem.Draw(a.screen)
		a.settingsPage.Render(a.screen, a.renderEngine, a.width, a.height)
		a.screen.Show()
		return
	}

	frame.NewFrame(
		a.screen,
		a.renderEngine,
		a.particleSystem,
		a.state.Message,
		a.state.Player,
		a.state.PlayerEnabled,
		!a.state.IsMessageProtected() && !a.state.Message.Persistent,
	).Render(dt)
}
