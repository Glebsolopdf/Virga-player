package app

import (
	"fmt"
	"strings"
	"time"

	"virga-player/music"
)

func trackIdentityKey(track *music.TrackInfo) string {
	if track == nil {
		return ""
	}
	return fmt.Sprintf("%s\x00%s\x00%s\x00%s", track.Source, strings.TrimSpace(track.Artist), strings.TrimSpace(track.Title), strings.TrimSpace(track.Album))
}

func (a *App) effectiveTrackElapsed(track *music.TrackInfo, now time.Time) int {
	if track == nil || track.Source == "default" {
		a.timelineTrackKey = ""
		a.timelineUsingFallback = false
		a.timelineLastElapsed = 0
		return 0
	}

	key := trackIdentityKey(track)
	incoming := track.Elapsed
	if incoming < 0 {
		incoming = 0
	}

	if key != a.timelineTrackKey {
		a.timelineTrackKey = key
		a.timelineLastElapsed = incoming
		a.timelineFallbackFrom = now.Add(-time.Duration(incoming) * time.Second)
		a.timelineUsingFallback = incoming <= 0
		if a.timelineUsingFallback && a.debug != nil {
			a.debug.Debugf("timeline fallback enabled source=%s artist=%q track=%q", track.Source, track.Artist, track.Title)
		}
	}

	if incoming > 0 {
		a.timelineLastElapsed = incoming
		a.timelineFallbackFrom = now.Add(-time.Duration(incoming) * time.Second)
		a.timelineUsingFallback = false
		return clampElapsedByDuration(incoming, track.Duration)
	}

	if track.Paused {
		a.timelineFallbackFrom = now.Add(-time.Duration(a.timelineLastElapsed) * time.Second)
		return clampElapsedByDuration(a.timelineLastElapsed, track.Duration)
	}

	if !a.timelineUsingFallback {
		a.timelineUsingFallback = true
		a.timelineFallbackFrom = now.Add(-time.Duration(a.timelineLastElapsed) * time.Second)
		if a.debug != nil {
			a.debug.Debugf("timeline fallback switched on source=%s artist=%q track=%q", track.Source, track.Artist, track.Title)
		}
	}

	effective := int(now.Sub(a.timelineFallbackFrom).Seconds())
	if effective < 0 {
		effective = 0
	}
	a.timelineLastElapsed = effective
	return clampElapsedByDuration(effective, track.Duration)
}

func clampElapsedByDuration(elapsed, duration int) int {
	if elapsed < 0 {
		return 0
	}
	if duration > 0 && elapsed > duration {
		return duration
	}
	return elapsed
}
