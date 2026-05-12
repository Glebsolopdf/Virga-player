package artwork

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"os/exec"
	"sync"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

var (
	sixelDataCacheMu sync.RWMutex
	sixelDataCache   = map[string][]byte{}
)

func (a *Artwork) renderSixel(screen tcell.Screen) {
	theme := settings.CurrentTheme()
	if a.getCoverImg() == nil {
		a.renderTextOnly(screen)
		return
	}

	a.mu.RLock()
	hasSixel := len(a.SixelData) > 0
	a.mu.RUnlock()
	if !hasSixel {
		a.prepareSixelDataAsync()
		a.renderTextOnly(screen)
		return
	}

	screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.TrackTitle))

	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	centerX += int(math.Round(a.RainOffsetX))
	centerY += int(math.Round(a.RainOffsetY))

	imgY := centerY - 12
	imgX := centerX - 20
	if imgX < 0 {
		imgX = 0
	}
	if imgY < 0 {
		imgY = 0
	}

	fmt.Printf("\x1B[%d;%dH", imgY+1, imgX+1)
	os.Stdout.Write(a.SixelData)

	textWidth := 40
	if textWidth > w-10 {
		textWidth = w - 10
	}
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+3, a.Title, theme.TrackTitle)
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+5, a.Artist, theme.TrackArtist)
	a.drawCenteredInArea(screen, centerX-textWidth/2, textWidth, centerY+7, a.Album, theme.TrackAlbum)

	if a.Duration > 0 {
		a.drawTimeline(screen, centerX, centerY+9, 28)
		timeStr := formatTime(a.Elapsed) + " / " + formatTime(a.Duration)
		a.drawText(screen, centerX-len(timeStr)/2, centerY+11, timeStr, theme.TrackTime)
	}
}

func getCachedSixelData(imagePath string) ([]byte, bool) {
	if imagePath == "" {
		return nil, false
	}
	sixelDataCacheMu.RLock()
	data, ok := sixelDataCache[imagePath]
	sixelDataCacheMu.RUnlock()
	if !ok || len(data) == 0 {
		return nil, false
	}
	copyData := make([]byte, len(data))
	copy(copyData, data)
	return copyData, true
}

func storeCachedSixelData(imagePath string, data []byte) {
	if imagePath == "" || len(data) == 0 {
		return
	}
	copyData := make([]byte, len(data))
	copy(copyData, data)
	sixelDataCacheMu.Lock()
	sixelDataCache[imagePath] = copyData
	sixelDataCacheMu.Unlock()
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

	sixelCmd := exec.Command(
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
	if err != nil || len(output) == 0 {
		return nil, false
	}
	return output, true
}
