package player

import (
	"fmt"
	"time"
	"virga-player/app/artwork"
)

type Player struct {
	Title            string
	Artist           string
	Album            string
	Duration         int
	Elapsed          int
	CoverX           int
	CoverY           int
	CoverW           int
	CoverH           int
	TextX            int
	TextY            int
	Hidden           map[int]bool
	HitTimes         map[int]time.Time
	RecoveryDuration time.Duration
	Artwork          *artwork.Artwork
	ArtworkPath      string
	ArtworkURL       string
}

type HitArea struct {
	X     int
	Y     int
	W     int
	H     int
	Field string
}

func New(centerX, centerY int) *Player {
	art := artwork.NewArtwork("", "Track Title", "Artist Name", "Album", 180, 0)

	return &Player{
		Title:            "Track Title",
		Artist:           "Artist Name",
		Album:            "Album",
		Duration:         180,
		Elapsed:          0,
		CoverX:           centerX - 20,
		CoverY:           centerY - 5,
		CoverW:           16,
		CoverH:           10,
		TextX:            centerX + 5,
		TextY:            centerY - 3,
		Hidden:           make(map[int]bool),
		HitTimes:         make(map[int]time.Time),
		RecoveryDuration: 500 * time.Millisecond,
		Artwork:          art,
		ArtworkPath:      "",
		ArtworkURL:       "",
	}
}

// SetTrackInfo
func (p *Player) SetTrackInfo(title, artist, album string, durationStr, elapsedStr string, durationSec, elapsedSec int) {
	p.Title = title
	p.Artist = artist
	p.Album = album
	p.Duration = durationSec
	p.Elapsed = elapsedSec

	// Обновляем Artwork
	if p.Artwork != nil {
		p.Artwork.UpdateTrackInfo(title, artist, album, durationSec, elapsedSec)
	}
}

// SetTrackInfoWithArtwork
func (p *Player) SetTrackInfoWithArtwork(title, artist, album string, durationStr, elapsedStr string, durationSec, elapsedSec int, artworkPath string) {
	trackChanged := title != p.Title || artist != p.Artist || album != p.Album
	p.SetTrackInfo(title, artist, album, durationStr, elapsedStr, durationSec, elapsedSec)
	if trackChanged {
		p.Elapsed = 0
	}

	if p.Artwork == nil {
		p.ArtworkPath = artworkPath
		p.Artwork = artwork.NewArtwork(artworkPath, title, artist, album, durationSec, elapsedSec)
		return
	}

	if artworkPath != "" && p.ArtworkPath != artworkPath {
		p.ArtworkPath = artworkPath
		p.Artwork = artwork.NewArtwork(artworkPath, title, artist, album, durationSec, elapsedSec)
		return
	}

	if trackChanged && artworkPath == "" {
		p.ArtworkPath = ""
		p.Artwork = artwork.NewArtwork("", title, artist, album, durationSec, elapsedSec)
		return
	}

	if p.Artwork != nil {
		p.Artwork.UpdateTrackInfo(title, artist, album, durationSec, p.Elapsed)
	}
}

func (p *Player) UpdatePosition(width, height int) {
	centerX := width / 2
	centerY := height / 2
	p.CoverX = centerX - 20
	p.CoverY = centerY - 5
	p.TextX = centerX + 5
	p.TextY = centerY - 3
}

func (p *Player) GetHitAreas() []HitArea {
	timeStr := fmt.Sprintf("%d:%02d / %d:%02d", p.Elapsed/60, p.Elapsed%60, p.Duration/60, p.Duration%60)

	areas := []HitArea{
		{
			X:     p.CoverX,
			Y:     p.CoverY,
			W:     p.CoverW,
			H:     p.CoverH,
			Field: "cover",
		},
		{
			X:     p.TextX,
			Y:     p.TextY,
			W:     len(p.Title),
			H:     1,
			Field: "title",
		},
		{
			X:     p.TextX,
			Y:     p.TextY + 1,
			W:     len(p.Artist),
			H:     1,
			Field: "artist",
		},
		{
			X:     p.TextX,
			Y:     p.TextY + 2,
			W:     len(p.Album),
			H:     1,
			Field: "album",
		},
		{
			X:     p.TextX,
			Y:     p.TextY + 3,
			W:     len(timeStr),
			H:     1,
			Field: "duration",
		},
	}
	return areas
}

func (p *Player) MarkHit(field string, x, y int) {
	key := y*1000 + x
	p.Hidden[key] = true
	p.HitTimes[key] = time.Now()
}

func (p *Player) IsHidden(field string, x, y int) bool {
	key := y*1000 + x
	if !p.Hidden[key] {
		return false
	}

	if hitTime, ok := p.HitTimes[key]; ok {
		if time.Since(hitTime) > p.RecoveryDuration {
			delete(p.Hidden, key)
			delete(p.HitTimes, key)
			return false
		}
	}

	return true
}

func (p *Player) Restore() {
	p.Hidden = make(map[int]bool)
	p.HitTimes = make(map[int]time.Time)
}
