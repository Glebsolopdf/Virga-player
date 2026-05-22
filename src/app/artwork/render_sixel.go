package artwork

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"math"
	"os"
	"os/exec"
	"sync"
	"time"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

var (
	sixelDataCacheMu sync.RWMutex
	sixelDataCache   = map[string]sixelCacheEntry{}
)

const (
	sixelBuildTimeout = 2 * time.Second
	sixelCacheMax     = 16
)

type sixelCacheEntry struct {
	data     []byte
	storedAt time.Time
}

func (a *Artwork) renderSixel(screen tcell.Screen, state artworkSnapshot) {
	theme := settings.CurrentTheme()
	if state.coverImg == nil {
		a.renderTextOnly(screen, state)
		return
	}

	hasSixel := len(state.sixelData) > 0
	if !hasSixel {
		a.prepareSixelDataAsync()
		a.renderTextOnly(screen, state)
		return
	}

	screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.TrackTitle))

	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	centerX += int(math.Round(state.rainOffsetX))
	centerY += int(math.Round(state.rainOffsetY))

	imgY := centerY - 12
	imgX := centerX - 20
	if imgX < 0 {
		imgX = 0
	}
	if imgY < 0 {
		imgY = 0
	}

	fmt.Printf("\x1B[%d;%dH", imgY+1, imgX+1)
	_, _ = os.Stdout.Write(state.sixelData)

	textWidth := 40
	if textWidth > w-10 {
		textWidth = w - 10
	}
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+3, state.title, theme.TrackTitle)
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+5, state.artist, theme.TrackArtist)
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+7, state.album, theme.TrackAlbum)

	if state.duration > 0 {
		a.drawTimeline(screen, centerX, centerY+9, 28, state.elapsed, state.duration)
		timeStr := formatTime(state.elapsed) + " / " + formatTime(state.duration)
		a.drawText(screen, centerX-len(timeStr)/2, centerY+11, timeStr, theme.TrackTime)
	}
}

func getCachedSixelData(imagePath string) ([]byte, bool) {
	if imagePath == "" {
		return nil, false
	}
	sixelDataCacheMu.RLock()
	entry, ok := sixelDataCache[imagePath]
	sixelDataCacheMu.RUnlock()
	if !ok || len(entry.data) == 0 {
		return nil, false
	}
	copyData := make([]byte, len(entry.data))
	copy(copyData, entry.data)
	return copyData, true
}

func storeCachedSixelData(imagePath string, data []byte) {
	if imagePath == "" || len(data) == 0 {
		return
	}
	copyData := make([]byte, len(data))
	copy(copyData, data)
	sixelDataCacheMu.Lock()
	sixelDataCache[imagePath] = sixelCacheEntry{data: copyData, storedAt: time.Now()}
	pruneSixelCacheLocked()
	sixelDataCacheMu.Unlock()
}

func pruneSixelCacheLocked() {
	if len(sixelDataCache) <= sixelCacheMax {
		return
	}
	for len(sixelDataCache) > sixelCacheMax {
		oldestKey := ""
		var oldestTime time.Time
		for key, entry := range sixelDataCache {
			if oldestKey == "" || entry.storedAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = entry.storedAt
			}
		}
		if oldestKey == "" {
			return
		}
		delete(sixelDataCache, oldestKey)
	}
}

func (a *Artwork) prepareSixelDataAsync() {
	a.mu.Lock()
	if a.sixelBuilding || len(a.SixelData) > 0 {
		a.mu.Unlock()
		return
	}
	img := a.CoverImg
	imagePath := a.ImagePath
	a.sixelBuilding = true
	a.mu.Unlock()

	if img == nil {
		a.mu.Lock()
		a.sixelBuilding = false
		a.mu.Unlock()
		return
	}

	go func(localImg image.Image, path string) {
		output, ok := buildSixelData(localImg)
		a.mu.Lock()
		a.sixelBuilding = false
		if ok {
			a.SixelData = output
			a.Mode = DisplaySixel
		}
		a.mu.Unlock()
		if ok {
			storeCachedSixelData(path, output)
		}
	}(img, imagePath)
}

func buildSixelData(img image.Image) ([]byte, bool) {
	if img == nil {
		return nil, false
	}

	var pngData bytes.Buffer
	if err := imageEncodePNG(&pngData, img); err != nil {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), sixelBuildTimeout)
	defer cancel()

	sixelCmd := exec.CommandContext(
		ctx,
		"convert",
		"png:-",
		"-filter", "Lanczos",
		"-resize", "256x256^",
		"-gravity", "center",
		"-extent", "256x256",
		"sixel:-",
	)
	sixelCmd.Stdin = &pngData
	output, err := sixelCmd.Output()
	if ctx.Err() != nil || err != nil || len(output) == 0 {
		return nil, false
	}
	return output, true
}
