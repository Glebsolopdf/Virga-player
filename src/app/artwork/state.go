package artwork

import "image/color"

func (a *Artwork) snapshot() artworkSnapshot {
	a.mu.RLock()
	defer a.mu.RUnlock()

	state := artworkSnapshot{
		mode:             a.Mode,
		imagePath:        a.ImagePath,
		coverImg:         a.CoverImg,
		title:            a.Title,
		artist:           a.Artist,
		album:            a.Album,
		duration:         a.Duration,
		elapsed:          a.Elapsed,
		animationEnabled: a.AnimationEnabled,
		fade:             a.Fade,
		pulse:            a.Pulse,
		rainOffsetX:      a.RainOffsetX,
		rainOffsetY:      a.RainOffsetY,
		sixelData:        a.SixelData,
	}
	return state
}

func (a *Artwork) UpdateTrackInfo(title, artist, album string, duration, elapsed int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Title = title
	a.Artist = artist
	a.Album = album
	a.Duration = duration
	a.Elapsed = elapsed
}

func (a *Artwork) computeAverageColor() {
	a.mu.RLock()
	img := a.CoverImg
	a.mu.RUnlock()
	if img == nil {
		a.mu.Lock()
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		a.mu.Unlock()
		return
	}

	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		a.mu.Lock()
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		a.mu.Unlock()
		return
	}

	var rSum, gSum, bSum, count float64
	stepX := (b.Dx() / 64) + 1
	stepY := (b.Dy() / 64) + 1
	for y := b.Min.Y; y < b.Max.Y; y += stepY {
		for x := b.Min.X; x < b.Max.X; x += stepX {
			r, g, b, _ := color.NRGBAModel.Convert(img.At(x, y)).RGBA()
			rSum += float64(r >> 8)
			gSum += float64(g >> 8)
			bSum += float64(b >> 8)
			count += 1
		}
	}
	if count == 0 {
		a.mu.Lock()
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		a.mu.Unlock()
		return
	}
	a.mu.Lock()
	a.AverageColor = color.RGBA{
		R: uint8(clampFloat(rSum/count, 0, 255)),
		G: uint8(clampFloat(gSum/count, 0, 255)),
		B: uint8(clampFloat(bSum/count, 0, 255)),
		A: 255,
	}
	a.mu.Unlock()
}
