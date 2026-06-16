package playerlyrics

import "strings"

const (
	appearCenter = iota
	appearLeft
	appearRight
)

const (
	exitRain = iota
	exitSink
	exitShear
)

func appearText(text string, progress float64, variant int) string {
	if progress >= 1 {
		return text
	}
	runes := []rune(text)
	visible := int(float64(len(runes)) * progress)
	if visible < 1 {
		visible = 1
	}
	switch variant {
	case appearLeft:
		return string(runes[:visible])
	case appearRight:
		return strings.Repeat(" ", len(runes)-visible) + string(runes[len(runes)-visible:])
	default:
		pad := strings.Repeat(" ", (len(runes)-visible)/2)
		return pad + string(runes[:visible])
	}
}

func washText(text string, progress float64) string {
	if progress >= 1 {
		return ""
	}
	runes := []rune(text)
	step := int(progress * 5)
	for i := step; i < len(runes); i += 4 {
		runes[i] = ' '
	}
	return string(runes)
}

func exitFrame(text string, progress float64, variant int) (string, int, int) {
	switch variant {
	case exitSink:
		return strings.TrimRight(washText(text, progress), " "), 0, int(progress * 6)
	case exitShear:
		return washText(text, progress), int(progress * 4), int(progress * 2)
	default:
		return washText(text, progress), 0, int(progress * 4)
	}
}
