package rain

import "math"

func applySpectrum(ps *ParticleSystem, dt, low, mid, high, envelope float64) {
	if !ps.musicOn {
		resetSpectrum(ps)
		ps.beatTimer = 0
		ps.beatInterval = 0
		ps.adaptiveSpeed = 0
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
			updatePulse(ps, dt, 0, 0, 0, 0)
		} else {
			resetPulseState(ps)
		}
		ps.lastEnvelope = envelope
		return
	}
	ps.silenced = false

	envelope = max(0., min(envelope, 1.))
	ps.lowEnergy = low
	ps.midEnergy = mid
	ps.highEnergy = high

	if !ps.pulseEnabled {
		resetPulseState(ps)
	} else {
		updatePulse(ps, dt, low, mid, high, envelope)
	}

	ps.lastEnvelope = envelope

	if ps.visualizer {
		if ps.pulseEnabled {
			reactiveEnergy := max(0., min(high*0.62+mid*0.28+low*0.10, 1.))
			ps.speedMul = ps.baseSpeed * max(0.6, min(1.0+(reactiveEnergy*0.85+ps.pulse*0.55)*2.1*ps.intensity, 3.2))
			ps.spawnMul = 0.4
			ps.energyMul = max(0., min(low*0.6+mid*0.7+high*0.9, 1.))
			ps.spawnRate = ps.baseSpawn * max(0.8, min(1.0+reactiveEnergy*1.5*ps.intensity, 5.0))
		} else {
			ps.speedMul = ps.baseSpeed * max(0.6, min(1.0+envelope*1.0*ps.intensity, 2.0))
			ps.spawnMul = 0.4
			ps.energyMul = max(0., min(low*0.6+mid*0.7+high*0.9, 1.))
			ps.spawnRate = ps.baseSpawn * max(0.6, min(1.0+envelope*1.0*ps.intensity, 3.0))
		}
		return
	}

	if !ps.pulseEnabled {
		pulseEnergy := max(0., min(high*0.55+mid*0.30+low*0.15, 1.))
		speedEnergy := max(0., min(low*0.30+mid*0.40+high*0.10, 1.))
		ps.energyMul = max(0.12, min(speedEnergy*6.2*ps.intensity, 7.0))
		ps.spawnMul = max(0.12, min(pulseEnergy*8.0*ps.intensity, 9.0))
		ps.spawnRate = ps.baseSpawn * max(0.6, min(1.0+(high*0.4+mid*0.3+low*0.2)*ps.intensity, 3.0))
		return
	}

	densityEnergy := max(0., min(high*0.52+mid*0.33+low*0.15, 1.))
	speedEnergy := max(0., min(high*0.46+mid*0.30+low*0.12+ps.pulse*0.12, 1.))
	ps.energyMul = max(0.18, min(speedEnergy*6.2*ps.intensity, 7.0))
	ps.spawnMul = max(0.25, min(densityEnergy*8.2*ps.intensity, 9.0))
	ps.spawnRate = ps.baseSpawn * max(0.9, min(1.0+densityEnergy*1.75*ps.intensity, 4.0))
}

func resetSpectrum(ps *ParticleSystem) {
	ps.speedMul = ps.baseSpeed
	ps.spawnMul = 1.0
	ps.spawnRate = ps.baseSpawn
	ps.lastEnvelope = 0
	resetPulseState(ps)
	ps.silenced = false
}

func resetPulseState(ps *ParticleSystem) {
	ps.pulse = 0
	ps.pulseTarget = 0
	ps.pulseActive = false
	ps.pulseAttack = false
}

func updatePulse(ps *ParticleSystem, dt, low, mid, high, envelope float64) {
	pulseKey := math.Max(
		max(0., min(high*0.58+mid*0.24+low*0.08+envelope*0.10, 1.)),
		math.Max(
			max(0., min(mid*0.50+high*0.22+low*0.16+envelope*0.12, 1.)),
			max(0., min(low*0.42+mid*0.26+high*0.14+envelope*0.18, 1.)),
		),
	)
	burst := max(0., min((pulseKey-ps.lastPulseKey)*(5.4+ps.pulseBias*2.6), 1.))
	speed := adaptivePulseRate(ps, dt, pulseKey, burst)
	peak := max(0., min(pulseKey*0.84+burst*(0.95+ps.pulseBias*0.95), 1.))
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
		resetPulseState(ps)
		return
	}
	ps.lastPulseKey = pulseKey
}

func adaptivePulseRate(ps *ParticleSystem, dt, pulseKey, burst float64) float64 {
	ps.beatTimer += dt
	base := max(0.25, min(ps.pulseSpeed, 3.0))
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
		tempoFactor := max(0.72, min(0.46/ps.beatInterval, 1.75))
		intensityFactor := 0.88 + max(0., min(pulseKey*0.46+burst*0.90, 1.))*0.44
		target = max(0.25, min(base*tempoFactor*intensityFactor, 3.0))
	}
	smooth := 1 - math.Exp(-6.0*dt)
	ps.adaptiveSpeed += (target - ps.adaptiveSpeed) * smooth
	return max(0.25, min(ps.adaptiveSpeed, 3.0))
}
