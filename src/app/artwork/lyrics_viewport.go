package artwork

import "time"

func lyricViewport(text string, width int) (string, bool) {
	runes := []rune(text)
	if width <= 0 || len(runes) == 0 {
		return "", true
	}
	if len(runes) <= width {
		return string(runes), true
	}

	gap := maxInt(width/2, lyricsBasePadding+2)
	marquee := make([]rune, 0, len(runes)+gap*2)
	for i := 0; i < gap; i++ {
		marquee = append(marquee, ' ')
	}
	marquee = append(marquee, runes...)
	for i := 0; i < gap; i++ {
		marquee = append(marquee, ' ')
	}

	cycleLen := len(marquee) - width + 1
	if cycleLen <= 1 {
		return string(runes[:width]), false
	}
	step := int(time.Now().UnixMilli() / int64(lyricsScrollStep))
	start := step % cycleLen
	return string(marquee[start : start+width]), false
}
