package artwork

import (
	"math"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

// renderTextOnly
func (a *Artwork) renderTextOnly(screen tcell.Screen) {
	theme := settings.CurrentTheme()
	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	centerX += int(math.Round(a.RainOffsetX))
	centerY += int(math.Round(a.RainOffsetY))
	infoBlockH := 7

	hasCover := a.getCoverImg() != nil
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
	if img := a.getCoverImg(); img != nil {
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

	contentH := infoBlockH
	if hasCover && boxH > contentH {
		contentH = boxH
	}
	contentY := centerY - contentH/2
	boxY := contentY
	boxX := 0
	infoX := centerX - infoW/2
	infoY := centerY - infoBlockH/2

	if hasCover {
		totalW := boxW + infoGap + infoW
		boxX = centerX - totalW/2
		infoX = boxX + boxW + infoGap
		infoY = boxY + (boxH-infoBlockH)/2

		a.drawImageInBox(screen, boxX, boxY, coverInnerW, coverInnerH)
	}

	a.drawCenteredInArea(screen, infoX, infoW, infoY, a.Title, theme.TrackTitle)
	a.drawCenteredInArea(screen, infoX, infoW, infoY+1, a.Artist, theme.TrackArtist)
	a.drawCenteredInArea(screen, infoX, infoW, infoY+2, a.Album, theme.TrackAlbum)

	barWidth := infoW + 6
	if barWidth > w-10 {
		barWidth = w - 10
	}
	if barWidth < 16 {
		barWidth = 16
	}
	if a.Duration > 0 {
		barCenterX := infoX + infoW/2
		a.drawTimeline(screen, barCenterX, infoY+4, barWidth)
		timeStr := formatTime(a.Elapsed) + " / " + formatTime(a.Duration)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+6, timeStr, theme.TrackTime)
	} else {
		timeStr := formatTime(a.Elapsed) + " / --:--"
		a.drawTimeline(screen, infoX+infoW/2, infoY+4, barWidth)
		a.drawCenteredInArea(screen, infoX, infoW, infoY+6, timeStr, theme.TrackTime)
	}
}
