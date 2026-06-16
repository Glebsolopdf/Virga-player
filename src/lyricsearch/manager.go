package lyricsearch

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ModeDisabled      = "disabled"
	ModeRAMOnly       = "ram-only"
	ModeRAMWithAuto   = "ram-with-auto-save"
	ModeRAMWithPrompt = "ram-with-save-prompt"
	ModeDirectToDisk  = "direct-to-disk"

	// Legacy modes kept for backward compatibility with older configs.
	ModeLocal = "local"
	ModeAuto  = "auto"
)

var (
	ErrLyricsDisabled    = errors.New("lyrics mode is disabled")
	ErrLyricsNotFound    = errors.New("synced lyrics not found")
	ErrInstrumentalTrack = errors.New("track is instrumental")
	ErrMissingMetadata   = errors.New("lyrics search requires artist and track")
)

// GetLyrics returns synchronized lyrics according to the requested mode.
func GetLyrics(artist, track, mode string, saveToCache bool) (string, error) {
	artist = strings.TrimSpace(artist)
	track = strings.TrimSpace(track)
	mode = normalizeMode(mode)

	if mode == ModeDisabled {
		log.Printf("lyricsearch: skipped mode=disabled artist=%q track=%q", artist, track)
		return "", ErrLyricsDisabled
	}

	if artist == "" || track == "" {
		log.Printf("lyricsearch: missing metadata artist=%q track=%q", artist, track)
		return "", ErrMissingMetadata
	}

	var (
		lyrics string
		err    error
	)

	useCache := shouldUseCache(mode, saveToCache)
	if useCache {
		lyrics, err = readCachedLyrics(artist, track)
		switch {
		case err == nil:
			if !hasTimedLyrics(lyrics) {
				log.Printf("lyricsearch: cache entry is not time-synced artist=%q track=%q", artist, track)
				if mode == ModeLocal {
					return "", fmt.Errorf("%w: %s - %s", ErrLyricsNotFound, artist, track)
				}
				break
			}
			log.Printf("lyricsearch: cache hit artist=%q track=%q", artist, track)
			return lyrics, nil
		case errors.Is(err, os.ErrNotExist):
			log.Printf("lyricsearch: cache miss artist=%q track=%q", artist, track)
		default:
			log.Printf("lyricsearch: cache read failed artist=%q track=%q err=%v", artist, track, err)
			if mode == ModeLocal {
				return "", fmt.Errorf("read cached lyrics: %w", err)
			}
		}
	} else if mode == ModeAuto {
		log.Printf("lyricsearch: cache bypass enabled artist=%q track=%q", artist, track)
	}

	if mode == ModeLocal {
		return "", fmt.Errorf("%w: %s - %s", ErrLyricsNotFound, artist, track)
	}

	lyrics, err = fetchSyncedLyrics(artist, track)
	if err != nil {
		log.Printf("lyricsearch: remote fetch failed artist=%q track=%q err=%v", artist, track, err)
		return "", err
	}

	if !useCache {
		log.Printf("lyricsearch: returning remote lyrics without caching artist=%q track=%q", artist, track)
		return lyrics, nil
	}

	if err := writeCachedLyrics(artist, track, lyrics); err != nil {
		log.Printf("lyricsearch: cache write failed artist=%q track=%q err=%v", artist, track, err)
		return lyrics, fmt.Errorf("write cached lyrics: %w", err)
	}

	log.Printf("lyricsearch: cached remote lyrics artist=%q track=%q", artist, track)
	return lyrics, nil
}

func normalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case ModeRAMOnly:
		return ModeRAMOnly
	case ModeRAMWithAuto:
		return ModeRAMWithAuto
	case ModeRAMWithPrompt:
		return ModeRAMWithPrompt
	case ModeDirectToDisk:
		return ModeDirectToDisk
	case ModeLocal:
		return ModeLocal
	case ModeAuto:
		return ModeAuto
	default:
		return ModeDisabled
	}
}

func shouldUseCache(mode string, saveToCache bool) bool {
	return mode == ModeLocal || (mode == ModeAuto && saveToCache)
}
