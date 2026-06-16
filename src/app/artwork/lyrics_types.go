package artwork

import "github.com/gdamore/tcell/v2"

const (
	maxVisibleLyrics  = 3
	lyricsBasePadding = 3
	lyricsDefaultGap  = 2
	lyricsScrollStep  = 220
)

type lyricCue struct {
	atMillis int
	text     string
}

type lyricDisplayLine struct {
	text    string
	color   tcell.Color
	active  bool
	hasAnim bool
	anim    float64
}
