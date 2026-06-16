package manager

import (
	"runtime"
	"time"
)

func (m *Manager) UpdateRuntime(dt float64, particles, particlesMax, targetFPS int) {
	if !m.Enabled() {
		return
	}
	now := time.Now()
	instFPS := 0.0
	if dt > 0 {
		instFPS = 1.0 / dt
	}
	needSample, prevSample, prevCPU := m.setRuntimeFrame(now, instFPS, particles, particlesMax, targetFPS)
	if !needSample {
		return
	}
	procCPU := processCPUSeconds()
	cpuPct := runtimeCPUPercent(now, prevSample, procCPU, prevCPU)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.mu.Lock()
	m.lastCPUTime = procCPU
	m.cpuPercent = cpuPct
	m.memMiB = float64(mem.Alloc) / (1024.0 * 1024.0)
	m.goroutines = runtime.NumGoroutine()
	m.mu.Unlock()
}

func (m *Manager) setRuntimeFrame(now time.Time, instFPS float64, particles, particlesMax, targetFPS int) (bool, time.Time, float64) {
	m.mu.Lock()
	if instFPS > 0 {
		if m.fps == 0 {
			m.fps = instFPS
		} else {
			m.fps = m.fps*0.86 + instFPS*0.14
		}
	}
	m.particles = particles
	m.particlesMax = particlesMax
	m.targetFPS = targetFPS
	needSample := m.lastCPUSample.IsZero() || now.Sub(m.lastCPUSample) >= 500*time.Millisecond
	prevSample := m.lastCPUSample
	prevCPU := m.lastCPUTime
	if needSample {
		m.lastCPUSample = now
	}
	m.mu.Unlock()
	return needSample, prevSample, prevCPU
}

func runtimeCPUPercent(now, prevSample time.Time, procCPU, prevCPU float64) float64 {
	if prevSample.IsZero() {
		return 0
	}
	wall := now.Sub(prevSample).Seconds()
	deltaCPU := procCPU - prevCPU
	if wall <= 0 || deltaCPU < 0 {
		return 0
	}
	return (deltaCPU / wall) * 100.0
}
