package music

import (
	"sync"
	"time"
)

type TrackInfo struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int    `json:"duration"`
	Elapsed     int    `json:"elapsed"`
	ArtworkURL  string `json:"artwork_url"`
	ArtworkPath string `json:"artwork_path,omitempty"`
	Source      string
}

var (
	trackCacheMu sync.RWMutex
	cachedTrack  *TrackInfo
)

const (
	fallbackDurationSeconds  = 4 * 60
	trackInfoRefreshInterval = 1 * time.Second
)

func init() {
	go trackInfoRefresher()
}

func trackInfoRefresher() {
	refreshTrackInfo()
	ticker := time.NewTicker(trackInfoRefreshInterval)
	defer ticker.Stop()
	for range ticker.C {
		refreshTrackInfo()
	}
}

func refreshTrackInfo() {
	track := fetchTrackInfo()
	if track == nil {
		track = getDefaultTrack()
	}
	trackCacheMu.Lock()
	cachedTrack = track
	trackCacheMu.Unlock()
}

func fetchTrackInfo() *TrackInfo {
	track := getPlayerctlTrack()
	if track != nil && track.Title != "" {
		normalizeTrackDuration(track)
		track.ArtworkPath = track.GetArtworkPath()
		return track
	}

	track = getJSONTrack()
	if track != nil && track.Title != "" {
		normalizeTrackDuration(track)
		track.ArtworkPath = track.GetArtworkPath()
		return track
	}

	track = getDefaultTrack()
	normalizeTrackDuration(track)
	track.ArtworkPath = track.GetArtworkPath()
	return track
}

func GetTrackInfo() *TrackInfo {
	trackCacheMu.RLock()
	track := cachedTrack
	trackCacheMu.RUnlock()
	if track == nil {
		return getDefaultTrack()
	}
	return track.Clone()
}

func (t *TrackInfo) Clone() *TrackInfo {
	if t == nil {
		return nil
	}
	return &TrackInfo{
		Title:       t.Title,
		Artist:      t.Artist,
		Album:       t.Album,
		Duration:    t.Duration,
		Elapsed:     t.Elapsed,
		ArtworkURL:  t.ArtworkURL,
		ArtworkPath: t.ArtworkPath,
		Source:      t.Source,
	}
}

func normalizeTrackDuration(track *TrackInfo) {
	if track == nil {
		return
	}

	if track.Elapsed < 0 {
		track.Elapsed = 0
	}

	if track.Duration <= 0 {
		if track.Title == "No Track Playing" {
			return
		}
		track.Duration = fallbackDurationSeconds
		if track.Elapsed >= track.Duration {
			track.Elapsed = track.Duration - 1
		}
		return
	}

	if track.Elapsed > track.Duration {
		track.Elapsed = track.Duration
	}
}
