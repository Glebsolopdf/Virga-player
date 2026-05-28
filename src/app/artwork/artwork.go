package artwork

import (
	"image"
	"image/color"
	"sync"

	textrender "virga-player/app/artwork/textrender"
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
	textRenderCache  textrender.Cache
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
	RainTimer        float64
	RainOffsetX      float64
	RainOffsetY      float64
	RainResistance   float64
	SixelData        []byte
	pulseTarget      float64
	lastPulseKey     float64
	beatTimer        float64
	beatInterval     float64
	adaptiveSpeed    float64
	pulseActive      bool
	pulseAttacking   bool
	mu               sync.RWMutex
	loadStarted      bool
	sixelBuilding    bool
}

type artworkSnapshot struct {
	mode             DisplayMode
	imagePath        string
	coverImg         image.Image
	title            string
	artist           string
	album            string
	duration         int
	elapsed          int
	animationEnabled bool
	fade             float64
	pulse            float64
	rainOffsetX      float64
	rainOffsetY      float64
	sixelData        []byte
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
