package music

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type artworkPathCacheEntry struct {
	path     string
	failedAt time.Time
}

var artworkPathCache = map[string]artworkPathCacheEntry{}

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
		cacheDir := filepath.Join(os.TempDir(), "virga-player", "artwork")
		os.MkdirAll(cacheDir, 0755)

		cacheFile := filepath.Join(cacheDir, fmt.Sprintf("%x.jpg", hashString(t.ArtworkURL)))

		if entry, ok := artworkPathCache[t.ArtworkURL]; ok {
			if entry.path != "" {
				if _, err := os.Stat(entry.path); err == nil && isValidImage(entry.path) {
					return entry.path
				}
			}
			if !entry.failedAt.IsZero() && time.Since(entry.failedAt) < 10*time.Minute {
				return ""
			}
		}

		if _, err := os.Stat(cacheFile); err == nil {
			if isValidImage(cacheFile) {
				artworkPathCache[t.ArtworkURL] = artworkPathCacheEntry{path: cacheFile}
				return cacheFile
			}
			_ = os.Remove(cacheFile)
		}

		if err := downloadFile(t.ArtworkURL, cacheFile); err == nil {
			if isValidImage(cacheFile) {
				artworkPathCache[t.ArtworkURL] = artworkPathCacheEntry{path: cacheFile}
				return cacheFile
			}
			_ = os.Remove(cacheFile)
		}

		artworkPathCache[t.ArtworkURL] = artworkPathCacheEntry{failedAt: time.Now()}
	}

	return ""
}

func downloadFile(url, filepath string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func isValidImage(filepath string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	_, format, err := image.DecodeConfig(f)
	return err == nil && (format == "jpeg" || format == "png")
}

func hashString(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = ((h << 5) - h) + uint32(c)
	}
	return h
}
