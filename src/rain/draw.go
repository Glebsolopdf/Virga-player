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

			color := getLayerColor(p.Layer, theme, getDropColor(p.Opacity, theme))
			char := getDropChar(i, p.Length, p.VelX, theme)
			style := getDropStyle(color, p.Layer, i, p.Length, p.Opacity, p.MaxOpacity, ps.pulse)
			screen.SetContent(x, dropY, char, nil, style)
		}
	}
}

func getDropColor(opacity int, theme settings.Theme) tcell.Color {
	if opacity > 1 {
		return theme.RainHead
	}
	return theme.RainTail
}

func getDropStyle(color tcell.Color, layer, position, length, opacity, maxOpacity int, pulse float64) tcell.Style {
	style := tcell.StyleDefault.Foreground(color)
	pulse = clamp(pulse, 0, 1)

	switch layer {
	case layerVeryNear:
		style = style.Foreground(color).Bold(true)
	case layerNear:
		style = style.Foreground(color).Bold(true)
	case layerMid:
		style = style.Foreground(color)
	case layerFar:
		style = style.Foreground(color).Dim(true)
	case layerVeryFar:
		style = style.Foreground(color).Dim(true)
	}

	if pulse > 0.65 {
		style = style.Bold(true)
	} else if pulse < 0.22 {
		style = style.Dim(true)
	}

	opacityRatio := 0.0
	if maxOpacity > 0 {
		opacityRatio = float64(opacity) / float64(maxOpacity)
	}
	if opacityRatio < 0.3 {
		style = style.Dim(true)
	} else if opacityRatio > 0.9 {
		style = style.Bold(true)
	}

	if position < length-1 {
		tailDistance := length - 1 - position
		if tailDistance > 1 && pulse < 0.45 {
			style = style.Dim(true)
		}
	}
	return style
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

func getLayerColor(layer int, theme settings.Theme, fallback tcell.Color) tcell.Color {
	switch layer {
	case layerVeryNear:
		return theme.RainVeryNear
	case layerNear:
		return theme.RainNear
	case layerMid:
		return theme.RainMid
	case layerFar:
		return theme.RainFar
	case layerVeryFar:
		return theme.RainVeryFar
	default:
		return fallback
	}
}
