package artwork

import (
	"math"

	"virga-player/settings/theme"

	"github.com/gdamore/tcell/v2"
)

func (a *Artwork) drawDotsAnimation(screen tcell.Screen, centerX, y, maxWidth int, progress float64, currentTheme theme.Theme) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	screenWidth, _ := screen.Size()
	if screenWidth <= 0 {
		return
	}
	if maxWidth <= 0 {
		maxWidth = screenWidth - 2
	}

	bandWidth := maxWidth
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

	innerPadding := lyricsBasePadding
	maxPadding := maxInt(1, bandWidth/4)
	if innerPadding > maxPadding {
		innerPadding = maxPadding
	}
	textWidth := bandWidth - innerPadding*2
	if textWidth < 4 {
		textWidth = bandWidth
		innerPadding = 0
	}

	drawX := bandX + innerPadding
	if textWidth <= 0 {
		return
	}

	const animRadius = 30
	desiredWidth := animRadius*2 + 1

	animWidth := textWidth
	if animWidth > desiredWidth {
		animWidth = desiredWidth
	}
	if animWidth < 6 {
		if textWidth >= 6 {
			animWidth = 6
		} else {
			animWidth = textWidth
		}
	}

	animX := drawX + (textWidth-animWidth)/2
	if animX < drawX {
		animX = drawX
	}
	if animX+animWidth > drawX+textWidth {
		animX = drawX + textWidth - animWidth
	}

	a.fillLine(screen, animX, y, animWidth, currentTheme.LyricsBackground)

	centerPos := animX + animWidth/2

	leftCount := 5
	dotColor := currentTheme.LyricsPulse
	bg := currentTheme.LyricsBackground
	leftHalf := animWidth / 2
	if leftHalf < 1 {
		leftHalf = 1
	}

	for j := 0; j < leftCount; j++ {
		var leftInit int
		if leftCount == 1 {
			leftInit = animX
		} else {
			leftInit = animX + int(math.Round(float64(j)*float64(leftHalf-1)/float64(leftCount-1)))
		}
		leftFinal := centerPos - (leftCount - j)
		pos := leftInit + int(math.Round(float64(leftFinal-leftInit)*progress))
		if pos < animX {
			pos = animX
		}
		if pos >= animX+animWidth {
			pos = animX + animWidth - 1
		}
		screen.SetContent(pos, y, '.', nil, tcell.StyleDefault.Foreground(dotColor).Background(bg))
		var rightInit int
		if leftCount == 1 {
			rightInit = animX + animWidth - 1
		} else {
			rightInit = animX + animWidth - 1 - int(math.Round(float64(j)*float64(leftHalf-1)/float64(leftCount-1)))
		}
		rightFinal := centerPos + (j + 1)
		rpos := rightInit + int(math.Round(float64(rightFinal-rightInit)*progress))
		if rpos < animX {
			rpos = animX
		}
		if rpos >= animX+animWidth {
			rpos = animX + animWidth - 1
		}
		screen.SetContent(rpos, y, '.', nil, tcell.StyleDefault.Foreground(dotColor).Background(bg))
	}
}
