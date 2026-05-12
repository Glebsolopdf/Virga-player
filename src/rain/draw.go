package rain

import (
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (ps *ParticleSystem) Draw(screen tcell.Screen) {
	if !ps.enabled {
		return
	}
	theme := settings.CurrentTheme()
	for _, p := range ps.particles {
		x := int(p.X)
		y := int(p.Y)

		if x < 0 || x >= ps.width || y < 0 || y >= ps.height {
			continue
		}

		for i := 0; i < p.Length; i++ {
			dropY := y + i
			if dropY >= ps.height {
				break
			}

			color := getDropColor(p.Opacity, theme)
			char := getDropChar(i, p.Length, p.VelX, theme)
			screen.SetContent(x, dropY, char, nil, tcell.StyleDefault.Foreground(color))
		}
	}
}

func getDropColor(opacity int, theme settings.Theme) tcell.Color {
	if opacity > 1 {
		return theme.RainHead
	}
	return theme.RainTail
}

func getDropChar(position, length int, velX float64, theme settings.Theme) rune {
	if velX > 0 {
		return theme.RainRightRune
	}
	if velX < 0 {
		return theme.RainLeftRune
	}
	if position == length-1 {
		return theme.RainHeadRune
	}
	return theme.RainBodyRune
}
