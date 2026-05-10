package music

import (
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type artworkLookupCacheEntry struct {
	url       string
	updatedAt time.Time
}

var artworkLookupCache = map[string]artworkLookupCacheEntry{}

var ogImageRegexp = regexp.MustCompile(`(?i)property=["']og:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*property=["']og:image["']`)
var twitterImageRegexp = regexp.MustCompile(`(?i)name=["']twitter:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*name=["']twitter:image["']`)

func getArtworkURL() string {
	keys := []string{"mpris:artUrl", "xesam:artUrl", "artUrl", "xesam:artwork", "thumbnail"}
	for _, key := range keys {
		if value := getPlayerctlMetadataValue(key); value != "" {
			return value
		}
	}

	if dump := getPlayerctlMetadataDump(); dump != "" {
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

	trackURL := getPlayerctlMetadataValue("xesam:url")
	if trackURL != "" {
		return getArtworkFromTrackPage(trackURL)
	}

	return ""
}

func getArtworkFromTrackPage(trackURL string) string {
	if cached, ok := artworkLookupCache[trackURL]; ok {
		if time.Since(cached.updatedAt) < 10*time.Minute {
			return cached.url
		}
	}

	client := &http.Client{Timeout: 4 * time.Second}
	req, err := http.NewRequest(http.MethodGet, trackURL, nil)
	if err != nil {
		artworkLookupCache[trackURL] = artworkLookupCacheEntry{url: "", updatedAt: time.Now()}
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		artworkLookupCache[trackURL] = artworkLookupCacheEntry{url: "", updatedAt: time.Now()}
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		artworkLookupCache[trackURL] = artworkLookupCacheEntry{url: "", updatedAt: time.Now()}
		return ""
	}

	htmlBody := string(body)
	url := extractMetaImageURL(htmlBody)
	artworkLookupCache[trackURL] = artworkLookupCacheEntry{url: url, updatedAt: time.Now()}
	return url
}

func extractMetaImageURL(htmlBody string) string {
	for _, re := range []*regexp.Regexp{ogImageRegexp, twitterImageRegexp} {
		matches := re.FindStringSubmatch(htmlBody)
		if len(matches) < 2 {
			continue
		}

		for i := 1; i < len(matches); i++ {
			candidate := strings.TrimSpace(html.UnescapeString(matches[i]))
			if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") {
				return candidate
			}
		}
	}

	return ""
}
