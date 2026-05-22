package music

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	playerctlCheckOnce sync.Once
	playerctlAvailable bool
)

const playerctlCommandTimeout = 750 * time.Millisecond

func hasPlayerctl() bool {
	playerctlCheckOnce.Do(func() {
		_, err := exec.LookPath("playerctl")
		playerctlAvailable = err == nil
	})
	return playerctlAvailable
}

func runPlayerctl(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), playerctlCommandTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "playerctl", args...).Output()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func cleanPlayerctlValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(strings.ToLower(value), "no player could handle this command") {
		return ""
	}
	return value
}

func getFirstPlayerctlMetadataValue(keys ...string) string {
	for _, key := range keys {
		value, err := runPlayerctl("metadata", key)
		if err != nil {
			continue
		}
		if value = cleanPlayerctlValue(value); value != "" {
			return value
		}
	}
	return ""
}

func parsePlayerctlDuration(value string) int {
	if value == "" || value == "0" {
		return 0
	}
	duration, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int(duration / 1000000)
}

func getPlayerctlPositionSeconds() int {
	value, err := runPlayerctl("position")
	if err != nil || value == "" {
		return 0
	}
	position, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return int(position)
}

func getPlayerctlTrack() *TrackInfo {
	if !hasPlayerctl() {
		return nil
	}

	metadataOutput, err := runPlayerctl("metadata", "--format", "{{xesam:title}}\t{{xesam:artist}}\t{{xesam:album}}\t{{mpris:artUrl}}\t{{mpris:length}}\t{{xesam:url}}")
	if err != nil {
		return getPlayerctlTrackFallback()
	}

	parts := strings.SplitN(metadataOutput, "\t", 6)
	if len(parts) == 0 {
		return nil
	}

	title := cleanPlayerctlValue(parts[0])
	if title == "" || title == "No player could handle this command" {
		return nil
	}

	artist := ""
	if len(parts) > 1 {
		artist = strings.TrimSpace(parts[1])
	}

	album := ""
	if len(parts) > 2 {
		album = strings.TrimSpace(parts[2])
	}

	artworkURL := ""
	if len(parts) > 3 {
		artworkURL = strings.TrimSpace(parts[3])
	}

	duration := 0
	if len(parts) > 4 {
		duration = parsePlayerctlDuration(strings.TrimSpace(parts[4]))
	}

	trackURL := ""
	if len(parts) > 5 {
		trackURL = strings.TrimSpace(parts[5])
	}

	if artworkURL == "" {
		artworkURL = getArtworkURLFromTrackURL(trackURL)
	}

	elapsed := getPlayerctlPositionSeconds()

	return &TrackInfo{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Duration:   duration,
		Elapsed:    elapsed,
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}

func getPlayerctlTrackFallback() *TrackInfo {
	title := getFirstPlayerctlMetadataValue("xesam:title", "title")
	if title == "" || title == "No player could handle this command" {
		return nil
	}

	artist := getFirstPlayerctlMetadataValue("xesam:artist", "artist")
	album := getFirstPlayerctlMetadataValue("xesam:album", "album")

	artworkURL := getArtworkURL()

	duration := parsePlayerctlDuration(getFirstPlayerctlMetadataValue("mpris:length"))
	if duration == 0 {
		duration = parsePlayerctlDuration(getFirstPlayerctlMetadataValue("xesam:length"))
	}

	elapsed := getPlayerctlPositionSeconds()

	return &TrackInfo{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Duration:   duration,
		Elapsed:    elapsed,
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}

func getPlayerctlMetadataValue(key string) string {
	if !hasPlayerctl() {
		return ""
	}

	out, err := runPlayerctl("metadata", key)
	if err != nil {
		return ""
	}
	return cleanPlayerctlValue(out)
}

func getPlayerctlMetadataDump() string {
	if !hasPlayerctl() {
		return ""
	}

	out, err := runPlayerctl("metadata")
	if err != nil {
		return ""
	}
	return out
}
