package artwork

import (
	"image"
	"image/color"
	textrender "virga-player/app/artwork/textrender"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

// drawImageInBox
func (a *Artwork) drawImageInBox(screen tcell.Screen, x, y, w, h int, img image.Image, fade, pulse float64) {
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
	cells := a.cachedTextRenderCells(img, w, h)

	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			xPos := x + col
			yPos := y + row
			if xPos < 0 || xPos >= screenW || yPos < 0 || yPos >= screenH {
				continue
			}

			cell := cells[row*w+col]
			topColor := rgbaToTcell(cell.Top, bg, fade, pulse)
			botColor := rgbaToTcell(cell.Bottom, bg, fade, pulse)
			style := tcell.StyleDefault.Foreground(topColor).Background(botColor)
			screen.SetContent(xPos, yPos, theme.ArtworkBlockRune, nil, style)
		}
	}
}

func (a *Artwork) cachedTextRenderCells(img image.Image, w, h int) []textrender.CellSample {
	a.mu.RLock()
	if a.textRenderCache.Width == w && a.textRenderCache.Height == h && len(a.textRenderCache.Cells) == w*h {
		cells := a.textRenderCache.Cells
		a.mu.RUnlock()
		return cells
	}
	a.mu.RUnlock()

	cells := textrender.BuildCells(img, w, h)

	a.mu.Lock()
	a.textRenderCache.Width = w
	a.textRenderCache.Height = h
	a.textRenderCache.Cells = cells
	a.mu.Unlock()

	return cells
}

func rgbaToTcell(src color.NRGBA, background color.NRGBA, fade, pulse float64) tcell.Color {
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
