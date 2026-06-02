package spawnlogic

import (
	"math/rand"

	"virga-player/rain/separating_frequencies"
)

const (
	layerVeryNear = iota
	layerNear
	layerMid
	layerFar
	layerVeryFar
)

type LayerProps struct {
	MinLength int
	MaxLength int
	MinSpeed  float64
	MaxSpeed  float64
	Opacity   int
	GrowTime  float64
	Delay     float64
}

type State struct {
	Width         int
	Height        int
	MaxSize       int
	ParticleCount int
	Direction     int
	SeparateFreq  bool
	MusicOn       bool
	Enabled       bool
	LowEnergy     float64
	MidEnergy     float64
	HighEnergy    float64
	Pulse         float64
	Intensity     float64
	BaseSpeed     float64
	LifeMul       float64
}

func chooseLayer(state State) int {
	if !state.SeparateFreq || !state.MusicOn {
		r := rand.Float64()
		if r < 0.4 {
			return layerNear
		}
		if r < 0.75 {
			return layerMid
		}
		return layerFar
	}

	weightHigh := state.HighEnergy * 2.6
	weightMid := state.MidEnergy * 1.4
	weightLow := state.LowEnergy * 2.2

	// Keep a persistent foreground presence in separated mode so near layers
	// do not disappear when low frequencies dominate.
	weightVeryNear := 0.20 + weightHigh*0.68 + weightMid*0.18
	weightNear := 0.28 + weightHigh*0.42 + weightMid*0.36 + weightLow*0.08
	weightMidLayer := 0.20 + weightMid*0.34 + weightHigh*0.14 + weightLow*0.16
	weightFar := 0.12 + weightLow*0.30 + weightMid*0.16
	weightVeryFar := 0.05 + weightLow*0.14

	total := weightVeryNear + weightNear + weightMidLayer + weightFar + weightVeryFar
	if total <= 0 {
		r := rand.Float64()
		if r < 0.4 {
			return layerNear
		}
		if r < 0.75 {
			return layerMid
		}
		return layerFar
	}

	r := rand.Float64() * total
	if r < weightVeryNear {
		return layerVeryNear
	}
	r -= weightVeryNear
	if r < weightNear {
		return layerNear
	}
	r -= weightNear
	if r < weightMidLayer {
		return layerMid
	}
	r -= weightMidLayer
	if r < weightFar {
		return layerFar
	}
	return layerVeryFar
}

func layerProps(layer int) LayerProps {
	switch layer {
	case layerVeryNear:
		return LayerProps{MinLength: 6, MaxLength: 10, MinSpeed: 24.0, MaxSpeed: 30.0, Opacity: 6, GrowTime: 0.12, Delay: 0.04}
	case layerNear:
		return LayerProps{MinLength: 4, MaxLength: 8, MinSpeed: 20.0, MaxSpeed: 26.0, Opacity: 5, GrowTime: 0.15, Delay: 0.05}
	case layerMid:
		return LayerProps{MinLength: 3, MaxLength: 6, MinSpeed: 16.0, MaxSpeed: 20.0, Opacity: 4, GrowTime: 0.20, Delay: 0.08}
	case layerFar:
		return LayerProps{MinLength: 2, MaxLength: 4, MinSpeed: 12.0, MaxSpeed: 16.0, Opacity: 4, GrowTime: 0.26, Delay: 0.10}
	default:
		return LayerProps{MinLength: 1, MaxLength: 3, MinSpeed: 9.0, MaxSpeed: 13.0, Opacity: 3, GrowTime: 0.32, Delay: 0.15}
	}
}

func estimateLife(height int, startY, targetVelY, delay, baseSpeed, lifeMul float64) float64 {
	distance := float64(height) - startY
	if distance < 0 {
		distance = 0
	}
	if targetVelY <= 0 || baseSpeed <= 0 {
		return 0.9 + rand.Float64()*0.7
	}
	travelTime := delay + distance/(targetVelY*baseSpeed)
	life := travelTime * lifeMul
	if life < 0.6 {
		life = 0.6 + rand.Float64()*0.4
	}
	return life
}

func layerEnergy(state State, layer int) float64 {
	return separating_frequencies.LayerEnergy(state.SeparateFreq, layer, state.Pulse, state.LowEnergy, state.MidEnergy, state.HighEnergy)
}

func layerSpeedEnergy(state State, layer int) float64 {
	return separating_frequencies.LayerSpeedEnergy(state.SeparateFreq, layer, state.Pulse, state.LowEnergy, state.MidEnergy, state.HighEnergy)
}
