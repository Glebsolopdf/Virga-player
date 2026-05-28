package playerctlcmd

import (
	"strconv"
	"strings"
)

type TrackMetadata struct {
	Title      string
	Artist     string
	Album      string
	ArtworkURL string
	TrackURL   string
	Duration   int
}

func ReadTrackMetadata() (TrackMetadata, bool) {
	if !Available() {
		return TrackMetadata{}, false
	}

	output, err := run("metadata", "--format", "{{xesam:title}}\t{{xesam:artist}}\t{{xesam:album}}\t{{mpris:artUrl}}\t{{mpris:length}}\t{{xesam:url}}")
	if err != nil {
		return TrackMetadata{}, false
	}

	parts := strings.SplitN(output, "\t", 6)
	if len(parts) == 0 {
		return TrackMetadata{}, false
	}

	title := cleanValue(parts[0])
	if title == "" {
		return TrackMetadata{}, false
	}

	metadata := TrackMetadata{Title: title}
	if len(parts) > 1 {
		metadata.Artist = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		metadata.Album = strings.TrimSpace(parts[2])
	}
	if len(parts) > 3 {
		metadata.ArtworkURL = strings.TrimSpace(parts[3])
	}
	if len(parts) > 4 {
		metadata.Duration = DurationFromMicros(strings.TrimSpace(parts[4]))
	}
	if len(parts) > 5 {
		metadata.TrackURL = strings.TrimSpace(parts[5])
	}
	return metadata, true
}

func FirstMetadataValue(keys ...string) string {
	for _, key := range keys {
		value := MetadataValue(key)
		if value != "" {
			return value
		}
	}
	return ""
}

func DurationFromMicros(value string) int {
	if value == "" || value == "0" {
		return 0
	}
	duration, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int(duration / 1000000)
}

func PositionSeconds() int {
	value, err := run("position")
	if err != nil || value == "" {
		return 0
	}
	position, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return int(position)
}

func MetadataValue(key string) string {
	if !Available() {
		return ""
	}

	out, err := run("metadata", key)
	if err != nil {
		return ""
	}
	return cleanValue(out)
}

func MetadataDump() string {
	if !Available() {
		return ""
	}

	out, err := run("metadata")
	if err != nil {
		return ""
	}
	return out
}
