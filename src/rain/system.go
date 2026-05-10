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
	ps.musicOn = cfg.MusicReactive
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

	// Silence: envelope below noise floor → stop rain completely.
	const silenceThreshold = 0.06
	if envelope < silenceThreshold {
		ps.silenced = true
		ps.spawnMul = 0
		// speedMul stays — existing drops continue falling.
		return
	}
	ps.silenced = false

	// Direct per-band energy mapping — no state machine.
	// Low+mid drive speed (bass and body = flow intensity).
	// High+mid drive spawn rate (cymbals/hi-hats = drops appearing).
	speedEnergy := clamp(low*0.55+mid*0.35+high*0.10, 0, 1)
	spawnEnergy := clamp(high*0.50+mid*0.35+low*0.15, 0, 1)

	// energyMul: 0.05 (near-silent) → 4.5 (full drop). Intensity scales the ceiling.
	ps.energyMul = clamp(speedEnergy*4.5*ps.intensity, 0.05, 4.5)
	// spawnMul: scales spawn rate; highs spike it during percussion hits.
	ps.spawnMul = clamp(spawnEnergy*6.0*ps.intensity, 0.05, 6.0)
}

func (ps *ParticleSystem) ResetSpectrum() {
	ps.speedMul = ps.baseSpeed
	ps.spawnMul = 1.0
	ps.bassPhase = 0
	ps.silenced = false
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
