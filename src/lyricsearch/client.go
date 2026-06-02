package lyricsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	searchEndpoint = "https://lrclib.net/api/search"
	userAgent      = "VirgaPlayer/(https://github.com/Glebsolopdf/Virga-player)"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type searchResult struct {
	ID           int     `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

func fetchSyncedLyrics(artist, track string) (string, error) {
	query := url.Values{}
	query.Set("track_name", track)
	query.Set("artist_name", artist)

	req, err := http.NewRequest(http.MethodGet, searchEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("build lrclib request: %w", err)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", userAgent)

	log.Printf("lyricsearch: requesting lrclib artist=%q track=%q", artist, track)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request lrclib: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if readErr != nil {
			log.Printf("lyricsearch: failed to read lrclib error body artist=%q track=%q err=%v", artist, track, readErr)
		}
		trimmedBody := strings.TrimSpace(string(body))
		if trimmedBody == "" {
			log.Printf("lyricsearch: lrclib returned status=%s artist=%q track=%q", resp.Status, artist, track)
		} else {
			log.Printf("lyricsearch: lrclib returned status=%s artist=%q track=%q body=%q", resp.Status, artist, track, trimmedBody)
		}
		return "", fmt.Errorf("lrclib search returned %s", resp.Status)
	}

	var results []searchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Printf("lyricsearch: invalid lrclib json artist=%q track=%q err=%v", artist, track, err)
		return "", fmt.Errorf("decode lrclib response: %w", err)
	}

	log.Printf("lyricsearch: lrclib returned %d results artist=%q track=%q", len(results), artist, track)
	instrumental := false
	for _, result := range results {
		if result.Instrumental {
			instrumental = true
		}
		lyrics := strings.TrimSpace(result.SyncedLyrics)
		if lyrics == "" {
			continue
		}
		if !hasTimedLyrics(lyrics) {
			log.Printf("lyricsearch: skipping unsynced lyrics candidate id=%d artist=%q track=%q", result.ID, result.ArtistName, result.TrackName)
			continue
		}

		log.Printf("lyricsearch: selected synced lyrics id=%d artist=%q track=%q", result.ID, artist, track)
		return lyrics, nil
	}

	if instrumental {
		log.Printf("lyricsearch: track marked instrumental artist=%q track=%q", artist, track)
		return "", fmt.Errorf("%w: %s - %s", ErrInstrumentalTrack, artist, track)
	}

	log.Printf("lyricsearch: no synced lyrics found artist=%q track=%q", artist, track)
	return "", fmt.Errorf("%w: %s - %s", ErrLyricsNotFound, artist, track)
}

func hasTimedLyrics(lyrics string) bool {
	for _, line := range strings.Split(strings.ReplaceAll(lyrics, "\r\n", "\n"), "\n") {
		remaining := strings.TrimSpace(line)
		for strings.HasPrefix(remaining, "[") {
			end := strings.IndexByte(remaining, ']')
			if end <= 1 {
				break
			}
			if isTimestampTag(remaining[1:end]) {
				return true
			}
			remaining = strings.TrimLeft(remaining[end+1:], " \t")
		}
	}
	return false
}

func isTimestampTag(tag string) bool {
	parts := strings.SplitN(strings.TrimSpace(tag), ":", 2)
	if len(parts) != 2 {
		return false
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil || minutes < 0 {
		return false
	}

	secondsPart := strings.ReplaceAll(strings.TrimSpace(parts[1]), ",", ".")
	seconds, err := strconv.ParseFloat(secondsPart, 64)
	if err != nil || seconds < 0 {
		return false
	}

	return true
}
