package lyricsmanager

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewLyricsManager(cfg Config) (*LyricsManager, error) {
	normalized := normalizeManagerConfig(cfg)

	if normalized.Mode != normalized.ModeDisabled {
		if err := os.MkdirAll(normalized.PersistentDir, 0o755); err != nil {
			return nil, fmt.Errorf("create persistent lyrics directory: %w", err)
		}
	}

	if usesRAMStorage(normalized, normalized.Mode) {
		if err := os.MkdirAll(normalized.TempDir, 0o755); err != nil {
			return nil, fmt.Errorf("create temporary lyrics directory: %w", err)
		}
	}

	mgr := &LyricsManager{
		cfg:      normalized,
		ramIndex: make(map[string]*ramEntry, normalized.RAMMaxFiles),
		ramOrder: list.New(),
	}

	if usesRAMStorage(normalized, normalized.Mode) {
		mgr.cleanupCh = make(chan cleanupTask, 64)
		mgr.cleanupWG.Add(1)
		go mgr.cleanupWorker()
	}

	return mgr, nil
}

func normalizeManagerConfig(cfg Config) Config {
	if cfg.ModeDisabled == "" {
		cfg.ModeDisabled = "disabled"
	}
	if cfg.ModeLocal == "" {
		cfg.ModeLocal = "local"
	}
	if cfg.ModeAuto == "" {
		cfg.ModeAuto = "auto"
	}
	if cfg.ModeRAMOnly == "" {
		cfg.ModeRAMOnly = "ram-only"
	}
	if cfg.ModeRAMWithAuto == "" {
		cfg.ModeRAMWithAuto = "ram-with-auto-save"
	}
	if cfg.ModeRAMWithPrompt == "" {
		cfg.ModeRAMWithPrompt = "ram-with-save-prompt"
	}
	if cfg.ModeDirectToDisk == "" {
		cfg.ModeDirectToDisk = "direct-to-disk"
	}
	if cfg.ErrLyricsDisabled == nil {
		cfg.ErrLyricsDisabled = errors.New("lyrics mode is disabled")
	}
	if cfg.ErrMissingMetadata == nil {
		cfg.ErrMissingMetadata = errors.New("lyrics search requires artist and track")
	}

	if cfg.NormalizeMode != nil {
		cfg.Mode = cfg.NormalizeMode(cfg.Mode)
	}
	if cfg.Mode == cfg.ModeLocal {
		cfg.Mode = cfg.ModeRAMOnly
	}
	if cfg.Mode == cfg.ModeAuto {
		cfg.Mode = cfg.ModeRAMWithAuto
	}

	if strings.TrimSpace(cfg.TempDir) == "" {
		cfg.TempDir = filepath.Join(os.TempDir(), "virgaplayerlyrics")
	}
	if strings.TrimSpace(cfg.PersistentDir) == "" {
		cfg.PersistentDir = filepath.Join("/tmp", "virga-player", "lyrics")
	}
	if cfg.RAMMaxFiles <= 0 {
		cfg.RAMMaxFiles = defaultRAMMaxFiles
	}
	if cfg.RAMMaxBytes <= 0 {
		cfg.RAMMaxBytes = defaultRAMMaxBytes
	}
	if cfg.AutoSaveAfter <= 0 {
		cfg.AutoSaveAfter = defaultAutoSaveAfter
	}

	return cfg
}

func usesRAMStorage(cfg Config, mode string) bool {
	switch mode {
	case cfg.ModeRAMOnly, cfg.ModeRAMWithAuto, cfg.ModeRAMWithPrompt:
		return true
	default:
		return false
	}
}
