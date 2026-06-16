package artworklookup

import (
	"sync"
	"time"
)

const (
	positiveTTL = 10 * time.Minute
	negativeTTL = 30 * time.Second
	cacheMax    = 128
)

type cacheEntry struct {
	artworkURL string
	expiresAt  time.Time
	storedAt   time.Time
}

type inflightCall struct {
	done   chan struct{}
	result string
}

var (
	cacheMu    sync.RWMutex
	cache      = map[string]cacheEntry{}
	inFlightMu sync.Mutex
	inFlight   = map[string]*inflightCall{}
)

func Resolve(trackURL string) string {
	if cached, ok := getCached(trackURL); ok {
		return cached
	}

	inFlightMu.Lock()
	if call, ok := inFlight[trackURL]; ok {
		inFlightMu.Unlock()
		<-call.done
		return call.result
	}
	call := &inflightCall{done: make(chan struct{})}
	inFlight[trackURL] = call
	inFlightMu.Unlock()

	result := fetchTrackPageArtwork(trackURL)
	Remember(trackURL, result)

	inFlightMu.Lock()
	call.result = result
	delete(inFlight, trackURL)
	close(call.done)
	inFlightMu.Unlock()

	return result
}

func Remember(trackURL, artworkURL string) {
	ttl := positiveTTL
	if artworkURL == "" {
		ttl = negativeTTL
	}
	now := time.Now()

	cacheMu.Lock()
	cache[trackURL] = cacheEntry{
		artworkURL: artworkURL,
		expiresAt:  now.Add(ttl),
		storedAt:   now,
	}
	pruneCacheLocked(now)
	cacheMu.Unlock()
}

func getCached(trackURL string) (string, bool) {
	cacheMu.RLock()
	entry, ok := cache[trackURL]
	cacheMu.RUnlock()
	if !ok {
		return "", false
	}
	if time.Now().After(entry.expiresAt) {
		cacheMu.Lock()
		delete(cache, trackURL)
		cacheMu.Unlock()
		return "", false
	}
	return entry.artworkURL, true
}

func pruneCacheLocked(now time.Time) {
	if len(cache) <= cacheMax {
		return
	}

	for trackURL, entry := range cache {
		if now.After(entry.expiresAt) {
			delete(cache, trackURL)
		}
	}
	if len(cache) <= cacheMax {
		return
	}

	for len(cache) > cacheMax {
		oldestKey := ""
		var oldestTime time.Time
		for trackURL, entry := range cache {
			if oldestKey == "" || entry.storedAt.Before(oldestTime) {
				oldestKey = trackURL
				oldestTime = entry.storedAt
			}
		}
		if oldestKey == "" {
			return
		}
		delete(cache, oldestKey)
	}
}
