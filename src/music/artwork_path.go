package music

import (
	"os"
	"strings"
)

func (t *TrackInfo) GetArtworkPath() string {
	if t.ArtworkURL == "" {
		return ""
	}

	if strings.HasPrefix(t.ArtworkURL, "file://") {
		path := strings.TrimPrefix(t.ArtworkURL, "file://")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		return ""
	}

	if strings.HasPrefix(t.ArtworkURL, "http://") || strings.HasPrefix(t.ArtworkURL, "https://") {
		return t.ArtworkURL
	}

	if _, err := os.Stat(t.ArtworkURL); err == nil {
		return t.ArtworkURL
	}

	return ""
}
