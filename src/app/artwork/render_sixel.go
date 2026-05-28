package artwork

import (
	"fmt"
	"image"
	"math"
	"os"

	sixeldata "virga-player/app/artwork/sixeldata"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

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
		output, ok := sixeldata.Build(localImg)
		a.mu.Lock()
		a.sixelBuilding = false
		if ok {
			a.SixelData = output
			a.Mode = DisplaySixel
		}
		a.mu.Unlock()
		if ok {
			sixeldata.Store(path, output)
		}
	}(img, imagePath)
}
