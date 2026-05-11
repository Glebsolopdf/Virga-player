package artwork

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/gdamore/tcell/v2"
)

var (
	artworkLoadFailures   = make(map[string]struct{})
	artworkLoadFailuresMu sync.RWMutex
)

// DisplayMode
type DisplayMode int

const (
	DisplaySixel DisplayMode = iota
	DisplayTextOnly
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
	Duration         int
	Elapsed          int
	AnimationEnabled bool
	Fade             float64
	Pulse            float64
	Blink            float64
	LastEnvelope     float64
	ScrollTimer      float64
	ScrollOffset     int
	ScrollDir        int
	ScrollWidth      int
	SixelData        []byte
	mu               sync.RWMutex
	loadStarted      bool
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
		a.Mode = DisplayTextOnly
		a.beginCoverLoad()
	} else {
		a.Mode = DisplayTextOnly
	}

	return a
}

func loadFailed(imagePath string) bool {
	artworkLoadFailuresMu.RLock()
	defer artworkLoadFailuresMu.RUnlock()
	_, ok := artworkLoadFailures[imagePath]
	return ok
}

func markLoadFailed(imagePath string) {
	artworkLoadFailuresMu.Lock()
	artworkLoadFailures[imagePath] = struct{}{}
	artworkLoadFailuresMu.Unlock()
}

func (a *Artwork) beginCoverLoad() {
	if a.ImagePath == "" {
		return
	}

	if loadFailed(a.ImagePath) {
		return
	}

	a.mu.Lock()
	if a.loadStarted {
		a.mu.Unlock()
		return
	}
	a.loadStarted = true
	a.mu.Unlock()

	go func() {
		a.loadCoverImage()
		a.mu.RLock()
		img := a.CoverImg
		a.mu.RUnlock()
		if img == nil {
			markLoadFailed(a.ImagePath)
			return
		}
		a.computeAverageColor()
		a.mu.Lock()
		if DetectSixelSupport() && a.CoverImg != nil {
			a.Mode = DisplaySixel
		} else {
			a.Mode = DisplayTextOnly
		}
		a.mu.Unlock()
	}()
}

func (a *Artwork) getCoverImg() image.Image {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.CoverImg
}

func (a *Artwork) Render(screen tcell.Screen) {
	a.mu.RLock()
	mode := a.Mode
	a.mu.RUnlock()

	if a.AnimationEnabled && mode == DisplaySixel {
		a.renderTextOnly(screen)
		return
	}

	switch mode {
	case DisplaySixel:
		a.renderSixel(screen)
	case DisplayTextOnly:
		a.renderTextOnly(screen)
	}
}

func (a *Artwork) SetAnimationEnabled(enabled bool) {
	a.AnimationEnabled = enabled
	if !enabled {
		a.Fade = 1
		a.Pulse = 0
	}
}

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

	a.updateTextScroll(dt)
}

func (a *Artwork) updateTextScroll(dt float64) {
	if a.ScrollWidth <= 0 {
		a.ScrollOffset = 0
		return
	}

	titleLen := len([]rune(a.Title))
	artistLen := len([]rune(a.Artist))
	albumLen := len([]rune(a.Album))
	maxOffset := 0
	if titleLen > a.ScrollWidth {
		candidate := titleLen - a.ScrollWidth
		if candidate > maxOffset {
			maxOffset = candidate
		}
	}
	if artistLen > a.ScrollWidth {
		candidate := artistLen - a.ScrollWidth
		if candidate > maxOffset {
			maxOffset = candidate
		}
	}
	if albumLen > a.ScrollWidth {
		candidate := albumLen - a.ScrollWidth
		if candidate > maxOffset {
			maxOffset = candidate
		}
	}

	if maxOffset == 0 {
		a.ScrollOffset = 0
		return
	}

	if a.ScrollDir == 0 {
		a.ScrollDir = 1
	}

	a.ScrollTimer += dt
	const scrollInterval = 0.35
	if a.ScrollTimer < scrollInterval {
		return
	}
	a.ScrollTimer -= scrollInterval
	a.ScrollOffset += a.ScrollDir

	if a.ScrollOffset < 0 {
		a.ScrollOffset = 0
		a.ScrollDir = 1
	} else if a.ScrollOffset > maxOffset {
		a.ScrollOffset = maxOffset
		a.ScrollDir = -1
	}
}

// UpdateTrackInfo
func (a *Artwork) UpdateTrackInfo(title, artist, album string, duration, elapsed int) {
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
			r, g, b, _ := color.NRGBAModel.Convert(a.CoverImg.At(x, y)).RGBA()
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
