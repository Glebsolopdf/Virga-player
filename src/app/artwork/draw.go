package artwork

import (
	"fmt"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (a *Artwork) drawTimeline(screen tcell.Screen, centerX, y, width int) {
	theme := settings.CurrentTheme()
	w, _ := screen.Size()

	barWidth := width
	if barWidth > w-10 {
		barWidth = w - 10
	}
	if barWidth < 8 {
		barWidth = 8
	}

	barX := centerX - barWidth/2

	var currentPos int
	if a.Duration > 0 {
		currentPos = (a.Elapsed * barWidth) / a.Duration
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
	if len(text) <= maxLen {
		return text
	}
	if maxLen < 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

func (a *Artwork) drawCenteredInArea(screen tcell.Screen, x, w, y int, text string, color tcell.Color) {
	truncated := a.truncateText(text, w)
	tx := x + (w-len(truncated))/2
	a.drawText(screen, tx, y, truncated, color)
}

func (a *Artwork) drawText(screen tcell.Screen, x, y int, text string, color tcell.Color) {
	w, _ := screen.Size()
	for i, ch := range text {
		posX := x + i
		if posX >= 0 && posX < w && y >= 0 {
			screen.SetContent(posX, y, ch, nil, tcell.Style{}.Foreground(color))
		}
		if posX >= w {
			break
		}
	}
}

func formatTime(seconds int) string {
	mins := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}
