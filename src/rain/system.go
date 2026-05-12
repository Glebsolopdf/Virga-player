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
	ps.speedMul = ps.baseSpeed
	ps.musicOn = cfg.MusicReactive || cfg.RainVisualizer
	ps.visualizer = cfg.RainVisualizer
	ps.intensity = float64(cfg.MusicReactiveIntensity) / 100.0
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

func (ps *ParticleSystem) ApplySpectrum(low, mid, high, envelope float64) {
	if !ps.musicOn {
		ps.speedMul = ps.baseSpeed
		ps.spawnMul = 1.0
		ps.silenced = false
		return
	}

	const silenceThreshold = 0.06
	if envelope < silenceThreshold {
		ps.silenced = true
		ps.spawnMul = 0
		return
	}
	ps.silenced = false

	ps.lowEnergy = low
	ps.midEnergy = mid
	ps.highEnergy = high

	if ps.visualizer {
		ps.speedMul = ps.baseSpeed * clamp(1.0+envelope*2.0*ps.intensity, 0.5, 3.0)
		ps.spawnMul = 0.4
		ps.energyMul = clamp(low*0.6+mid*0.7+high*0.9, 0, 1)
		return
	}

	// Use envelope more actively so rain reacts faster to beat peaks.
	speedEnergy := clamp(low*0.30+mid*0.40+high*0.15+envelope*0.25, 0, 1)
	spawnEnergy := clamp(high*0.55+mid*0.30+low*0.05+envelope*0.40, 0, 1)
	ps.energyMul = clamp(speedEnergy*6.2*ps.intensity, 0.12, 7.0)
	ps.spawnMul = clamp(spawnEnergy*8.0*ps.intensity, 0.12, 9.0)
}

func (ps *ParticleSystem) ResetSpectrum() {
	ps.speedMul = ps.baseSpeed
	ps.spawnMul = 1.0
	ps.bassPhase = 0
	ps.lastEnvelope = 0
	ps.silenced = false
}

func (ps *ParticleSystem) RainForce() float64 {
	if ps.maxSize <= 0 {
		return 0
	}
	force := float64(len(ps.particles)) / float64(ps.maxSize)
	if force > 1 {
		force = 1
	}
	return force
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
