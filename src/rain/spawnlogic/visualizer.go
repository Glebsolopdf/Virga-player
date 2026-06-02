package spawnlogic

import "math/rand"

func PlanVisualizerSpawns(state State) []ParticlePlan {
	remaining := state.MaxSize - state.ParticleCount
	if remaining <= 0 || !state.Enabled {
		return nil
	}

	sections := []struct {
		energy float64
		x0, x1 int
		speed  float64
		length int
	}{
		{state.LowEnergy, 0, state.Width / 3, 18.0, 5},
		{state.MidEnergy, state.Width / 3, state.Width * 2 / 3, 20.0, 4},
		{state.HighEnergy, state.Width * 2 / 3, state.Width, 22.0, 3},
	}

	plans := make([]ParticlePlan, 0, remaining)
	for idx, section := range sections {
		if section.energy < 0.08 {
			continue
		}
		spawnCount := int(section.energy * 3.0 * state.Intensity)
		if rand.Float64() < section.energy*0.6*state.Intensity {
			spawnCount++
		}
		for i := 0; i < spawnCount && len(plans) < remaining; i++ {
			width := section.x1 - section.x0
			if width <= 0 {
				continue
			}
			x := float64(section.x0 + rand.Intn(width))
			startY := float64(rand.Intn(2) - 1)
			targetVelY := section.speed + section.energy*12.0
			energyLength := int((state.MidEnergy*1.2 + state.HighEnergy*0.8 + state.LowEnergy*0.5) * state.Intensity * 2.0)
			if energyLength > 3 {
				energyLength = 3
			}
			finalLength := section.length + rand.Intn(2) + energyLength
			if finalLength > section.length+3 {
				finalLength = section.length + 3
			}
			layer := chooseLayer(state)
			if state.SeparateFreq {
				switch idx {
				case 0:
					layer = layerNear
				case 1:
					layer = layerMid
				case 2:
					layer = layerFar
				}
			}
			props := layerProps(layer)
			fullVelY := targetVelY + state.HighEnergy*6.0*state.Intensity
			life := estimateLife(state.Height, startY, fullVelY, props.Delay, state.BaseSpeed, state.LifeMul)
			plans = append(plans, ParticlePlan{
				X:            x,
				Y:            startY,
				TargetVelY:   fullVelY,
				TargetLength: finalLength,
				Layer:        layer,
				GrowTime:     props.GrowTime,
				Delay:        props.Delay,
				Life:         life,
				MaxOpacity:   props.Opacity,
			})
		}
	}

	return plans
}

func PlanMessageSpawns(state State, startX, row int, message string, hidden []bool) ([]ParticlePlan, []int) {
	remaining := state.MaxSize - state.ParticleCount
	if !state.Enabled || remaining <= 0 {
		return nil, nil
	}

	messageRunes := []rune(message)
	plans := make([]ParticlePlan, 0, remaining)
	hiddenIndices := make([]int, 0, remaining)
	for i, ch := range messageRunes {
		if len(plans) >= remaining {
			break
		}
		if ch == ' ' || (i < len(hidden) && hidden[i]) {
			continue
		}
		velX := 0.0
		if state.Direction != 0 {
			velX = float64(state.Direction) * 2.0
		}
		layer := chooseLayer(state)
		props := layerProps(layer)
		fullVelY := 22.0
		life := estimateLife(state.Height, float64(row), fullVelY, props.Delay, state.BaseSpeed, state.LifeMul)
		plans = append(plans, ParticlePlan{
			X:            float64(startX + i),
			Y:            float64(row),
			VelX:         velX,
			TargetVelY:   fullVelY,
			TargetLength: 3,
			Layer:        layer,
			GrowTime:     props.GrowTime,
			Delay:        props.Delay,
			Life:         life,
			MaxOpacity:   props.Opacity,
		})
		hiddenIndices = append(hiddenIndices, i)
	}

	return plans, hiddenIndices
}
