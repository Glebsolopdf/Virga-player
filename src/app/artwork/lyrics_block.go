package artwork

import (
	"strings"

	"virga-player/settings/theme"

	"github.com/gdamore/tcell/v2"
)

func (a *Artwork) drawLyricsBlock(screen tcell.Screen, centerX, maxWidth, startY int, lines []lyricDisplayLine, pulse float64, currentTheme theme.Theme) {
	if len(lines) == 0 {
		return
	}

	for i, line := range lines {
		if strings.TrimSpace(line.text) == "" && !line.hasAnim {
			continue
		}
		if line.hasAnim {
			a.drawDotsAnimation(screen, centerX, startY+i, maxWidth, line.anim, currentTheme)
		} else {
			a.drawLyricLine(screen, centerX, startY+i, maxWidth, line, pulse, currentTheme)
		}
	}
}

// drawLyricLine draws a single lyric line centered in a band with background.
func (a *Artwork) drawLyricLine(screen tcell.Screen, centerX, y, maxWidth int, line lyricDisplayLine, pulse float64, currentTheme theme.Theme) {
	screenWidth, _ := screen.Size()
	if screenWidth <= 0 {
		return
	}
	if maxWidth <= 0 {
		maxWidth = screenWidth - 2
	}

	extraPad := 0
	if line.active {
		extraPad = int(pulse*2 + 0.5)
	}

	contentWidth := len([]rune(line.text))
	if contentWidth < 1 {
		contentWidth = 1
	}
	bandWidth := contentWidth + (lyricsBasePadding+extraPad)*2
	if bandWidth > maxWidth {
		bandWidth = maxWidth
	}
	if bandWidth > screenWidth-2 {
		bandWidth = screenWidth - 2
	}
	if bandWidth < 8 {
		bandWidth = minInt(screenWidth, 8)
	}

	bandX := centerX - bandWidth/2
	if bandX < 0 {
		bandX = 0
	}
	if bandX+bandWidth > screenWidth {
		bandWidth = screenWidth - bandX
	}

	a.fillLine(screen, bandX, y, bandWidth, currentTheme.LyricsBackground)

	innerPadding := lyricsBasePadding + extraPad
	maxPadding := maxInt(1, bandWidth/4)
	if innerPadding > maxPadding {
		innerPadding = maxPadding
	}
	textWidth := bandWidth - innerPadding*2
	if textWidth < 4 {
		textWidth = bandWidth
		innerPadding = 0
	}

	visible, centered := lyricViewport(line.text, textWidth)
	drawX := bandX + innerPadding
	if centered {
		drawX += (textWidth - len([]rune(visible))) / 2
	}
	a.drawTextWithBackground(screen, drawX, y, visible, line.color, currentTheme.LyricsBackground)
}

// lyricsVerticalGap determines vertical gap between rendered lyric lines.
func lyricsVerticalGap(lines []lyricDisplayLine) int {
	if len(lines) == 0 {
		return 0
	}
	gap := theme.CurrentTheme().LyricsGap
	if gap < 0 {
		gap = lyricsDefaultGap
	}
	if gap > 8 {
		return 8
	}
	return gap
}
