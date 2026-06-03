package lyricsmanager

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (m *LyricsManager) OnTrackStarted(track Track) (string, error) {
	artist := strings.TrimSpace(track.Artist)
	title := strings.TrimSpace(track.Title)
	if artist == "" || title == "" {
		m.cancelActiveTrack()
		return "", m.cfg.ErrMissingMetadata
	}

	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return "", ErrLyricsManagerClosed
	}
	mode := m.cfg.Mode
	tempDir := m.cfg.TempDir
	persistentDir := m.cfg.PersistentDir
	autoSaveAfter := m.cfg.AutoSaveAfter
	promptCb := m.cfg.Prompt
	m.mu.RUnlock()

	if mode == m.cfg.ModeDisabled {
		m.cancelActiveTrack()
		return "", m.cfg.ErrLyricsDisabled
	}

	key := artist + "\x00" + title
	m.setActiveTrack(key, nil)

	if usesRAMStorage(m.cfg, mode) {
		if lyrics, ok := m.getRAM(key); ok {
			m.debugf("ram cache hit artist=%q track=%q", artist, title)
			m.scheduleDeferredSave(mode, key, autoSaveAfter, promptCb, track, lyrics)
			return lyrics, nil
		}
		m.debugf("ram cache miss artist=%q track=%q", artist, title)
	}

	if m.cfg.ReadLyricsFromDir != nil {
		if lyrics, err := m.cfg.ReadLyricsFromDir(persistentDir, artist, title); err == nil {
			if m.cfg.HasTimedLyrics == nil || m.cfg.HasTimedLyrics(lyrics) {
				if usesRAMStorage(m.cfg, mode) {
					if err := m.putRAM(key, artist, title, lyrics, tempDir); err != nil {
						m.debugf("failed to write temp cache artist=%q track=%q err=%v", artist, title, err)
					}
				}
				return lyrics, nil
			}
		}
	}

	if m.cfg.FetchLyrics == nil {
		return "", fmt.Errorf("fetch lyrics callback is not configured")
	}
	lyrics, err := m.cfg.FetchLyrics(artist, title)
	if err != nil {
		return "", err
	}

	if mode == m.cfg.ModeDirectToDisk {
		if m.cfg.WriteLyricsToDir == nil {
			return lyrics, fmt.Errorf("write lyrics callback is not configured")
		}
		if _, err := m.cfg.WriteLyricsToDir(persistentDir, artist, title, lyrics); err != nil {
			return lyrics, fmt.Errorf("save lyrics to disk: %w", err)
		}
		m.debugf("direct save complete artist=%q track=%q bytes=%d", artist, title, len(lyrics))
		return lyrics, nil
	}

	if usesRAMStorage(m.cfg, mode) {
		if err := m.putRAM(key, artist, title, lyrics, tempDir); err != nil {
			return lyrics, fmt.Errorf("save lyrics to temporary cache: %w", err)
		}
	}

	m.scheduleDeferredSave(mode, key, autoSaveAfter, promptCb, track, lyrics)
	return lyrics, nil
}

func (m *LyricsManager) scheduleDeferredSave(mode, key string, delay time.Duration, prompt PromptCallback, track Track, lyrics string) {
	if mode != m.cfg.ModeRAMWithAuto && mode != m.cfg.ModeRAMWithPrompt {
		return
	}

	if !m.isActiveTrack(key) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.setActiveTrack(key, cancel)

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		if !m.isActiveTrack(key) {
			return
		}

		if m.hasPersistentTimedLyrics(track) {
			m.debugf("deferred save skipped (already persisted) artist=%q track=%q", track.Artist, track.Title)
			if mode == m.cfg.ModeRAMWithPrompt {
				m.notifyPromptStatus(prompt, track, lyrics, "Lyrics already saved")
			}
			return
		}

		if mode == m.cfg.ModeRAMWithPrompt {
			if prompt == nil {
				return
			}
			approved := prompt(ctx, PromptRequest{
				Track:  track,
				Lyrics: lyrics,
				Message: fmt.Sprintf(
					"Save lyrics to disk for %s - %s? (Y/N)",
					strings.TrimSpace(track.Artist),
					strings.TrimSpace(track.Title),
				),
			})
			if !approved {
				return
			}
		}

		m.mu.RLock()
		if m.closed {
			m.mu.RUnlock()
			return
		}
		persistentDir := m.cfg.PersistentDir
		writeLyrics := m.cfg.WriteLyricsToDir
		m.mu.RUnlock()
		if writeLyrics == nil {
			return
		}

		if _, err := writeLyrics(persistentDir, track.Artist, track.Title, lyrics); err != nil {
			m.debugf("deferred save failed artist=%q track=%q err=%v", track.Artist, track.Title, err)
			return
		}

		if mode == m.cfg.ModeRAMWithPrompt {
			m.notifyPromptStatus(prompt, track, lyrics, "Lyrics saved")
		}
	}()
}

func (m *LyricsManager) notifyPromptStatus(prompt PromptCallback, track Track, lyrics, message string) {
	if prompt == nil {
		return
	}
	prompt(context.Background(), PromptRequest{
		Track:   track,
		Lyrics:  lyrics,
		Message: message,
	})
}

func (m *LyricsManager) hasPersistentTimedLyrics(track Track) bool {
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return false
	}
	persistentDir := m.cfg.PersistentDir
	readLyrics := m.cfg.ReadLyricsFromDir
	hasTimed := m.cfg.HasTimedLyrics
	m.mu.RUnlock()

	if readLyrics == nil {
		return false
	}

	lyrics, err := readLyrics(persistentDir, track.Artist, track.Title)
	if err != nil {
		return false
	}

	if hasTimed == nil {
		return true
	}

	return hasTimed(lyrics)
}

func (m *LyricsManager) isActiveTrack(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return !m.closed && m.activeTrackKey == key
}
