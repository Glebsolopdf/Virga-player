package spectrum

import "math"

func updatePulse(state *State, dt, low, mid, high, envelope float64) {
	pulseKey := math.Max(
		clamp(high*0.58+mid*0.24+low*0.08+envelope*0.10, 0, 1),
		math.Max(
			clamp(mid*0.50+high*0.22+low*0.16+envelope*0.12, 0, 1),
			clamp(low*0.42+mid*0.26+high*0.14+envelope*0.18, 0, 1),
		),
	)
	burst := clamp((pulseKey-state.LastPulseKey)*(5.4+state.PulseBias*2.6), 0, 1)
	speed := adaptivePulseRate(state, dt, pulseKey, burst)
	peak := clamp(pulseKey*0.84+burst*(0.95+state.PulseBias*0.95), 0, 1)
	targetPulse := math.Pow(peak, 0.55) * 0.98
	attack := 1 - math.Exp(-(34.0*speed)*dt)
	release := 1 - math.Exp(-(18.5*speed)*dt)

	if !state.PulseActive && targetPulse > 0.12 && (burst > 0.05 || pulseKey > 0.58) {
		state.PulseActive = true
		state.PulseAttack = true
		state.PulseTarget = targetPulse
		if state.Pulse < 0.02 {
			state.Pulse = 0.02
		}
	}

	if !state.PulseActive {
		state.Pulse = 0
		return
	}

	if state.PulseAttack {
		state.Pulse += (state.PulseTarget - state.Pulse) * attack
		if state.PulseTarget-state.Pulse < 0.02 || state.Pulse >= state.PulseTarget*0.92 {
			state.Pulse = state.PulseTarget
			state.PulseAttack = false
		}
	} else {
		state.Pulse += (0 - state.Pulse) * release
	}

	if state.Pulse > 0.99 {
		state.Pulse = 0.99
	}
	if state.Pulse < 0.005 {
		state.LastPulseKey = pulseKey * 0.32
		resetPulseState(state)
		return
	}
	state.LastPulseKey = pulseKey
}

func adaptivePulseRate(state *State, dt, pulseKey, burst float64) float64 {
	state.BeatTimer += dt
	base := clamp(state.PulseSpeed, 0.25, 3.0)
	if state.AdaptiveSpeed == 0 {
		state.AdaptiveSpeed = base
	}

	if burst > 0.05 && pulseKey > 0.48 {
		if state.BeatTimer >= 0.18 && state.BeatTimer <= 1.20 {
			if state.BeatInterval == 0 {
				state.BeatInterval = state.BeatTimer
			} else {
				state.BeatInterval += (state.BeatTimer - state.BeatInterval) * 0.34
			}
		}
		state.BeatTimer = 0
	}

	target := base
	if state.BeatInterval > 0 {
		tempoFactor := clamp(0.46/state.BeatInterval, 0.72, 1.75)
		intensityFactor := 0.88 + clamp(pulseKey*0.46+burst*0.90, 0, 1)*0.44
		target = clamp(base*tempoFactor*intensityFactor, 0.25, 3.0)
	}
	smooth := 1 - math.Exp(-6.0*dt)
	state.AdaptiveSpeed += (target - state.AdaptiveSpeed) * smooth
	return clamp(state.AdaptiveSpeed, 0.25, 3.0)
}
