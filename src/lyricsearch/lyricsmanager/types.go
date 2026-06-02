package lyricsmanager

import (
	"container/list"
	"context"
	"errors"
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
