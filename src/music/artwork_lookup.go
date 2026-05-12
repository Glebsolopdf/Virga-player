package music

import (
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var ogImageRegexp = regexp.MustCompile(`(?i)property=["']og:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*property=["']og:image["']`)
var twitterImageRegexp = regexp.MustCompile(`(?i)name=["']twitter:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*name=["']twitter:image["']`)

const (
	artworkLookupPositiveTTL = 10 * time.Minute
	artworkLookupNegativeTTL = 30 * time.Second
)

type artworkLookupCacheEntry struct {
	artworkURL string
	expiresAt  time.Time
}

var (
	artworkLookupCacheMu sync.RWMutex
	artworkLookupCache   = map[string]artworkLookupCacheEntry{}
)

func getCachedArtworkURL(trackURL string) (string, bool) {
	artworkLookupCacheMu.RLock()
	entry, ok := artworkLookupCache[trackURL]
	artworkLookupCacheMu.RUnlock()
	if !ok {
		return "", false
	}
	if time.Now().After(entry.expiresAt) {
		artworkLookupCacheMu.Lock()
		delete(artworkLookupCache, trackURL)
		artworkLookupCacheMu.Unlock()
		return "", false
	}
	return entry.artworkURL, true
}

func storeCachedArtworkURL(trackURL, artworkURL string) {
	ttl := artworkLookupPositiveTTL
	if artworkURL == "" {
		ttl = artworkLookupNegativeTTL
	}

	artworkLookupCacheMu.Lock()
	artworkLookupCache[trackURL] = artworkLookupCacheEntry{
		artworkURL: artworkURL,
		expiresAt:  time.Now().Add(ttl),
	}
	artworkLookupCacheMu.Unlock()
}

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
		if cached, ok := getCachedArtworkURL(trackURL); ok {
			return cached
		}

		if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
			storeCachedArtworkURL(trackURL, "")
			return ""
		}

		artworkURL := getArtworkFromTrackPage(trackURL)
		storeCachedArtworkURL(trackURL, artworkURL)
		return artworkURL
	}

	return ""
}

func getArtworkURLFromTrackURL(trackURL string) string {
	if trackURL == "" {
		return ""
	}

	if cached, ok := getCachedArtworkURL(trackURL); ok {
		return cached
	}

	if !strings.HasPrefix(trackURL, "http://") && !strings.HasPrefix(trackURL, "https://") {
		storeCachedArtworkURL(trackURL, "")
		return ""
	}

	artworkURL := getArtworkFromTrackPage(trackURL)
	storeCachedArtworkURL(trackURL, artworkURL)
	return artworkURL
}

func getArtworkFromTrackPage(trackURL string) string {
	client := &http.Client{Timeout: 4 * time.Second}
	req, err := http.NewRequest(http.MethodGet, trackURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return ""
	}

	htmlBody := string(body)
	return extractMetaImageURL(htmlBody)
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
