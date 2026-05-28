package player

import "virga-player/app/artwork"

type Player struct {
	Title       string
	Artist      string
	Album       string
	Duration    int
	Elapsed     int
	CoverX      int
	CoverY      int
	CoverW      int
	CoverH      int
	TextX       int
	TextY       int
	Artwork     *artwork.Artwork
	ArtworkPath string
	ArtworkURL  string
}

func New(centerX, centerY int) *Player {
	art := artwork.NewArtwork("", "Track Title", "Artist Name", "Album", 180, 0)

	return &Player{
		Title:       "Track Title",
		Artist:      "Artist Name",
		Album:       "Album",
		Duration:    180,
		Elapsed:     0,
		CoverX:      centerX - 20,
		CoverY:      centerY - 5,
		CoverW:      16,
		CoverH:      10,
		TextX:       centerX + 5,
		TextY:       centerY - 3,
		Artwork:     art,
		ArtworkPath: "",
		ArtworkURL:  "",
	}
}

// SetTrackInfo
func (p *Player) SetTrackInfo(title, artist, album string, durationSec, elapsedSec int) {
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
func (p *Player) SetTrackInfoWithArtwork(title, artist, album string, durationSec, elapsedSec int, artworkPath string) {
	trackChanged := title != p.Title || artist != p.Artist || album != p.Album
	p.SetTrackInfo(title, artist, album, durationSec, elapsedSec)
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
