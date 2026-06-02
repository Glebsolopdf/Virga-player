package app

import (
	"strings"
	"time"

	"virga-player/lyricsearch"
)

func (a *App) resetLyricsManager() {
	a.closeLyricsManager()

	if a.cfg == nil {
		a.lyricsDoubleConfirm.Store(true)
		return
	}

	a.lyricsDoubleConfirm.Store(a.cfg.LyricsDoubleConfirm)

	var debugf func(format string, args ...any)
	if a.debug != nil {
		debugf = a.debug.Debugf
	}

	mgr, err := lyricsearch.NewLyricsManager(lyricsearch.Config{
		Mode:          string(a.cfg.LyricsMode),
		TempDir:       strings.TrimSpace(a.cfg.LyricsTempDir),
		PersistentDir: strings.TrimSpace(a.cfg.LyricsPersistentDir),
		AutoSaveAfter: time.Duration(a.cfg.LyricsAutoSaveAfterSec) * time.Second,
		Prompt:        a.onLyricsSavePrompt,
		Debugf:        debugf,
	})
	if err != nil {
		if a.debug != nil {
			a.debug.Warnf("lyrics manager init failed: %v", err)
		}
		return
	}

	a.lyricsManager = mgr
}

func (a *App) closeLyricsManager() {
	if a.lyricsManager == nil {
		return
	}
	if err := a.lyricsManager.Close(); err != nil && a.debug != nil {
		a.debug.Warnf("lyrics manager close failed: %v", err)
	}
	a.lyricsManager = nil
}
