package rain

import (
	"math/rand"

	"virga-player/settings"
)

func (ps *ParticleSystem) Resize(width, height int) {
	ps.width = width
	ps.height = height
}

func (ps *ParticleSystem) GetParticles() []Particle {
	return ps.particles
}

func (ps *ParticleSystem) SetSpawnPaused(paused bool) {
	ps.spawnPaused = paused
}

func (ps *ParticleSystem) ApplyConfig(cfg *settings.Config) {
	ps.maxSize = cfg.MaxParticles
	ps.baseSpeed = float64(cfg.RainSpeed) / 100.0
	ps.lifeMul = float64(cfg.RainLifetime) / 100.0
	ps.speedMul = ps.baseSpeed
	ps.musicOn = cfg.MusicReactive || cfg.RainVisualizer
	ps.visualizer = cfg.RainVisualizer
	ps.intensity = float64(cfg.MusicReactiveIntensity) / 100.0
	ps.pulseBias = float64(cfg.RainPulse) / 100.0
	ps.pulseSpeed = float64(cfg.PulseSpeed) / 100.0
	ps.pulseEnabled = cfg.PulseOnRain()
	ps.separateFreq = cfg.SeparateFrequencies
	ps.enabled = cfg.RainEnabled
	newDir := directionFromConfig(cfg.Direction)
	if cfg.Direction != settings.DirectionRandom || ps.direction != newDir {
		ps.direction = newDir
	}
	if len(ps.particles) > ps.maxSize {
		ps.particles = ps.particles[:ps.maxSize]
	}
	if cfg.MaxParticles > 100 {
		ps.baseSpawn = 200.0
	} else {
		ps.baseSpawn = 150.0
	}
	ps.spawnRate = ps.baseSpawn
	ps.spawnMul = 1.0
	if !ps.enabled {
		ps.particles = nil
		return
	}
	for i := range ps.particles {
		if ps.direction == 0 {
			ps.particles[i].VelX = 0.0
		} else {
			ps.particles[i].VelX = float64(ps.direction) * (4.0 + rand.Float64()*2.0)
		}
	}
}

func (ps *ParticleSystem) ApplySpectrum(dt, low, mid, high, envelope float64) {
	applySpectrum(ps, dt, low, mid, high, envelope)
}

func (ps *ParticleSystem) ResetSpectrum() {
	resetSpectrum(ps)
	ps.bassPhase = 0
}


