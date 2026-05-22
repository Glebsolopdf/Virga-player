package app

import (
	"strings"
	"time"

	"virga-player/app/frame"
	"virga-player/music"
)

func sanitizeArtworkURLForLog(raw string) string {
	if strings.HasPrefix(strings.ToLower(raw), "data:") {
		return "data:<redacted>"
	}
	return raw
}

func (a *App) onTick(dt float64) {
	coverPulse := 0.0
	coverLow := 0.0
	coverMid := 0.0
	coverHigh := 0.0
	musicPressure := 0.0
	if a.audioAnalyzer != nil {
		bands, ok := a.audioAnalyzer.Bands()
		if ok {
			a.particleSystem.ApplySpectrum(dt, bands.Low, bands.Mid, bands.High, bands.Envelope)
			coverLow = bands.Low
			coverMid = bands.Mid
			coverHigh = bands.High
			coverPulse = bands.Envelope
			audioPressure := bands.Low*0.55 + bands.Mid*0.35 + bands.High*0.20
			if audioPressure > 1 {
				audioPressure = 1
			}
			musicPressure = coverPulse*0.65 + audioPressure*0.55
			musicPressure *= float64(a.cfg.MusicPlayerIntensity) / 100.0
			if musicPressure > 1 {
				musicPressure = 1
			}
		} else {
			a.particleSystem.ResetSpectrum()
		}
	} else {
		a.particleSystem.ResetSpectrum()
	}

	if a.state.PlayerEnabled && a.state.Player != nil {
		track := music.GetTrackInfo()
		artworkPath := track.ArtworkPath
		if artworkPath == "" && a.state.Player.ArtworkPath != "" {
			artworkPath = a.state.Player.ArtworkPath
		}
		if track.ArtworkURL != a.state.Player.ArtworkURL {
			a.state.Player.ArtworkURL = track.ArtworkURL
			if a.debug != nil {
				a.debug.Infof("artwork source=%s artworkURL=%q artworkPath=%q fallbackPath=%q", track.Source, sanitizeArtworkURLForLog(track.ArtworkURL), track.ArtworkPath, a.state.Player.ArtworkPath)
			}
		}

		if a.debug != nil && (track.Title != a.state.Player.Title || track.Artist != a.state.Player.Artist || track.Album != a.state.Player.Album || artworkPath != a.state.Player.ArtworkPath) {
			a.debug.Infof("track update source=%s title=%q artist=%q album=%q artworkURL=%q artworkPath=%q", track.Source, track.Title, track.Artist, track.Album, sanitizeArtworkURLForLog(track.ArtworkURL), artworkPath)
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
			a.state.Player.Artwork.SetAnimationEnabled(a.cfg.PulseOnCover())
			a.state.Player.Artwork.UpdateAnimation(dt, coverLow, coverMid, coverHigh, coverPulse, float64(a.cfg.PulseSpeed)/100.0)
			if a.cfg.MusicPlayerAnimation {
				a.state.Player.Artwork.UpdateRainResistance(dt, musicPressure, a.cfg.MusicPlayerInvert)
			} else {
				a.state.Player.Artwork.UpdateRainResistance(dt, 0, false)
			}
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
		if a.debug != nil {
			a.debug.UpdateRuntime(dt, len(a.particleSystem.GetParticles()), a.cfg.MaxParticles, a.cfg.FPS)
		}
		if !a.state.IsMessageProtected() {
			a.particleSystem.HitMessage(a.state.Message.Text, a.state.Message.X, a.state.Message.Y, a.state.Message.Hidden)
		}
		a.particleSystem.Draw(a.screen)
		a.settingsPage.Render(a.screen, a.renderEngine, a.width, a.height)
		if a.debug != nil {
			a.debug.DrawOverlay(a.screen)
		}
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
		a.debug,
		a.cfg.MaxParticles,
		a.cfg.FPS,
		a.cfg.RainInFrontOfPlayer,
	).Render(dt)
}
