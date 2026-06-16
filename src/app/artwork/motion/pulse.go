package motion

import "math"

type PulseState struct {
	Pulse          float64
	PulseTarget    float64
	LastPulseKey   float64
	BeatTimer      float64
	BeatInterval   float64
	AdaptiveSpeed  float64
	PulseActive    bool
	PulseAttacking bool
}

func ResetPulse(state *PulseState) {
	state.Pulse = 0
	state.PulseTarget = 0
	state.PulseActive = false
	state.PulseAttacking = false
}

func UpdatePulse(state *PulseState, dt, low, mid, high, envelope, pulseSpeed float64) {
	pulseKey := math.Max(
		clampFloat(high*0.52+mid*0.24+low*0.10+envelope*0.14, 0, 1),
		math.Max(
			clampFloat(mid*0.48+high*0.22+low*0.16+envelope*0.14, 0, 1),
			clampFloat(low*0.40+mid*0.28+high*0.16+envelope*0.16, 0, 1),
		),
	)
	burst := clampFloat((pulseKey-state.LastPulseKey)*4.8, 0, 1)
	speed := adaptivePulseSpeed(state, dt, pulseKey, burst, pulseSpeed)
	peak := clampFloat(pulseKey*0.84+burst*1.05, 0, 1)
	targetPulse := math.Pow(peak, 0.55) * 0.98
	attack := 1 - math.Exp(-(34.0*speed)*dt)
	release := 1 - math.Exp(-(17.5*speed)*dt)

	if !state.PulseActive && targetPulse > 0.12 && (burst > 0.05 || pulseKey > 0.58) {
		state.PulseActive = true
		state.PulseAttacking = true
		state.PulseTarget = targetPulse
		if state.Pulse < 0.02 {
			state.Pulse = 0.02
		}
	}

	if !state.PulseActive {
		state.Pulse = 0
		state.LastPulseKey = pulseKey
		return
	}

	if state.PulseAttacking {
		state.Pulse += (state.PulseTarget - state.Pulse) * attack
		if state.PulseTarget-state.Pulse < 0.02 || state.Pulse >= state.PulseTarget*0.92 {
			state.Pulse = state.PulseTarget
			state.PulseAttacking = false
		}
	} else {
		state.Pulse += (0 - state.Pulse) * release
	}

	if state.Pulse > 0.99 {
		state.Pulse = 0.99
	}
	if state.Pulse < 0.005 {
		state.LastPulseKey = pulseKey * 0.32
		ResetPulse(state)
		return
	}
	state.LastPulseKey = pulseKey
}

func adaptivePulseSpeed(state *PulseState, dt, pulseKey, burst, baseSpeed float64) float64 {
	state.BeatTimer += dt
	base := clampFloat(baseSpeed, 0.25, 3.0)
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
		tempoFactor := clampFloat(0.46/state.BeatInterval, 0.72, 1.75)
		intensityFactor := 0.88 + clampFloat(pulseKey*0.50+burst*0.85, 0, 1)*0.42
		target = clampFloat(base*tempoFactor*intensityFactor, 0.25, 3.0)
	}
	smooth := 1 - math.Exp(-6.0*dt)
	state.AdaptiveSpeed += (target - state.AdaptiveSpeed) * smooth
	return clampFloat(state.AdaptiveSpeed, 0.25, 3.0)
}

func clampFloat(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
