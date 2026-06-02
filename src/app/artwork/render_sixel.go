package artwork

import (
	"fmt"
	"image"
	"os"

	sixeldata "virga-player/app/artwork/sixeldata"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

const (
	sixelImageWidth      = 40
	sixelImageHeight     = 24
	sixelTextWidth       = 40
	sixelTitleOffsetY    = 15
	sixelArtistOffsetY   = 17
	sixelAlbumOffsetY    = 19
	sixelTimelineOffsetY = 21
	sixelTimeOffsetY     = 23
	sixelPlayerBodyH     = 24
)

func (a *Artwork) renderSixel(screen tcell.Screen, state artworkSnapshot, drawInfo, drawLyrics bool) {
	theme := settings.CurrentTheme()
	if state.coverImg == nil {
		a.renderTextOnly(screen, state, drawInfo, drawLyrics)
		return
	}

	hasSixel := len(state.sixelData) > 0
	if !hasSixel {
		if drawInfo {
			a.prepareSixelDataAsync()
		}
		a.renderTextOnly(screen, state, drawInfo, drawLyrics)
		return
	}

	if drawInfo {
		screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.TrackTitle))
	}

	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	lyricsLines := lyricDisplayWindow(state.lyrics, state.elapsed, state.pulse, theme)
	if !drawLyrics {
		lyricsLines = nil
	}
	lyricsGap := lyricsVerticalGap(lyricsLines)
	totalH := sixelPlayerBodyH + lyricsGap + len(lyricsLines)
	blockTop := centerY - totalH/2

	imgY := blockTop
	imgX := centerX - sixelImageWidth/2
	if imgX < 0 {
		imgX = 0
	}
	if imgY < 0 {
		imgY = 0
	}

	if drawInfo && a.shouldRenderSixelAt(imgX, imgY, state.sixelData) {
		fmt.Printf("\x1B[%d;%dH", imgY+1, imgX+1)
		_, _ = os.Stdout.Write(state.sixelData)
	}

	textWidth := sixelTextWidth
	if textWidth > w-10 {
		textWidth = w - 10
	}
	textStartX := centerX - textWidth/2
	if drawInfo {
		a.drawCenteredInArea(screen, textStartX, textWidth, blockTop+sixelTitleOffsetY, state.title, theme.TrackTitle)
		a.drawCenteredInArea(screen, textStartX, textWidth, blockTop+sixelArtistOffsetY, state.artist, theme.TrackArtist)
		a.drawCenteredInArea(screen, textStartX, textWidth, blockTop+sixelAlbumOffsetY, state.album, theme.TrackAlbum)
	}

	lyricsY := blockTop + sixelPlayerBodyH + lyricsGap
	if drawInfo {
		a.drawTimeline(screen, centerX, blockTop+sixelTimelineOffsetY, 28, state.elapsed, state.duration)
		timeStr := formatTime(state.elapsed) + " / " + formatTime(state.duration)
		if state.duration <= 0 {
			timeStr = "00:00 / 00:00"
		}
		a.drawText(screen, centerX-len(timeStr)/2, blockTop+sixelTimeOffsetY, timeStr, theme.TrackTime)
	}

	if len(lyricsLines) > 0 {
		lyricsMaxWidth := w - 2
		if lyricsMaxWidth < 16 {
			lyricsMaxWidth = 16
		}
		a.drawLyricsBlock(screen, centerX, lyricsMaxWidth, lyricsY, lyricsLines, state.pulse, theme)
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
