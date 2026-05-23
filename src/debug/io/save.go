package io

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"virga-player/settings"
)

func SaveDumpToDisk(text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("no log data to save")
	}
	dir := filepath.Join(settings.ConfigDirPath(), "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("virga-debug-%s.log", time.Now().Format("20060102-150405"))
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
