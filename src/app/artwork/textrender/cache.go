package textrender

import (
	"image"
	"image/color"
)

type CellSample struct {
	Top    color.NRGBA
	Bottom color.NRGBA
}

type Cache struct {
	Width  int
	Height int
	Cells  []CellSample
}

func BuildCells(img image.Image, w, h int) []CellSample {
	b := img.Bounds()
	srcW := b.Dx()
	srcH := b.Dy()
	if srcW <= 0 || srcH <= 0 || w <= 0 || h <= 0 {
		return nil
	}

	cropSize := srcW
	if srcH < cropSize {
		cropSize = srcH
	}
	cropMinX := b.Min.X + (srcW-cropSize)/2
	cropMinY := b.Min.Y + (srcH-cropSize)/2
	targetPxH := h * 2
	cells := make([]CellSample, 0, w*h)

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

			cells = append(cells, CellSample{
				Top:    sampleCenterColor(img, srcX0, topSrcY0, srcX1, topSrcY1),
				Bottom: sampleCenterColor(img, srcX0, botSrcY0, srcX1, botSrcY1),
			})
		}
	}

	return cells
}
