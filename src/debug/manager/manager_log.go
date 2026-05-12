package manager

import (
	"fmt"
	"strings"
	"time"

	"virga-player/debug/storage"
)

func (m *Manager) Infof(format string, args ...any) {
	m.log(levelInfo, fmt.Sprintf(format, args...))
}

func (m *Manager) Warnf(format string, args ...any) {
	m.log(levelWarn, fmt.Sprintf(format, args...))
}

func (m *Manager) Errorf(format string, args ...any) {
	m.log(levelError, fmt.Sprintf(format, args...))
}

func (m *Manager) Debugf(format string, args ...any) {
	m.log(levelDebug, fmt.Sprintf(format, args...))
}

func (m *Manager) log(lvl level, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	m.buf.Append(storage.Entry{At: time.Now(), Level: string(lvl), Text: text})
}
