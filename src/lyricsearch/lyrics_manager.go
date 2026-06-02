package lyricsearch

import (
	"time"

	"virga-player/lyricsearch/lyricsmanager"
)

type PromptCallback = lyricsmanager.PromptCallback

type PromptRequest = lyricsmanager.PromptRequest

type Track = lyricsmanager.Track

type LyricsManager = lyricsmanager.LyricsManager

var ErrLyricsManagerClosed = lyricsmanager.ErrLyricsManagerClosed

type Config struct {
	Mode          string
	TempDir       string
	PersistentDir string

	RAMMaxFiles int
	RAMMaxBytes int64

	AutoSaveAfter time.Duration
	Prompt        PromptCallback
	Debugf        func(format string, args ...any)
}

func NewLyricsManager(cfg Config) (*LyricsManager, error) {
	managerCfg := lyricsmanager.Config{
		Mode:               cfg.Mode,
		TempDir:            cfg.TempDir,
		PersistentDir:      cfg.PersistentDir,
		RAMMaxFiles:        cfg.RAMMaxFiles,
		RAMMaxBytes:        cfg.RAMMaxBytes,
		AutoSaveAfter:      cfg.AutoSaveAfter,
		Prompt:             cfg.Prompt,
		Debugf:             cfg.Debugf,
		NormalizeMode:      normalizeMode,
		ReadLyricsFromDir:  readLyricsFromDir,
		WriteLyricsToDir:   writeLyricsToDir,
		FetchLyrics:        fetchSyncedLyrics,
		HasTimedLyrics:     hasTimedLyrics,
		ModeDisabled:       ModeDisabled,
		ModeLocal:          ModeLocal,
		ModeAuto:           ModeAuto,
		ModeRAMOnly:        ModeRAMOnly,
		ModeRAMWithAuto:    ModeRAMWithAuto,
		ModeRAMWithPrompt:  ModeRAMWithPrompt,
		ModeDirectToDisk:   ModeDirectToDisk,
		ErrLyricsDisabled:  ErrLyricsDisabled,
		ErrMissingMetadata: ErrMissingMetadata,
	}

	return lyricsmanager.NewLyricsManager(managerCfg)
}
