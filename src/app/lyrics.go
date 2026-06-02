package app

import (
	"errors"
	"strconv"
	"strings"

	"virga-player/lyricsearch"
	"virga-player/music"
	"virga-player/settings"
)

type lyricFetchResult struct {
	requestKey string
	artist     string
	track      string
	lyrics     string
	err        error
}

func (a *App) processLyricsResults() {
	if a.lyricsResults == nil {
		return
	}

	for {
		select {
		case result := <-a.lyricsResults:
			if result.requestKey != a.lyricsRequestKey {
				if a.debug != nil {
					a.debug.Debugf("lyrics result ignored as stale artist=%q track=%q key=%q current=%q", result.artist, result.track, result.requestKey, a.lyricsRequestKey)
				}
				continue
			}

			a.lyricsRequestKey = ""
			a.lyricsResultKey = result.requestKey
			a.currentLyrics = result.lyrics
			if a.state != nil && a.state.Player != nil {
				a.state.Player.SetSyncedLyrics(result.lyrics)
			}

			if result.lyrics != "" && a.debug != nil {
				a.debug.Infof("lyrics ready artist=%q track=%q bytes=%d", result.artist, result.track, len(result.lyrics))
			}

			if result.err == nil || a.debug == nil {
				continue
			}

			switch {
			case errors.Is(result.err, lyricsearch.ErrInstrumentalTrack):
				a.debug.Debugf("lyrics skipped instrumental artist=%q track=%q", result.artist, result.track)
			case errors.Is(result.err, lyricsearch.ErrLyricsNotFound):
				a.debug.Debugf("lyrics not found artist=%q track=%q", result.artist, result.track)
			case errors.Is(result.err, lyricsearch.ErrMissingMetadata):
				a.debug.Debugf("lyrics skipped missing metadata artist=%q track=%q", result.artist, result.track)
			case errors.Is(result.err, lyricsearch.ErrLyricsDisabled):
				a.debug.Debugf("lyrics disabled")
			case errors.Is(result.err, lyricsearch.ErrLyricsManagerClosed):
				a.debug.Debugf("lyrics manager closed")
			default:
				a.debug.Warnf("lyrics lookup failed artist=%q track=%q err=%v", result.artist, result.track, result.err)
			}
		default:
			return
		}
	}
}

func (a *App) syncLyrics(track *music.TrackInfo) {
	if a.cfg == nil || !a.cfg.Player {
		a.clearLyricsState("player-disabled")
		return
	}
	if !a.cfg.LyricsVisible {
		a.clearLyricsState("hidden")
		return
	}

	if a.cfg.LyricsMode == settings.LyricsModeDisabled {
		a.clearLyricsState("disabled")
		return
	}
	if a.lyricsManager == nil {
		a.clearLyricsState("manager-unavailable")
		return
	}

	artist, title, requestKey := a.lyricsRequestFor(track)
	if requestKey == "" {
		a.clearLyricsState("no-track")
		return
	}

	if requestKey == a.lyricsRequestKey || requestKey == a.lyricsResultKey {
		return
	}

	a.lyricsRequestKey = requestKey
	a.lyricsResultKey = ""
	a.currentLyrics = ""
	if a.state != nil && a.state.Player != nil {
		a.state.Player.SetSyncedLyrics("")
	}

	if a.debug != nil {
		a.debug.Infof("lyrics lookup started mode=%s artist=%q track=%q", a.cfg.LyricsMode, artist, title)
	}

	lyricsMgr := a.lyricsManager
	go func(requestKey, artist, title string, mgr *lyricsearch.LyricsManager, results chan<- lyricFetchResult) {
		lyrics, err := mgr.OnTrackStarted(lyricsearch.Track{Artist: artist, Title: title})
		results <- lyricFetchResult{
			requestKey: requestKey,
			artist:     artist,
			track:      title,
			lyrics:     lyrics,
			err:        err,
		}
	}(requestKey, artist, title, lyricsMgr, a.lyricsResults)
}

func (a *App) lyricsRequestFor(track *music.TrackInfo) (artist, title, requestKey string) {
	if track == nil || track.Source == "default" {
		return "", "", ""
	}

	artist = strings.TrimSpace(track.Artist)
	title = strings.TrimSpace(track.Title)
	if artist == "" || title == "" {
		return "", "", ""
	}

	requestKey = string(a.cfg.LyricsMode) +
		"\x00" + artist +
		"\x00" + title +
		"\x00" + strconv.Itoa(a.cfg.LyricsAutoSaveAfterSec) +
		"\x00" + strconv.FormatBool(a.cfg.LyricsDoubleConfirm)
	return artist, title, requestKey
}

func (a *App) clearLyricsState(reason string) {
	if a.lyricsRequestKey == "" && a.lyricsResultKey == "" && a.currentLyrics == "" {
		return
	}

	a.lyricsRequestKey = ""
	a.lyricsResultKey = ""
	a.currentLyrics = ""
	if a.state != nil && a.state.Player != nil {
		a.state.Player.SetSyncedLyrics("")
	}

	if a.debug != nil {
		a.debug.Debugf("lyrics state cleared reason=%s", reason)
	}
}
