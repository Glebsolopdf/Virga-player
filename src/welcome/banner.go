package welcome

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func (m model) renderBanner() string {
	state := m.bannerAnim.state(m.frame)
	innerWidth := fitWidth(m.width-10, 34, 54)
	textProgress := easeOut(state.textProgress)
	titleColor := fadeColor(
		colorful.LinearRgb(0.28, 0.30, 0.34),
		colorful.LinearRgb(1, 1, 1),
		textProgress,
	)
	subtitleColor := fadeColor(
		colorful.LinearRgb(0.22, 0.24, 0.27),
		colorful.LinearRgb(0.72, 0.74, 0.78),
		textProgress,
	)

	lines := []string{
		revealText(centerText("Welcome to Virga", innerWidth), textProgress, titleColor, true),
		revealText(centerText("Made by SubstituteMe", innerWidth), textProgress, subtitleColor, false),
	}
	return renderAnimatedBox(lines, innerWidth, state, float64(m.frame), borderStartTopLeft)
}

func easeOut(value float64) float64 {
	value = math.Max(0, math.Min(1, value))
	return 1 - math.Pow(1-value, 3)
}
