package music

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getJSONTrack() *TrackInfo {
	trackPath := filepath.Join("/tmp", "virga-player", "track.json")
	data, err := os.ReadFile(trackPath)
	if err != nil {
		return nil
	}

	var track TrackInfo
	if err := json.Unmarshal(data, &track); err != nil {
		return nil
	}

	if track.Title != "" && track.Title != "No Track" {
		return &track
	}

	return nil
}

func getDefaultTrack() *TrackInfo {
	return &TrackInfo{
		Title:      "No Track Playing",
		Artist:     "Start your music player",
		Album:      "or update /tmp/virga-player/track.json",
		Duration:   0,
		Elapsed:    0,
		ArtworkURL: "",
	}
}
