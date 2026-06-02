package manager

import "fmt"

func (m *Manager) runtimeSummary() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf(
		"FPS %.1f/%d | Particles %d/%d | CPU %.1f%% | RAM %.1f MiB | G %d",
		m.fps,
		m.targetFPS,
		m.particles,
		m.particlesMax,
		m.cpuPercent,
		m.memMiB,
		m.goroutines,
	)
}
