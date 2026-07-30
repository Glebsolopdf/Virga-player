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

var ogImageRe = regexp.MustCompile(`(?i)property=["']og:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*property=["']og:image["']`)
var twitterImageRe = regexp.MustCompile(`(?i)name=["']twitter:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*name=["']twitter:image["']`)

const (
	artCachePosTTL = 10 * time.Minute
	artCacheNegTTL = 30 * time.Second
	artCacheMax    = 128
)

type artCacheEntry struct {
	url       string
	expiresAt time.Time
	storedAt  time.Time
}

type artInflight struct {
	done   chan struct{}
	result string
}

var (
	artCacheMu    sync.RWMutex
	artCache      = map[string]artCacheEntry{}
	artFlightMu   sync.Mutex
	artInFlight   = map[string]*artInflight{}
)

func artworkResolve(trackURL string) string {
	if cached, ok := artCacheGet(trackURL); ok {
		return cached
	}

	artFlightMu.Lock()
	if call, ok := artInFlight[trackURL]; ok {
		artFlightMu.Unlock()
		<-call.done
		return call.result
	}
	call := &artInflight{done: make(chan struct{})}
	artInFlight[trackURL] = call
	artFlightMu.Unlock()

	result := fetchArtwork(trackURL)
	artworkRemember(trackURL, result)

	artFlightMu.Lock()
	call.result = result
	delete(artInFlight, trackURL)
	close(call.done)
	artFlightMu.Unlock()

	return result
}

func artworkRemember(trackURL, artworkURL string) {
	ttl := artCachePosTTL
	if artworkURL == "" {
		ttl = artCacheNegTTL
	}
	now := time.Now()

	artCacheMu.Lock()
	artCache[trackURL] = artCacheEntry{
		url:       artworkURL,
		expiresAt: now.Add(ttl),
		storedAt:  now,
	}
	artCachePrune(now)
	artCacheMu.Unlock()
}

func artCacheGet(trackURL string) (string, bool) {
	artCacheMu.RLock()
	entry, ok := artCache[trackURL]
	artCacheMu.RUnlock()
	if !ok {
		return "", false
	}
	if time.Now().After(entry.expiresAt) {
		artCacheMu.Lock()
		delete(artCache, trackURL)
		artCacheMu.Unlock()
		return "", false
	}
	return entry.url, true
}

func artCachePrune(now time.Time) {
	if len(artCache) <= artCacheMax {
		return
	}
	for k, e := range artCache {
		if now.After(e.expiresAt) {
			delete(artCache, k)
		}
	}
	if len(artCache) <= artCacheMax {
		return
	}
	for len(artCache) > artCacheMax {
		oldestKey := ""
		var oldestTime time.Time
		for k, e := range artCache {
			if oldestKey == "" || e.storedAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = e.storedAt
			}
		}
		if oldestKey == "" {
			return
		}
		delete(artCache, oldestKey)
	}
}

func fetchArtwork(trackURL string) string {
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
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return ""
	}
	return extractMetaImage(string(body))
}

func extractMetaImage(htmlBody string) string {
	for _, re := range []*regexp.Regexp{ogImageRe, twitterImageRe} {
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
