package artwork

import (
	"fmt"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (a *Artwork) drawTimeline(screen tcell.Screen, centerX, y, width int, elapsed, duration int) {
	theme := settings.CurrentTheme()
	w, h := screen.Size()
	if y < 0 || y >= h {
		return
	}

	barWidth := width
	if barWidth > w-10 {
		barWidth = w - 10
	}
	if barWidth < 8 {
		barWidth = 8
	}

	barX := centerX - barWidth/2

	var currentPos int
	if duration > 0 {
		currentPos = (elapsed * barWidth) / duration
	}

	if currentPos < 0 {
		currentPos = 0
	}
	if currentPos >= barWidth {
		currentPos = barWidth - 1
	}

	for i := 0; i < barWidth; i++ {
		var ch rune
		var color tcell.Color

		if i == 0 || i == barWidth-1 {
			if i == 0 {
				ch = theme.TimelineLeftRune
			} else {
				ch = theme.TimelineRightRune
			}
			color = theme.TimelineBracket
		} else if i < currentPos {
			ch = theme.TimelinePlayedRune
			color = theme.TimelinePlayed
		} else if i == currentPos {
			ch = theme.TimelineCurrentRune
			color = theme.TimelineCurrent
		} else {
			ch = theme.TimelineEmptyRune
			color = theme.TimelineRemaining
		}

		if barX+i >= 0 && barX+i < w {
			screen.SetContent(barX+i, y, ch, nil, tcell.Style{}.Foreground(color))
		}
	}
}

func (a *Artwork) truncateText(text string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	if maxLen < 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

func (a *Artwork) drawCenteredInArea(screen tcell.Screen, x, w, y int, text string, color tcell.Color) {
	truncated := a.truncateText(text, w)
	tx := x + (w-len([]rune(truncated)))/2
	a.drawText(screen, tx, y, truncated, color)
}

func (a *Artwork) drawText(screen tcell.Screen, x, y int, text string, color tcell.Color) {
	w, h := screen.Size()
	if y < 0 || y >= h {
		return
	}
	for offset, ch := range []rune(text) {
		posX := x + offset
		if posX >= 0 && posX < w {
			screen.SetContent(posX, y, ch, nil, tcell.Style{}.Foreground(color))
		}
		if posX >= w {
			break
		}
	}
}

func (a *Artwork) drawTextWithBackground(screen tcell.Screen, x, y int, text string, foreground, background tcell.Color) {
	w, h := screen.Size()
	if y < 0 || y >= h {
		return
	}
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	for offset, ch := range []rune(text) {
		posX := x + offset
		if posX >= 0 && posX < w {
			screen.SetContent(posX, y, ch, nil, style)
		}
		if posX >= w {
			break
		}
	}
}

func (a *Artwork) fillLine(screen tcell.Screen, x, y, width int, background tcell.Color) {
	w, h := screen.Size()
	if y < 0 || y >= h || width <= 0 {
		return
	}
	startX := maxInt(0, x)
	endX := minInt(w, x+width)
	style := tcell.StyleDefault.Background(background).Foreground(background)
	for posX := startX; posX < endX; posX++ {
		screen.SetContent(posX, y, ' ', nil, style)
	}
}

func formatTime(seconds int) string {
	mins := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}
