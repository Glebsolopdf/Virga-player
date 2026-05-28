package artwork

import motion "virga-player/app/artwork/motion"

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
	state := a.pulseState()
	motion.ResetPulse(&state)
	a.applyPulseState(state)
}

func (a *Artwork) updatePulse(dt, low, mid, high, envelope, pulseSpeed float64) {
	state := a.pulseState()
	motion.UpdatePulse(&state, dt, low, mid, high, envelope, pulseSpeed)
	a.applyPulseState(state)
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

	state := motion.RainState{
		Timer:      a.RainTimer,
		OffsetX:    a.RainOffsetX,
		OffsetY:    a.RainOffsetY,
		Resistance: a.RainResistance,
	}
	motion.UpdateRain(&state, dt, pressure, invert)
	a.RainTimer = state.Timer
	a.RainOffsetX = state.OffsetX
	a.RainOffsetY = state.OffsetY
	a.RainResistance = state.Resistance
}

func (a *Artwork) pulseState() motion.PulseState {
	return motion.PulseState{
		Pulse:          a.Pulse,
		PulseTarget:    a.pulseTarget,
		LastPulseKey:   a.lastPulseKey,
		BeatTimer:      a.beatTimer,
		BeatInterval:   a.beatInterval,
		AdaptiveSpeed:  a.adaptiveSpeed,
		PulseActive:    a.pulseActive,
		PulseAttacking: a.pulseAttacking,
	}
}

func (a *Artwork) applyPulseState(state motion.PulseState) {
	a.Pulse = state.Pulse
	a.pulseTarget = state.PulseTarget
	a.lastPulseKey = state.LastPulseKey
	a.beatTimer = state.BeatTimer
	a.beatInterval = state.BeatInterval
	a.adaptiveSpeed = state.AdaptiveSpeed
	a.pulseActive = state.PulseActive
	a.pulseAttacking = state.PulseAttacking
}
