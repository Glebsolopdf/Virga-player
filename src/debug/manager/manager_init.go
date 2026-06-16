package manager

import (
	"io"

	"virga-player/debug/storage"
)

func NewManager(enabled, forced bool) *Manager {
	m := &Manager{
		enabled: enabled || forced,
		forced:  forced,
		buf:     storage.NewRingBuffer(1000),
	}
	m.log(levelInfo, "debug manager initialized")
	if forced {
		m.log(levelInfo, "debug mode enabled by --debug flag")
	}
	return m
}

func (m *Manager) Writer() io.Writer {
	return &lineWriter{mgr: m}
}
