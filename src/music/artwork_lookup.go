package music

import (
	"strings"

	artworklookup "virga-player/music/artworklookup"
	playerctlcmd "virga-player/music/playerctlcmd"
)

func getArtworkURL() string {
	keys := []string{"mpris:artUrl", "xesam:artUrl", "artUrl", "xesam:artwork", "thumbnail"}
	for _, key := range keys {
		if value := playerctlcmd.MetadataValue(key); value != "" {
			return value
		}
	}

	if dump := playerctlcmd.MetadataDump(); dump != "" {
		for _, line := range strings.Split(dump, "\n") {
			lower := strings.ToLower(line)
			if !strings.Contains(lower, "art") && !strings.Contains(lower, "thumb") && !strings.Contains(lower, "image") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}

			candidate := strings.Trim(fields[len(fields)-1], "\"'")
			if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") || strings.HasPrefix(candidate, "file://") {
				return candidate
			}
		}
	}

	trackURL := playerctlcmd.MetadataValue("xesam:url")
	if trackURL != "" {
		if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
			artworklookup.Remember(trackURL, "")
			return ""
		}

		return artworklookup.Resolve(trackURL)
	}

	return ""
}

func getArtworkURLFromTrackURL(trackURL string) string {
	if trackURL == "" {
		return ""
	}

	if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
		artworklookup.Remember(trackURL, "")
		return ""
	}

	return artworklookup.Resolve(trackURL)
}
