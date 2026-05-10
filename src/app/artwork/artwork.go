package artwork

import (
	"image"
	"image/color"
	"math"

	"github.com/gdamore/tcell/v2"
)

// DisplayMode определяет как отображать информацию о треке
type DisplayMode int

const (
	DisplaySixel    DisplayMode = iota // Реальная картинка через SIXEL
	DisplayTextOnly                    // Только текстовая информация
)

// Artwork
type Artwork struct {
	Mode             DisplayMode
	ImagePath        string
	CoverImg         image.Image
	AverageColor     color.Color
	Title            string
	Artist           string
	Album            string
	Duration         int // в секундах
	Elapsed          int // в секундах
	AnimationEnabled bool
	Fade             float64
	Pulse            float64
	Blink            float64
	LastEnvelope     float64
	SixelData        []byte
	SixelPath        string
}

// NewArtwork
func NewArtwork(imagePath, title, artist, album string, duration, elapsed int) *Artwork {
	a := &Artwork{
		ImagePath:    imagePath,
		Title:        title,
		Artist:       artist,
		Album:        album,
		Duration:     duration,
		Elapsed:      elapsed,
		AverageColor: color.RGBA{255, 255, 255, 255},
		Fade:         0,
		Pulse:        0,
		LastEnvelope: 0,
	}

	if imagePath != "" {
		a.loadCoverImage()
		a.computeAverageColor()
	}

	// Определяем поддержку терминала
	if DetectSixelSupport() && imagePath != "" {
		a.Mode = DisplaySixel
	} else {
		a.Mode = DisplayTextOnly
	}

	return a
}

// Render отображает обложку/информацию
func (a *Artwork) Render(screen tcell.Screen) {
	switch a.Mode {
	case DisplaySixel:
		a.renderSixel(screen)
	case DisplayTextOnly:
		a.renderTextOnly(screen)
	}
}

// SetAnimationEnabled toggles cover animation mode.
func (a *Artwork) SetAnimationEnabled(enabled bool) {
	a.AnimationEnabled = enabled
	if !enabled {
		a.Fade = 1
		a.Pulse = 0
	}
}

// UpdateAnimation обновляет внутреннее состояние эффектов обложки.
func (a *Artwork) UpdateAnimation(dt, envelope float64) {
	if !a.AnimationEnabled {
		a.Fade = 1
		a.Pulse = 0
		return
	}
	if dt < 0 {
		dt = 0
	}

	if a.Fade < 1 {
		a.Fade += dt * 1.1
		if a.Fade > 1 {
			a.Fade = 1
		}
	}
	envelope = clampFloat(envelope, 0, 1)
	burst := clampFloat((envelope-a.LastEnvelope)*3.0, 0, 1)
	peak := clampFloat(envelope*0.65+burst*0.9, 0, 1)
	targetPulse := math.Pow(peak, 0.65) * 0.96
	attack := 1 - math.Exp(-6.0*dt)
	release := 1 - math.Exp(-2.8*dt)

	if targetPulse > a.Pulse {
		if burst > 0.15 {
			a.Pulse = a.Pulse*0.2 + targetPulse*0.8
		} else {
			a.Pulse += (targetPulse - a.Pulse) * attack
		}
	} else {
		a.Pulse += (targetPulse - a.Pulse) * release
	}
	if a.Pulse > 0.99 {
		a.Pulse = 0.99
	}
	if a.Pulse < 0.005 {
		a.Pulse = 0
	}
	a.LastEnvelope = envelope
}

// UpdateTrackInfo обновляет информацию о треке
func (a *Artwork) UpdateTrackInfo(title, artist, album string, duration, elapsed int) {
	a.Title = title
	a.Artist = artist
	a.Album = album
	a.Duration = duration
	a.Elapsed = elapsed
}

func (a *Artwork) computeAverageColor() {
	if a.CoverImg == nil {
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		return
	}

	b := a.CoverImg.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		return
	}

	var rSum, gSum, bSum, count float64
	stepX := (b.Dx() / 64) + 1
	stepY := (b.Dy() / 64) + 1
	for y := b.Min.Y; y < b.Max.Y; y += stepY {
		for x := b.Min.X; x < b.Max.X; x += stepX {
			r, g, b, _ := color.NRGBAModel.Convert(a.CoverImg.At(x, y)).RGBA()
			rSum += float64(r >> 8)
			gSum += float64(g >> 8)
			bSum += float64(b >> 8)
			count += 1
		}
	}
	if count == 0 {
		a.AverageColor = color.RGBA{255, 255, 255, 255}
		return
	}
	a.AverageColor = color.RGBA{
		R: uint8(clampFloat(rSum/count, 0, 255)),
		G: uint8(clampFloat(gSum/count, 0, 255)),
		B: uint8(clampFloat(bSum/count, 0, 255)),
		A: 255,
	}
}
