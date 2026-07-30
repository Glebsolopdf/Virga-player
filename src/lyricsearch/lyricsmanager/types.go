package lyricsmanager

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	defaultRAMMaxFiles         = 10
	defaultRAMMaxBytes   int64 = 15 * 1024 * 1024
	defaultAutoSaveAfter       = 30 * time.Second
)

var ErrLyricsManagerClosed = errors.New("lyrics manager is closed")

type PromptCallback func(context.Context, PromptRequest) bool

type PromptRequest struct {
	Track   Track
	Lyrics  string
	Message string
}

type Track struct {
	Artist string
	Title  string
}

type Config struct {
	Mode          string
	TempDir       string
	PersistentDir string

	RAMMaxFiles int
	RAMMaxBytes int64

	AutoSaveAfter time.Duration
	Prompt        PromptCallback
	Debugf        func(format string, args ...any)

	NormalizeMode     func(string) string
	ReadLyricsFromDir func(baseDir, artist, track string) (string, error)
	WriteLyricsToDir  func(baseDir, artist, track, lyrics string) (string, error)
	FetchLyrics       func(artist, track string) (string, error)
	HasTimedLyrics    func(lyrics string) bool

	ModeDisabled      string
	ModeLocal         string
	ModeAuto          string
	ModeRAMOnly       string
	ModeRAMWithAuto   string
	ModeRAMWithPrompt string
	ModeDirectToDisk  string

	ErrLyricsDisabled  error
	ErrMissingMetadata error
}

type ramEntry struct {
	key     string
	artist  string
	title   string
	lyrics  string
	path    string
	size    int64
	element *list.Element
}

type cleanupTask struct {
	files      []string
	artistDirs []string
}

type LyricsManager struct {
	mu sync.RWMutex

	cfg Config

	ramIndex map[string]*ramEntry
	ramOrder *list.List
	ramBytes int64

	activeTrackKey    string
	activeTrackCancel context.CancelFunc

	closed bool

	cleanupCh chan cleanupTask
	cleanupWG sync.WaitGroup
}

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
