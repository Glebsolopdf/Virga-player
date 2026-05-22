package artwork

import "math"

func (a *Artwork) SetAnimationEnabled(enabled bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.AnimationEnabled = enabled
	if !enabled {
		a.Fade = 1
		a.resetPulseState()
	}
}

func (a *Artwork) resetPulseState() {
	a.Pulse = 0
	a.pulseTarget = 0
	a.pulseActive = false
	a.pulseAttacking = false
}

func (a *Artwork) adaptivePulseSpeed(dt, pulseKey, burst, baseSpeed float64) float64 {
	a.beatTimer += dt
	base := clampFloat(baseSpeed, 0.25, 3.0)
	if a.adaptiveSpeed == 0 {
		a.adaptiveSpeed = base
	}

	if burst > 0.05 && pulseKey > 0.48 {
		if a.beatTimer >= 0.18 && a.beatTimer <= 1.20 {
			if a.beatInterval == 0 {
				a.beatInterval = a.beatTimer
			} else {
				a.beatInterval += (a.beatTimer - a.beatInterval) * 0.34
			}
		}
		a.beatTimer = 0
	}

	target := base
	if a.beatInterval > 0 {
		tempoFactor := clampFloat(0.46/a.beatInterval, 0.72, 1.75)
		intensityFactor := 0.88 + clampFloat(pulseKey*0.50+burst*0.85, 0, 1)*0.42
		target = clampFloat(base*tempoFactor*intensityFactor, 0.25, 3.0)
	}
	smooth := 1 - math.Exp(-6.0*dt)
	a.adaptiveSpeed += (target - a.adaptiveSpeed) * smooth
	return clampFloat(a.adaptiveSpeed, 0.25, 3.0)
}

func (a *Artwork) updatePulse(dt, low, mid, high, envelope, pulseSpeed float64) {
	pulseKey := math.Max(
		clampFloat(high*0.52+mid*0.24+low*0.10+envelope*0.14, 0, 1),
		math.Max(
			clampFloat(mid*0.48+high*0.22+low*0.16+envelope*0.14, 0, 1),
			clampFloat(low*0.40+mid*0.28+high*0.16+envelope*0.16, 0, 1),
		),
	)
	burst := clampFloat((pulseKey-a.lastPulseKey)*4.8, 0, 1)
	speed := a.adaptivePulseSpeed(dt, pulseKey, burst, pulseSpeed)
	peak := clampFloat(pulseKey*0.84+burst*1.05, 0, 1)
	targetPulse := math.Pow(peak, 0.55) * 0.98
	attack := 1 - math.Exp(-(34.0*speed)*dt)
	release := 1 - math.Exp(-(17.5*speed)*dt)

	if !a.pulseActive && targetPulse > 0.12 && (burst > 0.05 || pulseKey > 0.58) {
		a.pulseActive = true
		a.pulseAttacking = true
		a.pulseTarget = targetPulse
		if a.Pulse < 0.02 {
			a.Pulse = 0.02
		}
	}

	if !a.pulseActive {
		a.Pulse = 0
		a.lastPulseKey = pulseKey
		return
	}

	if a.pulseAttacking {
		a.Pulse += (a.pulseTarget - a.Pulse) * attack
		if a.pulseTarget-a.Pulse < 0.02 || a.Pulse >= a.pulseTarget*0.92 {
			a.Pulse = a.pulseTarget
			a.pulseAttacking = false
		}
	} else {
		a.Pulse += (0 - a.Pulse) * release
	}

	if a.Pulse > 0.99 {
		a.Pulse = 0.99
	}
	if a.Pulse < 0.005 {
		a.lastPulseKey = pulseKey * 0.32
		a.resetPulseState()
		return
	}
	a.lastPulseKey = pulseKey
}

func (a *Artwork) UpdateAnimation(dt, low, mid, high, envelope, pulseSpeed float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if dt < 0 {
		dt = 0
	}

	if !a.AnimationEnabled {
		a.Fade = 1
		a.resetPulseState()
		a.LastEnvelope = 0
		a.beatTimer = 0
		a.beatInterval = 0
		a.adaptiveSpeed = 0
		return
	}

	if a.Fade < 1 {
		a.Fade += dt * 1.1
		if a.Fade > 1 {
			a.Fade = 1
		}
	}
	low = clampFloat(low, 0, 1)
	mid = clampFloat(mid, 0, 1)
	high = clampFloat(high, 0, 1)
	envelope = clampFloat(envelope, 0, 1)
	a.updatePulse(dt, low, mid, high, envelope, pulseSpeed)
	a.LastEnvelope = envelope
}

func (a *Artwork) UpdateRainResistance(dt, pressure float64, invert bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if dt < 0 {
		dt = 0
	}
	pressure = clampFloat(pressure, 0, 1)
	target := clampFloat(math.Pow(pressure, 0.75)*0.8+pressure*0.22, 0, 1)
	rise := 1 - math.Exp(-5.2*dt)
	fall := 1 - math.Exp(-1.2*dt)
	alpha := rise
	if target < a.RainResistance {
		alpha = fall
	}
	a.RainResistance += (target - a.RainResistance) * clampFloat(alpha, 0, 1)
	speedScale := 0.85 + pressure*1.4
	if speedScale < 0.85 {
		speedScale = 0.85
	}
	a.RainTimer += dt*speedScale + a.RainResistance*0.1 + pressure*0.15
	amplitude := 3.6 + pressure*2.2
	amplitudeY := 2.4 + pressure*1.6
	offsetX := ((math.Sin(a.RainTimer*1.1) + math.Sin(a.RainTimer*0.4+1.3)) * 0.5) * a.RainResistance * amplitude
	offsetY := ((math.Cos(a.RainTimer*0.85) + math.Cos(a.RainTimer*0.55+0.7)) * 0.5) * a.RainResistance * amplitudeY
	if invert {
		offsetX = -offsetX
		offsetY = -offsetY
	}
	a.RainOffsetX = offsetX
	a.RainOffsetY = offsetY
}
