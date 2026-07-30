package rain

import "math/rand"

type particlePlan struct {
	X            float64
	Y            float64
	VelX         float64
	TargetVelY   float64
	TargetLength int
	Layer        int
	GrowTime     float64
	Delay        float64
	Life         float64
	MaxOpacity   int
}

func initialCount(state spawnState) int {
	if !state.Enabled {
		return 0
	}
	return state.MaxSize / 10
}

func planSpawn(state spawnState) (particlePlan, bool) {
	if state.ParticleCount >= state.MaxSize || state.Width <= 0 || state.Height <= 0 {
		return particlePlan{}, false
	}

	startX := float64(rand.Intn(state.Width))
	startY := float64(rand.Intn(2) - 1)
	velX := 0.0
	sideSpawnHeight := state.Height / 2
	if sideSpawnHeight < 1 {
		sideSpawnHeight = 1
	}

	if state.Direction != 0 && rand.Intn(8) == 0 {
		if state.Direction > 0 {
			startX = -2
			startY = float64(rand.Intn(sideSpawnHeight))
			velX = 6.0 + rand.Float64()*2.0
		} else {
			startX = float64(state.Width + 2)
			startY = float64(rand.Intn(sideSpawnHeight))
			velX = -(6.0 + rand.Float64()*2.0)
		}
	} else if state.Direction != 0 {
		velX = float64(state.Direction) * (4.0 + rand.Float64()*2.0)
	}

	layer := chooseLayer(state)
	props := layerProps(layer)
	energyLength := int(layerEnergy(state.SeparateFreq, layer, state.Pulse, state.LowEnergy, state.MidEnergy, state.HighEnergy) * state.Intensity * 4.0)
	finalLength := props.MinLength + rand.Intn(props.MaxLength-props.MinLength+1) + energyLength
	if finalLength > props.MaxLength+3 {
		finalLength = props.MaxLength + 3
	}
	targetVelY := props.MinSpeed + rand.Float64()*(props.MaxSpeed-props.MinSpeed) + layerSpeedEnergy(state.SeparateFreq, layer, state.Pulse, state.LowEnergy, state.MidEnergy, state.HighEnergy)*12.0
	life := estimateLife(state.Height, startY, targetVelY, props.Delay, state.BaseSpeed, state.LifeMul)

	return particlePlan{
		X:            startX,
		Y:            startY,
		VelX:         velX,
		TargetVelY:   targetVelY,
		TargetLength: finalLength,
		Layer:        layer,
		GrowTime:     props.GrowTime,
		Delay:        props.Delay,
		Life:         life,
		MaxOpacity:   props.Opacity,
	}, true
}
