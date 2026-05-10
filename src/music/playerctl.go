package music

import (
	"os/exec"
	"strconv"
	"strings"
)

func getPlayerctlTrack() *TrackInfo {
	_, err := exec.LookPath("playerctl")
	if err != nil {
		return nil
	}

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
	}
}

func getPlayerctlMetadataValue(key string) string {
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
	cmd := exec.Command("playerctl", "metadata")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}
