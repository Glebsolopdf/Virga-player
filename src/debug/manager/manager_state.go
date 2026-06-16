package manager

import "time"

func (m *Manager) Enabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

func (m *Manager) SetEnabled(enabled bool) {
	m.mu.Lock()
	if m.forced {
		enabled = true
	}
	prev := m.enabled
	m.enabled = enabled
	m.mu.Unlock()
	if prev != enabled {
		if enabled {
			m.log(levelInfo, "debug mode enabled")
		} else {
			m.log(levelInfo, "debug mode disabled")
		}
	}
}

func (m *Manager) setStatus(status string) {
	m.mu.Lock()
	m.status = status
	m.statusAt = time.Now()
	m.mu.Unlock()
	m.log(levelInfo, status)
}
