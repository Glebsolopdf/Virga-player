package lyricsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	searchEndpoint    = "https://lrclib.net/api/search"
	userAgent         = "VirgaPlayer/(https://github.com/Glebsolopdf/Virga-player)"
	maxSearchAttempts = 3
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

type retryableStatusError struct {
	status string
}

func (e retryableStatusError) Error() string {
	return fmt.Sprintf("lrclib search returned %s", e.status)
}

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

	log.Printf("lyricsearch: requesting lrclib artist=%q track=%q", artist, track)

	var results []searchResult
	var lastErr error
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		results, lastErr = searchOnce(query, artist, track)
		if lastErr == nil {
			break
		}
		if !isRetryableSearchError(lastErr) || attempt == maxSearchAttempts {
			return "", lastErr
		}
		backoff := time.Duration(attempt*attempt) * 350 * time.Millisecond
		log.Printf("lyricsearch: transient error, retrying attempt=%d/%d after=%s artist=%q track=%q err=%v", attempt+1, maxSearchAttempts, backoff, artist, track, lastErr)
		time.Sleep(backoff)
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

func searchOnce(query url.Values, artist, track string) ([]searchResult, error) {
	req, err := http.NewRequest(http.MethodGet, searchEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build lrclib request: %w", err)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request lrclib: %w", err)
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

		switch {
		case resp.StatusCode == http.StatusTooManyRequests:
			return nil, retryableStatusError{status: resp.Status}
		case resp.StatusCode >= 500 && resp.StatusCode <= 599:
			return nil, retryableStatusError{status: resp.Status}
		default:
			return nil, fmt.Errorf("lrclib search returned %s", resp.Status)
		}
	}

	var results []searchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Printf("lyricsearch: invalid lrclib json artist=%q track=%q err=%v", artist, track, err)
		return nil, fmt.Errorf("decode lrclib response: %w", err)
	}

	return results, nil
}

func isRetryableSearchError(err error) bool {
	var retryStatus retryableStatusError
	if errors.As(err, &retryStatus) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "timeout") || strings.Contains(errText, "deadline exceeded") || strings.Contains(errText, "connection reset")
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
