package music

import (
	"strconv"
	"strings"
)

type trackMetadata struct {
	Title      string
	Artist     string
	Album      string
	ArtworkURL string
	TrackURL   string
	Duration   int
}

func readTrackMeta() (trackMetadata, bool) {
	if !playerctlAvailable() {
		return trackMetadata{}, false
	}

	output, err := playerctlRun("metadata", "--format", "{{xesam:title}}\t{{xesam:artist}}\t{{xesam:album}}\t{{mpris:artUrl}}\t{{mpris:length}}\t{{xesam:url}}")
	if err != nil {
		return trackMetadata{}, false
	}

	parts := strings.SplitN(output, "\t", 6)
	if len(parts) == 0 {
		return trackMetadata{}, false
	}

	title := playerctlClean(parts[0])
	if title == "" {
		return trackMetadata{}, false
	}

	meta := trackMetadata{Title: title}
	if len(parts) > 1 {
		meta.Artist = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		meta.Album = strings.TrimSpace(parts[2])
	}
	if len(parts) > 3 {
		meta.ArtworkURL = strings.TrimSpace(parts[3])
	}
	if len(parts) > 4 {
		meta.Duration = durationFromMicros(strings.TrimSpace(parts[4]))
	}
	if len(parts) > 5 {
		meta.TrackURL = strings.TrimSpace(parts[5])
	}
	return meta, true
}

func firstMetaValue(keys ...string) string {
	for _, key := range keys {
		val := metaValue(key)
		if val != "" {
			return val
		}
	}
	return ""
}

func durationFromMicros(val string) int {
	if val == "" || val == "0" {
		return 0
	}
	d, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return int(d / 1000000)
}

func playerctlPosition() int {
	val, err := playerctlRun("position")
	if err != nil || val == "" {
		return 0
	}
	pos, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return int(pos)
}

func playerctlStatus() string {
	val, err := playerctlRun("status")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(val)
}

func playerctlIsPaused() bool {
	return strings.EqualFold(playerctlStatus(), "Paused")
}

func metaValue(key string) string {
	if !playerctlAvailable() {
		return ""
	}
	out, err := playerctlRun("metadata", key)
	if err != nil {
		return ""
	}
	return playerctlClean(out)
}

func metaDump() string {
	if !playerctlAvailable() {
		return ""
	}
	out, err := playerctlRun("metadata")
	if err != nil {
		return ""
	}
	return out
}
