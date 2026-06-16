package sixeldata

import (
	"sync"
	"time"
)

const cacheMax = 16

type cacheEntry struct {
	data     []byte
	storedAt time.Time
}

var (
	cacheMu sync.RWMutex
	cache   = map[string]cacheEntry{}
)

func Get(imagePath string) ([]byte, bool) {
	if imagePath == "" {
		return nil, false
	}
	cacheMu.RLock()
	entry, ok := cache[imagePath]
	cacheMu.RUnlock()
	if !ok || len(entry.data) == 0 {
		return nil, false
	}
	copyData := make([]byte, len(entry.data))
	copy(copyData, entry.data)
	return copyData, true
}

func Store(imagePath string, data []byte) {
	if imagePath == "" || len(data) == 0 {
		return
	}
	copyData := make([]byte, len(data))
	copy(copyData, data)

	cacheMu.Lock()
	cache[imagePath] = cacheEntry{data: copyData, storedAt: time.Now()}
	pruneLocked()
	cacheMu.Unlock()
}

func pruneLocked() {
	if len(cache) <= cacheMax {
		return
	}
	for len(cache) > cacheMax {
		oldestKey := ""
		var oldestTime time.Time
		for key, entry := range cache {
			if oldestKey == "" || entry.storedAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = entry.storedAt
			}
		}
		if oldestKey == "" {
			return
		}
		delete(cache, oldestKey)
	}
}
