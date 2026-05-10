package artwork

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func (a *Artwork) loadCoverImage() {
	f, err := os.Open(a.ImagePath)
	if err != nil {
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	a.CoverImg = downscaleCoverImage(img)
}

func downscaleCoverImage(img image.Image) image.Image {
	maxDim := 256
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= maxDim && height <= maxDim {
		return img
	}

	scale := float64(maxDim) / float64(width)
	if float64(height)*scale > float64(maxDim) {
		scale = float64(maxDim) / float64(height)
	}

	newW := int(float64(width) * scale)
	newH := int(float64(height) * scale)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
	return resized
}

func imageEncodePNG(buf *bytes.Buffer, img image.Image) error {
	if img == nil {
		return fmt.Errorf("nil image")
	}
	return png.Encode(buf, img)
}
