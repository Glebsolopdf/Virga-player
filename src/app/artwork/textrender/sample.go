package textrender

import (
	"image"
	"image/color"
)

func sampleCenterColor(img image.Image, x0, y0, x1, y1 int) color.NRGBA {
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
	return color.NRGBAModel.Convert(img.At(cx, cy)).(color.NRGBA)
}
