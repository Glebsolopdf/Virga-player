package rain

import (
	"math"
	"math/rand"

	"virga-player/settings"
)

func (ps *ParticleSystem) resetPulseState() {
	ps.pulse = 0
	ps.pulseTarget = 0
	ps.pulseActive = false
	ps.pulseAttack = false
}

func (ps *ParticleSystem) adaptivePulseRate(dt, pulseKey, burst float64) float64 {
	ps.beatTimer += dt
	base := clamp(ps.pulseSpeed, 0.25, 3.0)
	if ps.adaptiveSpeed == 0 {
		ps.adaptiveSpeed = base
	}

	if burst > 0.05 && pulseKey > 0.48 {
		if ps.beatTimer >= 0.18 && ps.beatTimer <= 1.20 {
			if ps.beatInterval == 0 {
				ps.beatInterval = ps.beatTimer
			} else {
				ps.beatInterval += (ps.beatTimer - ps.beatInterval) * 0.34
			}
		}
		ps.beatTimer = 0
	}

	target := base
	if ps.beatInterval > 0 {
		tempoFactor := clamp(0.46/ps.beatInterval, 0.72, 1.75)
		intensityFactor := 0.88 + clamp(pulseKey*0.46+burst*0.90, 0, 1)*0.44
		target = clamp(base*tempoFactor*intensityFactor, 0.25, 3.0)
	}
	smooth := 1 - math.Exp(-6.0*dt)
	ps.adaptiveSpeed += (target - ps.adaptiveSpeed) * smooth
	return clamp(ps.adaptiveSpeed, 0.25, 3.0)
}

func (ps *ParticleSystem) updatePulse(dt, low, mid, high, envelope float64) {
	pulseKey := math.Max(
		clamp(high*0.58+mid*0.24+low*0.08+envelope*0.10, 0, 1),
		math.Max(
			clamp(mid*0.50+high*0.22+low*0.16+envelope*0.12, 0, 1),
			clamp(low*0.42+mid*0.26+high*0.14+envelope*0.18, 0, 1),
		),
	)
	burst := clamp((pulseKey-ps.lastPulseKey)*(5.4+ps.pulseBias*2.6), 0, 1)
	speed := ps.adaptivePulseRate(dt, pulseKey, burst)
	peak := clamp(pulseKey*0.84+burst*(0.95+ps.pulseBias*0.95), 0, 1)
	targetPulse := math.Pow(peak, 0.55) * 0.98
	attack := 1 - math.Exp(-(34.0*speed)*dt)
	release := 1 - math.Exp(-(18.5*speed)*dt)

	if !ps.pulseActive && targetPulse > 0.12 && (burst > 0.05 || pulseKey > 0.58) {
		ps.pulseActive = true
		ps.pulseAttack = true
		ps.pulseTarget = targetPulse
		if ps.pulse < 0.02 {
			ps.pulse = 0.02
		}
	}

	if !ps.pulseActive {
		ps.pulse = 0
		return
	}

	if ps.pulseAttack {
		ps.pulse += (ps.pulseTarget - ps.pulse) * attack
		if ps.pulseTarget-ps.pulse < 0.02 || ps.pulse >= ps.pulseTarget*0.92 {
			ps.pulse = ps.pulseTarget
			ps.pulseAttack = false
		}
	} else {
		ps.pulse += (0 - ps.pulse) * release
	}

	if ps.pulse > 0.99 {
		ps.pulse = 0.99
	}
	if ps.pulse < 0.005 {
		ps.lastPulseKey = pulseKey * 0.32
		ps.resetPulseState()
		return
	}
	ps.lastPulseKey = pulseKey
}

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
	if !ps.musicOn {
		ps.speedMul = ps.baseSpeed
		ps.spawnMul = 1.0
		ps.spawnRate = ps.baseSpawn
		ps.silenced = false
		ps.resetPulseState()
		ps.beatTimer = 0
		ps.beatInterval = 0
		ps.adaptiveSpeed = 0
		ps.lastEnvelope = 0
		return
	}

	if dt < 0 {
		dt = 0
	}

	const silenceThreshold = 0.06
	if envelope < silenceThreshold {
		ps.silenced = true
		ps.spawnMul = 0
		ps.spawnRate = ps.baseSpawn
		if ps.pulseEnabled {
			ps.pulseAttack = false
			ps.updatePulse(dt, 0, 0, 0, 0)
		} else {
			ps.resetPulseState()
		}
		ps.lastEnvelope = envelope
		return
	}
	ps.silenced = false

	envelope = clamp(envelope, 0, 1)
	ps.lowEnergy = low
	ps.midEnergy = mid
	ps.highEnergy = high

	if !ps.pulseEnabled {
		ps.resetPulseState()
	} else {
		ps.updatePulse(dt, low, mid, high, envelope)
	}

	ps.lastEnvelope = envelope

	if ps.visualizer {
		if ps.pulseEnabled {
			reactiveEnergy := clamp(high*0.62+mid*0.28+low*0.10, 0, 1)
			ps.speedMul = ps.baseSpeed * clamp(1.0+(reactiveEnergy*0.85+ps.pulse*0.55)*2.1*ps.intensity, 0.6, 3.2)
			ps.spawnMul = 0.4
			ps.energyMul = clamp(low*0.6+mid*0.7+high*0.9, 0, 1)
			ps.spawnRate = ps.baseSpawn * clamp(1.0+reactiveEnergy*1.5*ps.intensity, 0.8, 5.0)
		} else {
			ps.speedMul = ps.baseSpeed * clamp(1.0+envelope*1.0*ps.intensity, 0.6, 2.0)
			ps.spawnMul = 0.4
			ps.energyMul = clamp(low*0.6+mid*0.7+high*0.9, 0, 1)
			ps.spawnRate = ps.baseSpawn * clamp(1.0+envelope*1.0*ps.intensity, 0.6, 3.0)
		}
		return
	}

	if !ps.pulseEnabled {
		pulseEnergy := clamp(high*0.55+mid*0.30+low*0.15, 0, 1)
		speedEnergy := clamp(low*0.30+mid*0.40+high*0.10, 0, 1)
		ps.energyMul = clamp(speedEnergy*6.2*ps.intensity, 0.12, 7.0)
		ps.spawnMul = clamp(pulseEnergy*8.0*ps.intensity, 0.12, 9.0)
		ps.spawnRate = ps.baseSpawn * clamp(1.0+(high*0.4+mid*0.3+low*0.2)*ps.intensity, 0.6, 3.0)
		return
	}

	densityEnergy := clamp(high*0.52+mid*0.33+low*0.15, 0, 1)
	speedEnergy := clamp(high*0.46+mid*0.30+low*0.12+ps.pulse*0.12, 0, 1)
	ps.energyMul = clamp(speedEnergy*6.2*ps.intensity, 0.18, 7.0)
	ps.spawnMul = clamp(densityEnergy*8.2*ps.intensity, 0.25, 9.0)
	ps.spawnRate = ps.baseSpawn * clamp(1.0+densityEnergy*1.75*ps.intensity, 0.9, 4.0)
}

func (ps *ParticleSystem) ResetSpectrum() {
	ps.speedMul = ps.baseSpeed
	ps.spawnMul = 1.0
	ps.spawnRate = ps.baseSpawn
	ps.bassPhase = 0
	ps.lastEnvelope = 0
	ps.resetPulseState()
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
