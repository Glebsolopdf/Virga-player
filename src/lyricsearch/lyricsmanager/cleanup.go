package lyricsmanager

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (m *LyricsManager) enqueueCleanup(task cleanupTask) {
	if len(task.files) == 0 && len(task.artistDirs) == 0 {
		return
	}

	m.mu.RLock()
	ch := m.cleanupCh
	closed := m.closed
	m.mu.RUnlock()
	if closed || ch == nil {
		m.runCleanup(task)
		return
	}

	select {
	case ch <- task:
	default:
		go m.runCleanup(task)
	}
}

func (m *LyricsManager) cleanupWorker() {
	defer m.cleanupWG.Done()
	for task := range m.cleanupCh {
		m.runCleanup(task)
	}
}

func (m *LyricsManager) runCleanup(task cleanupTask) {
	for _, filePath := range task.files {
		if strings.TrimSpace(filePath) == "" {
			continue
		}
		if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			m.debugf("temp cleanup remove failed path=%q err=%v", filePath, err)
			continue
		}
		m.debugf("temp cleanup removed file=%q", filePath)
	}

	for _, dirPath := range task.artistDirs {
		m.pruneArtistDirIfEmpty(dirPath)
	}
}

func (m *LyricsManager) pruneArtistDirIfEmpty(dirPath string) {
	if strings.TrimSpace(dirPath) == "" {
		return
	}

	m.mu.RLock()
	tempRoot := m.cfg.TempDir
	m.mu.RUnlock()

	cleanDir := filepath.Clean(dirPath)
	cleanRoot := filepath.Clean(tempRoot)
	if cleanDir == cleanRoot {
		return
	}
	if !strings.HasPrefix(cleanDir, cleanRoot+string(os.PathSeparator)) {
		return
	}

	entries, err := os.ReadDir(cleanDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			m.debugf("temp cleanup readdir failed dir=%q err=%v", cleanDir, err)
		}
		return
	}
	if len(entries) != 0 {
		return
	}

	if err := os.Remove(cleanDir); err != nil && !errors.Is(err, os.ErrNotExist) {
		m.debugf("temp cleanup remove artist dir failed dir=%q err=%v", cleanDir, err)
		return
	}
	m.debugf("temp cleanup removed empty artist dir=%q", cleanDir)
}

func (m *LyricsManager) debugf(format string, args ...any) {
	if m.cfg.Debugf != nil {
		m.cfg.Debugf("lyrics manager: "+format, args...)
		return
	}
	log.Printf("lyricsearch: "+format, args...)
}
