package spectrum

func Apply(state *State, dt, low, mid, high, envelope float64, musicOn bool) {
	if !musicOn {
		Reset(state)
		state.BeatTimer = 0
		state.BeatInterval = 0
		state.AdaptiveSpeed = 0
		return
	}

	if dt < 0 {
		dt = 0
	}

	const silenceThreshold = 0.06
	if envelope < silenceThreshold {
		state.Silenced = true
		state.SpawnMul = 0
		state.SpawnRate = state.BaseSpawn
		if state.PulseEnabled {
			state.PulseAttack = false
			updatePulse(state, dt, 0, 0, 0, 0)
		} else {
			resetPulseState(state)
		}
		state.LastEnvelope = envelope
		return
	}
	state.Silenced = false

	envelope = clamp(envelope, 0, 1)
	state.LowEnergy = low
	state.MidEnergy = mid
	state.HighEnergy = high

	if !state.PulseEnabled {
		resetPulseState(state)
	} else {
		updatePulse(state, dt, low, mid, high, envelope)
	}

	state.LastEnvelope = envelope

	if state.Visualizer {
		if state.PulseEnabled {
			reactiveEnergy := clamp(high*0.62+mid*0.28+low*0.10, 0, 1)
			state.SpeedMul = state.BaseSpeed * clamp(1.0+(reactiveEnergy*0.85+state.Pulse*0.55)*2.1*state.Intensity, 0.6, 3.2)
			state.SpawnMul = 0.4
			state.EnergyMul = clamp(low*0.6+mid*0.7+high*0.9, 0, 1)
			state.SpawnRate = state.BaseSpawn * clamp(1.0+reactiveEnergy*1.5*state.Intensity, 0.8, 5.0)
		} else {
			state.SpeedMul = state.BaseSpeed * clamp(1.0+envelope*1.0*state.Intensity, 0.6, 2.0)
			state.SpawnMul = 0.4
			state.EnergyMul = clamp(low*0.6+mid*0.7+high*0.9, 0, 1)
			state.SpawnRate = state.BaseSpawn * clamp(1.0+envelope*1.0*state.Intensity, 0.6, 3.0)
		}
		return
	}

	if !state.PulseEnabled {
		pulseEnergy := clamp(high*0.55+mid*0.30+low*0.15, 0, 1)
		speedEnergy := clamp(low*0.30+mid*0.40+high*0.10, 0, 1)
		state.EnergyMul = clamp(speedEnergy*6.2*state.Intensity, 0.12, 7.0)
		state.SpawnMul = clamp(pulseEnergy*8.0*state.Intensity, 0.12, 9.0)
		state.SpawnRate = state.BaseSpawn * clamp(1.0+(high*0.4+mid*0.3+low*0.2)*state.Intensity, 0.6, 3.0)
		return
	}

	densityEnergy := clamp(high*0.52+mid*0.33+low*0.15, 0, 1)
	speedEnergy := clamp(high*0.46+mid*0.30+low*0.12+state.Pulse*0.12, 0, 1)
	state.EnergyMul = clamp(speedEnergy*6.2*state.Intensity, 0.18, 7.0)
	state.SpawnMul = clamp(densityEnergy*8.2*state.Intensity, 0.25, 9.0)
	state.SpawnRate = state.BaseSpawn * clamp(1.0+densityEnergy*1.75*state.Intensity, 0.9, 4.0)
}
