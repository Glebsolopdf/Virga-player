package artwork

import (
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

// renderTextOnly выводит текстовый интерфейс с обложкой слева.
func (a *Artwork) renderTextOnly(screen tcell.Screen) {
	theme := settings.CurrentTheme()
	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2
	infoBlockH := 7

	hasCover := a.CoverImg != nil
	infoGap := 4
	infoW := 34

	// Обложка фиксированного размера для стабильности.
	coverInnerH := 10
	coverInnerW := coverInnerH * 2
	maxCoverByWidth := w - infoW - infoGap - 4
	maxCoverByHeight := h - 4

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

	if len(a.Title) > 30 {
		a.drawCenteredInArea(screen, infoX, infoW, infoY, a.Title[:30], theme.TrackTitle)
	} else {
		a.drawCenteredInArea(screen, infoX, infoW, infoY, a.Title, theme.TrackTitle)
	}

	if len(a.Artist) > 30 {
		a.drawCenteredInArea(screen, infoX, infoW, infoY+1, a.Artist[:30], theme.TrackArtist)
	} else {
		a.drawCenteredInArea(screen, infoX, infoW, infoY+1, a.Artist, theme.TrackArtist)
	}

	if len(a.Album) > 30 {
		a.drawCenteredInArea(screen, infoX, infoW, infoY+2, a.Album[:30], theme.TrackAlbum)
	} else {
		a.drawCenteredInArea(screen, infoX, infoW, infoY+2, a.Album, theme.TrackAlbum)
	}

	barWidth := infoW - 6
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
