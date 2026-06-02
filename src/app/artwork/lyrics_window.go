package artwork

import (
	"sort"

	"virga-player/settings/theme"

	"github.com/gdamore/tcell/v2"
)

func lyricDisplayWindow(cues []lyricCue, elapsedSeconds int, currentPulse float64, currentTheme theme.Theme) []lyricDisplayLine {
	if len(cues) == 0 {
		return nil
	}

	elapsedMillis := elapsedSeconds * 1000
	if elapsedMillis < 0 {
		elapsedMillis = 0
	}

	current := currentLyricIndex(cues, elapsedMillis)
	if current < 0 {
		limit := minInt(len(cues), maxVisibleLyrics)
		lines := make([]lyricDisplayLine, 0, limit)
		for i := 0; i < limit; i++ {
			lines = append(lines, lyricDisplayLine{text: cues[i].text, color: currentTheme.LyricsInactive})
		}
		return lines
	}

	start := maxInt(0, current-1)
	end := minInt(len(cues), current+2)
	if end-start < maxVisibleLyrics {
		start = maxInt(0, end-maxVisibleLyrics)
	}

	lines := make([]lyricDisplayLine, 0, end-start)
	for i := start; i < end; i++ {
		lineColor := currentTheme.LyricsInactive
		active := i == current
		if active {
			lineColor = pulsingLyricColor(currentTheme, currentPulse)
		}
		lines = append(lines, lyricDisplayLine{text: cues[i].text, color: lineColor, active: active})

		if active && i+1 < len(cues) {
			nextAt := cues[i+1].atMillis
			totalGap := nextAt - cues[i].atMillis
			if totalGap > 15000 && elapsedMillis < nextAt {
				passed := elapsedMillis - cues[i].atMillis
				if passed < 0 {
					passed = 0
				}
				var progress float64
				if totalGap > 0 {
					progress = float64(passed) / float64(totalGap)
					if progress < 0 {
						progress = 0
					}
					if progress > 1 {
						progress = 1
					}
				}
				lines = append(lines, lyricDisplayLine{text: "", color: currentTheme.LyricsInactive, active: false, hasAnim: true, anim: progress})
			}
		}
	}
	return lines
}

func currentLyricIndex(cues []lyricCue, elapsedMillis int) int {
	index := sort.Search(len(cues), func(i int) bool {
		return cues[i].atMillis > elapsedMillis
	})
	return index - 1
}

func pulsingLyricColor(currentTheme theme.Theme, currentPulse float64) tcell.Color {
	if currentPulse >= 0.35 {
		return currentTheme.LyricsPulse
	}
	return currentTheme.LyricsCurrent
}
