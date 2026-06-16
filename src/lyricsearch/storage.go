package lyricsearch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readCachedLyrics(artist, track string) (string, error) {
	path, err := lyricsPath(artist, track)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeCachedLyrics(artist, track, lyrics string) error {
	path, err := lyricsPath(artist, track)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lyrics cache directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(lyrics), 0o644); err != nil {
		return fmt.Errorf("write lyrics cache file: %w", err)
	}

	return nil
}

func lyricsPath(artist, track string) (string, error) {
	return lyricsPathIn(defaultPersistentLyricsDir(), artist, track)
}

func defaultPersistentLyricsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/tmp", "virga-player", "lyrics")
	}

	return filepath.Join(
		homeDir,
		".config",
		"virga-player",
		"lyrics",
	)
}

func DefaultPersistentDir() string {
	return defaultPersistentLyricsDir()
}

func defaultTempLyricsDir() string {
	return filepath.Join(os.TempDir(), "virgaplayerlyrics")
}

func lyricsPathIn(baseDir, artist, track string) (string, error) {
	if strings.TrimSpace(baseDir) == "" {
		return "", fmt.Errorf("lyrics directory is empty")
	}

	return filepath.Join(
		baseDir,
		sanitizePathSegment(artist),
		sanitizePathSegment(track)+".lrc",
	), nil
}

func readLyricsFromDir(baseDir, artist, track string) (string, error) {
	path, err := lyricsPathIn(baseDir, artist, track)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeLyricsToDir(baseDir, artist, track, lyrics string) (string, error) {
	path, err := lyricsPathIn(baseDir, artist, track)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create lyrics directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(lyrics), 0o644); err != nil {
		return "", fmt.Errorf("write lyrics file: %w", err)
	}

	return path, nil
}

func sanitizePathSegment(value string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		"?", "_",
		"*", "_",
		":", "_",
		"|", "_",
		`"`, "_",
		"<", "_",
		">", "_",
	)

	sanitized := strings.TrimSpace(replacer.Replace(value))
	sanitized = strings.Trim(sanitized, ".")
	if sanitized == "" || sanitized == "." || sanitized == ".." {
		return "unknown"
	}

	return sanitized
}
