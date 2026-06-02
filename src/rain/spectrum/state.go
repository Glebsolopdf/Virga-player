package spectrum

type State struct {
	BaseSpeed     float64
	BaseSpawn     float64
	Intensity     float64
	PulseBias     float64
	PulseSpeed    float64
	SpeedMul      float64
	SpawnMul      float64
	SpawnRate     float64
	EnergyMul     float64
	LowEnergy     float64
	MidEnergy     float64
	HighEnergy    float64
	LastEnvelope  float64
	Pulse         float64
	PulseTarget   float64
	LastPulseKey  float64
	BeatTimer     float64
	BeatInterval  float64
	AdaptiveSpeed float64
	Silenced      bool
	PulseActive   bool
	PulseAttack   bool
	PulseEnabled  bool
	Visualizer    bool
}

func Reset(state *State) {
	state.SpeedMul = state.BaseSpeed
	state.SpawnMul = 1.0
	state.SpawnRate = state.BaseSpawn
	state.LastEnvelope = 0
	resetPulseState(state)
	state.Silenced = false
}

func resetPulseState(state *State) {
	state.Pulse = 0
	state.PulseTarget = 0
	state.PulseActive = false
	state.PulseAttack = false
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
