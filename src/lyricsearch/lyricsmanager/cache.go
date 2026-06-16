package lyricsmanager

import (
	"context"
	"fmt"
	"path/filepath"
)

func (m *LyricsManager) getRAM(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.ramIndex[key]
	if !ok {
		return "", false
	}

	m.ramOrder.MoveToFront(entry.element)
	return entry.lyrics, true
}

func (m *LyricsManager) putRAM(key, artist, title, lyrics, tempDir string) error {
	if m.cfg.WriteLyricsToDir == nil {
		return fmt.Errorf("write lyrics callback is not configured")
	}
	path, err := m.cfg.WriteLyricsToDir(tempDir, artist, title, lyrics)
	if err != nil {
		return err
	}

	size := int64(len([]byte(lyrics)))

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return ErrLyricsManagerClosed
	}

	cleanup := cleanupTask{}

	if existing, ok := m.ramIndex[key]; ok {
		m.ramBytes -= existing.size
		existing.artist = artist
		existing.title = title
		existing.lyrics = lyrics
		existing.path = path
		existing.size = size
		m.ramBytes += size
		m.ramOrder.MoveToFront(existing.element)
		cleanup = m.evictIfNeededLocked()
		entries := len(m.ramIndex)
		bytes := m.ramBytes
		m.mu.Unlock()
		m.enqueueCleanup(cleanup)
		m.debugf("ram cache update artist=%q track=%q entries=%d bytes=%d", artist, title, entries, bytes)
		return nil
	}

	entry := &ramEntry{
		key:    key,
		artist: artist,
		title:  title,
		lyrics: lyrics,
		path:   path,
		size:   size,
	}
	entry.element = m.ramOrder.PushFront(entry)
	m.ramIndex[key] = entry
	m.ramBytes += size
	cleanup = m.evictIfNeededLocked()
	entries := len(m.ramIndex)
	bytes := m.ramBytes
	m.mu.Unlock()
	m.enqueueCleanup(cleanup)
	m.debugf("ram cache insert artist=%q track=%q entries=%d bytes=%d", artist, title, entries, bytes)
	return nil
}

func (m *LyricsManager) evictIfNeededLocked() cleanupTask {
	task := cleanupTask{}
	for len(m.ramIndex) > m.cfg.RAMMaxFiles || m.ramBytes > m.cfg.RAMMaxBytes {
		tail := m.ramOrder.Back()
		if tail == nil {
			return task
		}

		entry, ok := tail.Value.(*ramEntry)
		if !ok || entry == nil {
			m.ramOrder.Remove(tail)
			continue
		}

		delete(m.ramIndex, entry.key)
		m.ramOrder.Remove(tail)
		m.ramBytes -= entry.size
		if entry.path != "" {
			task.files = append(task.files, entry.path)
			task.artistDirs = append(task.artistDirs, filepath.Dir(entry.path))
		}
		m.debugf("ram cache evict artist=%q track=%q entries=%d bytes=%d", entry.artist, entry.title, len(m.ramIndex), m.ramBytes)
	}

	return task
}

func (m *LyricsManager) setActiveTrack(key string, cancel context.CancelFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.setActiveTrackLocked(key, cancel)
}

func (m *LyricsManager) setActiveTrackLocked(key string, cancel context.CancelFunc) {
	if m.activeTrackCancel != nil {
		m.activeTrackCancel()
	}
	m.activeTrackKey = key
	m.activeTrackCancel = cancel
}

func (m *LyricsManager) cancelActiveTrack() {
	m.setActiveTrack("", nil)
}
