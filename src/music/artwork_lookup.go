package music

import "strings"

func getArtworkURL() string {
	keys := []string{"mpris:artUrl", "xesam:artUrl", "artUrl", "xesam:artwork", "thumbnail"}
	for _, key := range keys {
		if value := metaValue(key); value != "" {
			return value
		}
	}

	if dump := metaDump(); dump != "" {
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

	trackURL := metaValue("xesam:url")
	if trackURL != "" {
		if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
			artworkRemember(trackURL, "")
			return ""
		}

		return artworkResolve(trackURL)
	}

	return ""
}

func getArtworkURLFromTrackURL(trackURL string) string {
	if trackURL == "" {
		return ""
	}

	if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
		artworkRemember(trackURL, "")
		return ""
	}

	return artworkResolve(trackURL)
}
