package motion

import "math"

type RainState struct {
	Timer      float64
	OffsetX    float64
	OffsetY    float64
	Resistance float64
}

func UpdateRain(state *RainState, dt, pressure float64, invert bool) {
	if dt < 0 {
		dt = 0
	}
	pressure = clampFloat(pressure, 0, 1)
	target := clampFloat(math.Pow(pressure, 0.75)*0.8+pressure*0.22, 0, 1)
	rise := 1 - math.Exp(-5.2*dt)
	fall := 1 - math.Exp(-1.2*dt)
	alpha := rise
	if target < state.Resistance {
		alpha = fall
	}
	state.Resistance += (target - state.Resistance) * clampFloat(alpha, 0, 1)
	speedScale := 0.85 + pressure*1.4
	if speedScale < 0.85 {
		speedScale = 0.85
	}
	state.Timer += dt*speedScale + state.Resistance*0.1 + pressure*0.15
	amplitude := 3.6 + pressure*2.2
	amplitudeY := 2.4 + pressure*1.6
	offsetX := ((math.Sin(state.Timer*1.1) + math.Sin(state.Timer*0.4+1.3)) * 0.5) * state.Resistance * amplitude
	offsetY := ((math.Cos(state.Timer*0.85) + math.Cos(state.Timer*0.55+0.7)) * 0.5) * state.Resistance * amplitudeY
	if invert {
		offsetX = -offsetX
		offsetY = -offsetY
	}
	state.OffsetX = offsetX
	state.OffsetY = offsetY
}
