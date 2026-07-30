package lyricsearch

import "errors"

const (
	ModeDisabled      = "disabled"
	ModeRAMOnly       = "ram-only"
	ModeRAMWithAuto   = "ram-with-auto-save"
	ModeRAMWithPrompt = "ram-with-save-prompt"
	ModeDirectToDisk  = "direct-to-disk"

	// Legacy modes kept for backward compatibility with older configs.
	ModeLocal = "local"
	ModeAuto  = "auto"
)

var (
	ErrLyricsDisabled    = errors.New("lyrics mode is disabled")
	ErrLyricsNotFound    = errors.New("synced lyrics not found")
	ErrInstrumentalTrack = errors.New("track is instrumental")
	ErrMissingMetadata   = errors.New("lyrics search requires artist and track")
)
