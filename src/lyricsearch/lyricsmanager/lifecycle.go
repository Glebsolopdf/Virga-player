package lyricsmanager

import (
	"container/list"
	"fmt"
	"os"
)

func (m *LyricsManager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.setActiveTrackLocked("", nil)
	tempDir := m.cfg.TempDir
	mode := m.cfg.Mode
	cleanupCh := m.cleanupCh
	m.cleanupCh = nil

	m.ramIndex = map[string]*ramEntry{}
	m.ramOrder = list.New()
	m.ramBytes = 0
	m.mu.Unlock()

	if cleanupCh != nil {
		close(cleanupCh)
		m.cleanupWG.Wait()
	}

	if usesRAMStorage(m.cfg, mode) {
		if err := os.RemoveAll(tempDir); err != nil {
			return fmt.Errorf("remove temporary lyrics directory: %w", err)
		}
	}

	return nil
}
