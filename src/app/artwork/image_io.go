package artwork

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	imagedraw "image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	xdraw "golang.org/x/image/draw"
)

func (a *Artwork) loadCoverImage() {
	var img image.Image

	if strings.HasPrefix(a.ImagePath, "http://") || strings.HasPrefix(a.ImagePath, "https://") {
		client := http.Client{Timeout: 4 * time.Second}
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.ImagePath, nil)
		if err != nil {
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		img, _, err = image.Decode(io.LimitReader(resp.Body, 8*1024*1024))
		if err != nil {
			return
		}
	} else {
		f, err := os.Open(a.ImagePath)
		if err != nil {
			return
		}
		defer f.Close()

		img, _, err = image.Decode(f)
		if err != nil {
			return
		}
	}

	img = trimCoverPadding(img)
	cov := normalizeCoverImage(img)
	a.mu.Lock()
	a.CoverImg = cov
	a.mu.Unlock()
}

func trimCoverPadding(img image.Image) image.Image {
	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return img
	}

	trimRow := func(y int) bool {
		for x := b.Min.X; x < b.Max.X; x++ {
			if !isWhiteOrTransparent(img.At(x, y)) {
				return false
			}
		}
		return true
	}

	trimCol := func(x int) bool {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if !isWhiteOrTransparent(img.At(x, y)) {
				return false
			}
		}
		return true
	}

	minX, minY := b.Min.X, b.Min.Y
	maxX, maxY := b.Max.X-1, b.Max.Y-1

	for minY <= maxY && trimRow(minY) {
		minY++
	}
	for maxY >= minY && trimRow(maxY) {
		maxY--
	}
	for minX <= maxX && trimCol(minX) {
		minX++
	}
	for maxX >= minX && trimCol(maxX) {
		maxX--
	}

	if minX > maxX || minY > maxY {
		return img
	}

	crop := image.Rect(minX, minY, maxX+1, maxY+1)
	if crop.Eq(b) {
		return img
	}

	trimmed := image.NewRGBA(image.Rect(0, 0, crop.Dx(), crop.Dy()))
	imagedraw.Draw(trimmed, trimmed.Bounds(), img, crop.Min, imagedraw.Src)
	return trimmed
}

func isWhiteOrTransparent(source color.Color) bool {
	c := color.NRGBAModel.Convert(source).(color.NRGBA)
	if c.A <= 8 {
		return true
	}
	return c.R >= 232 && c.G >= 232 && c.B >= 232
}

func normalizeCoverImage(img image.Image) image.Image {
	const targetDim = 256
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return img
	}

	if width == targetDim && height == targetDim {
		return img
	}

	cropSize := width
	if height < cropSize {
		cropSize = height
	}
	srcMinX := bounds.Min.X + (width-cropSize)/2
	srcMinY := bounds.Min.Y + (height-cropSize)/2
	srcRect := image.Rect(srcMinX, srcMinY, srcMinX+cropSize, srcMinY+cropSize)

	normalized := image.NewRGBA(image.Rect(0, 0, targetDim, targetDim))
	xdraw.CatmullRom.Scale(normalized, normalized.Bounds(), img, srcRect, xdraw.Over, nil)

	if cropSize < targetDim {
		return normalized
	}

	return normalized
}

func imageEncodePNG(buf *bytes.Buffer, img image.Image) error {
	if img == nil {
		return fmt.Errorf("nil image")
	}
	return png.Encode(buf, img)
}
