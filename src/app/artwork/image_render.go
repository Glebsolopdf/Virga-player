package artwork

import (
	"image"
	"image/color"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

// drawImageInBox
func (a *Artwork) drawImageInBox(screen tcell.Screen, x, y, w, h int) {
	img := a.getCoverImg()
	if img == nil {
		return
	}

	screenW, screenH := screen.Size()
	b := img.Bounds()
	srcW := b.Dx()
	srcH := b.Dy()
	if srcW <= 0 || srcH <= 0 || w <= 0 || h <= 0 {
		return
	}
	theme := settings.CurrentTheme()
	bgR, bgG, bgB := theme.Background.RGB()
	bg := color.NRGBA{R: uint8(bgR), G: uint8(bgG), B: uint8(bgB), A: 255}

	cropSize := srcW
	if srcH < cropSize {
		cropSize = srcH
	}
	cropMinX := b.Min.X + (srcW-cropSize)/2
	cropMinY := b.Min.Y + (srcH-cropSize)/2
	targetPxH := h * 2

	for row := 0; row < h; row++ {
		topSrcY0 := cropMinY + ((row*2)*cropSize)/targetPxH
		topSrcY1 := cropMinY + (((row*2)+1)*cropSize)/targetPxH
		if topSrcY1 <= topSrcY0 {
			topSrcY1 = topSrcY0 + 1
		}
		botSrcY0 := cropMinY + (((row*2)+1)*cropSize)/targetPxH
		botSrcY1 := cropMinY + (((row*2)+2)*cropSize)/targetPxH
		if botSrcY1 <= botSrcY0 {
			botSrcY1 = botSrcY0 + 1
		}

		for col := 0; col < w; col++ {
			srcX0 := cropMinX + (col*cropSize)/w
			srcX1 := cropMinX + ((col+1)*cropSize)/w
			if srcX1 <= srcX0 {
				srcX1 = srcX0 + 1
			}

			xPos := x + col
			yPos := y + row
			if xPos < 0 || xPos >= screenW || yPos < 0 || yPos >= screenH {
				continue
			}

			topColor := rgbaToTcell(sampleCenterColor(img, srcX0, topSrcY0, srcX1, topSrcY1), bg, a.Fade, a.Pulse)
			botColor := rgbaToTcell(sampleCenterColor(img, srcX0, botSrcY0, srcX1, botSrcY1), bg, a.Fade, a.Pulse)
			style := tcell.StyleDefault.Foreground(topColor).Background(botColor)
			screen.SetContent(xPos, yPos, '▀', nil, style)
		}
	}
}

func rgbaToTcell(c color.Color, background color.NRGBA, fade, pulse float64) tcell.Color {
	src := color.NRGBAModel.Convert(c).(color.NRGBA)
	alpha := float64(src.A) / 255.0

	rf := (float64(src.R)*alpha + float64(background.R)*(1-alpha)) / 255.0
	gf := (float64(src.G)*alpha + float64(background.G)*(1-alpha)) / 255.0
	bf := (float64(src.B)*alpha + float64(background.B)*(1-alpha)) / 255.0

	if fade < 0 {
		fade = 0
	}
	if fade > 1 {
		fade = 1
	}
	rf *= fade
	gf *= fade
	bf *= fade

	if pulse > 0 {
		rf += (1 - rf) * pulse * 0.28
		gf += (1 - gf) * pulse * 0.28
		bf += (1 - bf) * pulse * 0.28
	}

	return tcell.NewRGBColor(int32(clampFloat(rf*255, 0, 255)), int32(clampFloat(gf*255, 0, 255)), int32(clampFloat(bf*255, 0, 255)))
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func sampleCenterColor(img image.Image, x0, y0, x1, y1 int) color.Color {
	b := img.Bounds()
	if x0 < b.Min.X {
		x0 = b.Min.X
	}
	if y0 < b.Min.Y {
		y0 = b.Min.Y
	}
	if x1 > b.Max.X {
		x1 = b.Max.X
	}
	if y1 > b.Max.Y {
		y1 = b.Max.Y
	}
	if x1 <= x0 {
		x1 = x0 + 1
	}
	if y1 <= y0 {
		y1 = y0 + 1
	}
	cx := x0 + (x1-x0)/2
	cy := y0 + (y1-y0)/2
	if cx >= b.Max.X {
		cx = b.Max.X - 1
	}
	if cy >= b.Max.Y {
		cy = b.Max.Y - 1
	}
	return img.At(cx, cy)
}
