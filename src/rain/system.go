package rain

import (
	"math/rand"

	rainspectrum "virga-player/rain/spectrum"
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
	state := ps.spectrumState()
	rainspectrum.Apply(&state, dt, low, mid, high, envelope, ps.musicOn)
	ps.applySpectrumState(state)
}

func (ps *ParticleSystem) ResetSpectrum() {
	state := ps.spectrumState()
	rainspectrum.Reset(&state)
	ps.applySpectrumState(state)
	ps.bassPhase = 0
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

func (ps *ParticleSystem) spectrumState() rainspectrum.State {
	return rainspectrum.State{
		BaseSpeed:     ps.baseSpeed,
		BaseSpawn:     ps.baseSpawn,
		Intensity:     ps.intensity,
		PulseBias:     ps.pulseBias,
		PulseSpeed:    ps.pulseSpeed,
		SpeedMul:      ps.speedMul,
		SpawnMul:      ps.spawnMul,
		SpawnRate:     ps.spawnRate,
		EnergyMul:     ps.energyMul,
		LowEnergy:     ps.lowEnergy,
		MidEnergy:     ps.midEnergy,
		HighEnergy:    ps.highEnergy,
		LastEnvelope:  ps.lastEnvelope,
		Pulse:         ps.pulse,
		PulseTarget:   ps.pulseTarget,
		LastPulseKey:  ps.lastPulseKey,
		BeatTimer:     ps.beatTimer,
		BeatInterval:  ps.beatInterval,
		AdaptiveSpeed: ps.adaptiveSpeed,
		Silenced:      ps.silenced,
		PulseActive:   ps.pulseActive,
		PulseAttack:   ps.pulseAttack,
		PulseEnabled:  ps.pulseEnabled,
		Visualizer:    ps.visualizer,
	}
}

func (ps *ParticleSystem) applySpectrumState(state rainspectrum.State) {
	ps.speedMul = state.SpeedMul
	ps.spawnMul = state.SpawnMul
	ps.spawnRate = state.SpawnRate
	ps.energyMul = state.EnergyMul
	ps.lowEnergy = state.LowEnergy
	ps.midEnergy = state.MidEnergy
	ps.highEnergy = state.HighEnergy
	ps.lastEnvelope = state.LastEnvelope
	ps.pulse = state.Pulse
	ps.pulseTarget = state.PulseTarget
	ps.lastPulseKey = state.LastPulseKey
	ps.beatTimer = state.BeatTimer
	ps.beatInterval = state.BeatInterval
	ps.adaptiveSpeed = state.AdaptiveSpeed
	ps.silenced = state.Silenced
	ps.pulseActive = state.PulseActive
	ps.pulseAttack = state.PulseAttack
}
