package manager

import (
	"bytes"
	"strings"
)

type lineWriter struct {
	mgr *Manager
}

func (w *lineWriter) Write(p []byte) (int, error) {
	if w == nil || w.mgr == nil {
		return len(p), nil
	}
	trimmed := bytes.TrimSpace(p)
	if len(trimmed) == 0 {
		return len(p), nil
	}
	for _, part := range strings.Split(string(trimmed), "\n") {
		line := strings.TrimSpace(part)
		if line == "" {
			continue
		}
		w.mgr.log(levelInfo, line)
	}
	return len(p), nil
}
