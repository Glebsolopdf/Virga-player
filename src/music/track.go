package music

import "time"

type TrackInfo struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Album      string `json:"album"`
	Duration   int    `json:"duration"`
	Elapsed    int    `json:"elapsed"`
	ArtworkURL string `json:"artwork_url"`
}

var lastTrack *TrackInfo
var lastUpdate time.Time

const fallbackDurationSeconds = 4 * 60

func GetTrackInfo() *TrackInfo {
	now := time.Now()
	if lastTrack != nil && now.Sub(lastUpdate) < 100*time.Millisecond {
		return lastTrack
	}

	track := getPlayerctlTrack()
	if track != nil && track.Title != "" {
		normalizeTrackDuration(track)
		lastTrack = track
		lastUpdate = now
		return track
	}

	track = getJSONTrack()
	if track != nil && track.Title != "" {
		normalizeTrackDuration(track)
		lastTrack = track
		lastUpdate = now
		return track
	}

	track = getDefaultTrack()
	normalizeTrackDuration(track)
	lastTrack = track
	lastUpdate = now
	return track
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
