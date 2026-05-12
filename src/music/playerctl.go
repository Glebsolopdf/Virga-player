package music

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	playerctlCheckOnce sync.Once
	playerctlAvailable bool
)

func hasPlayerctl() bool {
	playerctlCheckOnce.Do(func() {
		_, err := exec.LookPath("playerctl")
		playerctlAvailable = err == nil
	})
	return playerctlAvailable
}

func getPlayerctlTrack() *TrackInfo {
	if !hasPlayerctl() {
		return nil
	}

	metadataCmd := exec.Command("playerctl", "metadata", "--format", "{{xesam:title}}\t{{xesam:artist}}\t{{xesam:album}}\t{{mpris:artUrl}}\t{{mpris:length}}\t{{xesam:url}}")
	metadataOutput, err := metadataCmd.Output()
	if err != nil {
		return getPlayerctlTrackFallback()
	}

	parts := strings.SplitN(strings.TrimSpace(string(metadataOutput)), "\t", 6)
	if len(parts) == 0 {
		return nil
	}

	title := strings.TrimSpace(parts[0])
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
		if durationStr := strings.TrimSpace(parts[4]); durationStr != "" && durationStr != "0" {
			if d, err := strconv.ParseInt(durationStr, 10, 64); err == nil {
				duration = int(d / 1000000)
			}
		}
	}

	trackURL := ""
	if len(parts) > 5 {
		trackURL = strings.TrimSpace(parts[5])
	}

	if artworkURL == "" {
		artworkURL = getArtworkURLFromTrackURL(trackURL)
	}

	elapsed := 0
	positionCmd := exec.Command("playerctl", "position")
	positionOutput, _ := positionCmd.Output()
	if posStr := strings.TrimSpace(string(positionOutput)); posStr != "" {
		if p, err := strconv.ParseFloat(posStr, 64); err == nil {
			elapsed = int(p)
		}
	}

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
	titleCmd := exec.Command("playerctl", "metadata", "xesam:title")
	titleOutput, err := titleCmd.Output()
	if err != nil {
		titleCmd = exec.Command("playerctl", "metadata", "title")
		titleOutput, err = titleCmd.Output()
		if err != nil {
			return nil
		}
	}

	title := strings.TrimSpace(string(titleOutput))
	if title == "" || title == "No player could handle this command" {
		return nil
	}

	artistCmd := exec.Command("playerctl", "metadata", "xesam:artist")
	artistOutput, _ := artistCmd.Output()
	artist := strings.TrimSpace(string(artistOutput))
	if artist == "" {
		artistCmd = exec.Command("playerctl", "metadata", "artist")
		artistOutput, _ = artistCmd.Output()
		artist = strings.TrimSpace(string(artistOutput))
	}

	albumCmd := exec.Command("playerctl", "metadata", "xesam:album")
	albumOutput, _ := albumCmd.Output()
	album := strings.TrimSpace(string(albumOutput))
	if album == "" {
		albumCmd = exec.Command("playerctl", "metadata", "album")
		albumOutput, _ = albumCmd.Output()
		album = strings.TrimSpace(string(albumOutput))
	}

	artworkURL := getArtworkURL()

	duration := 0
	durationCmd := exec.Command("playerctl", "metadata", "mpris:length")
	durationOutput, _ := durationCmd.Output()
	if durationStr := strings.TrimSpace(string(durationOutput)); durationStr != "" && durationStr != "0" {
		if d, err := strconv.ParseInt(durationStr, 10, 64); err == nil {
			duration = int(d / 1000000)
		}
	}

	if duration == 0 {
		alternativeDurationCmd := exec.Command("playerctl", "metadata", "xesam:length")
		alternativeDurationOutput, _ := alternativeDurationCmd.Output()
		if altDurationStr := strings.TrimSpace(string(alternativeDurationOutput)); altDurationStr != "" {
			if d, err := strconv.ParseInt(altDurationStr, 10, 64); err == nil {
				duration = int(d / 1000000)
			}
		}
	}

	elapsed := 0
	positionCmd := exec.Command("playerctl", "position")
	positionOutput, _ := positionCmd.Output()
	if posStr := strings.TrimSpace(string(positionOutput)); posStr != "" {
		if p, err := strconv.ParseFloat(posStr, 64); err == nil {
			elapsed = int(p)
		}
	}

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

	cmd := exec.Command("playerctl", "metadata", key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	value := strings.TrimSpace(string(out))
	if value == "" || strings.Contains(strings.ToLower(value), "no player could handle this command") {
		return ""
	}

	return value
}

func getPlayerctlMetadataDump() string {
	if !hasPlayerctl() {
		return ""
	}

	cmd := exec.Command("playerctl", "metadata")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}
