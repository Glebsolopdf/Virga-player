package artwork

import (
	"math"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

// renderTextOnly
func (a *Artwork) renderTextOnly(screen tcell.Screen, state artworkSnapshot, drawInfo, drawLyrics bool) {
	theme := settings.CurrentTheme()
	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	centerX += int(math.Round(state.rainOffsetX))
	centerY += int(math.Round(state.rainOffsetY))
	lyricsLines := lyricDisplayWindow(state.lyrics, state.elapsed, state.pulse, theme)
	if !drawLyrics {
		lyricsLines = nil
	}
	lyricsGap := lyricsVerticalGap(lyricsLines)
	infoBlockH := 7

	hasCover := state.coverImg != nil
	infoGap := 4
	infoW := 34
	if w >= 100 {
		infoW = 28
	}
	if w >= 120 {
		infoW = 26
	}
	if w >= 150 {
		infoW = 24
	}

	if infoW > w/3 {
		infoW = w / 3
	}

	coverInnerH := 10
	coverInnerW := coverInnerH * 2
	if img := state.coverImg; img != nil {
		imgBounds := img.Bounds()
		if imgBounds.Dx() > 0 && imgBounds.Dy() > 0 {
			aspect := float64(imgBounds.Dx()) / float64(imgBounds.Dy())
			coverInnerH = imgBounds.Dy() / 24
			if coverInnerH < 10 {
				coverInnerH = 10
			}
			coverInnerW = int(float64(coverInnerH) * aspect * 1.8)
			if coverInnerW < coverInnerH*2 {
				coverInnerW = coverInnerH * 2
			}
		}
	}

	maxCoverByWidth := w - infoW - infoGap - 4
	maxCoverByHeight := h - 4

	if coverInnerW > maxCoverByWidth {
		coverInnerW = maxCoverByWidth
		coverInnerH = coverInnerW / 2
	}
	if coverInnerW > w/3 {
		coverInnerW = w / 3
		coverInnerH = coverInnerW / 2
	}
	if coverInnerH > maxCoverByHeight {
		coverInnerH = maxCoverByHeight
		coverInnerW = coverInnerH * 2
	}

	boxW := coverInnerW + 2
	boxH := coverInnerH + 2
	if maxCoverByHeight < coverInnerH || maxCoverByWidth < coverInnerW || coverInnerH < 8 {
		hasCover = false
	}

	bodyH := infoBlockH
	if hasCover && boxH > bodyH {
		bodyH = boxH
	}
	totalH := bodyH + lyricsGap + len(lyricsLines)
	contentY := centerY - totalH/2
	boxY := contentY + (bodyH-boxH)/2
	boxX := 0
	infoX := centerX - infoW/2
	infoY := contentY + (bodyH-infoBlockH)/2

	if hasCover {
		totalW := boxW + infoGap + infoW
		boxX = centerX - totalW/2
		infoX = boxX + boxW + infoGap

		if drawInfo {
			a.drawImageInBox(screen, boxX, boxY, coverInnerW, coverInnerH, state.coverImg, state.fade, state.pulse)
		}
	}

	if drawInfo {
		a.drawCenteredInArea(screen, infoX, infoW, infoY, state.title, theme.TrackTitle)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+1, state.artist, theme.TrackArtist)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+2, state.album, theme.TrackAlbum)
	}

	barWidth := infoW + 6
	if barWidth > w-10 {
		barWidth = w - 10
	}
	if barWidth < 16 {
		barWidth = 16
	}
	if drawInfo && state.duration > 0 {
		barCenterX := infoX + infoW/2
		a.drawTimeline(screen, barCenterX, infoY+4, barWidth, state.elapsed, state.duration)
		timeStr := formatTime(state.elapsed) + " / " + formatTime(state.duration)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+6, timeStr, theme.TrackTime)
	} else if drawInfo {
		timeStr := formatTime(state.elapsed) + " / --:--"
		a.drawTimeline(screen, infoX+infoW/2, infoY+4, barWidth, state.elapsed, state.duration)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+6, timeStr, theme.TrackTime)
	}

	if len(lyricsLines) > 0 {
		lyricsMaxWidth := w - 2
		if lyricsMaxWidth < 16 {
			lyricsMaxWidth = 16
		}
		a.drawLyricsBlock(screen, centerX, lyricsMaxWidth, contentY+bodyH+lyricsGap, lyricsLines, state.pulse, theme)
	}
}
