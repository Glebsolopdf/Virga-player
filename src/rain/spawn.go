package rain

import (
	"math/rand"

	"virga-player/rain/separating_frequencies"
)

type particleLayerProps struct {
	minLength, maxLength int
	minSpeed, maxSpeed   float64
	opacity              int
	growTime, delay      float64
}

func (ps *ParticleSystem) spawnInitial() {
	if !ps.enabled {
		return
	}
	for i := 0; i < ps.maxSize/10; i++ {
		ps.spawn()
	}
}

func (ps *ParticleSystem) chooseLayer() int {
	if !ps.separateFreq || !ps.musicOn {
		r := rand.Float64()
		if r < 0.4 {
			return layerNear
		}
		if r < 0.75 {
			return layerMid
		}
		return layerFar
	}

	weightLow := ps.lowEnergy * 3.0
	weightMid := ps.midEnergy * 1.2
	weightHigh := ps.highEnergy * 0.7
	total := weightLow + weightMid + weightHigh
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
	if r < weightLow*0.55 {
		return layerVeryNear
	}
	if r < weightLow {
		return layerNear
	}
	if r < weightLow+weightMid*0.6 {
		return layerMid
	}
	if r < weightLow+weightMid {
		return layerFar
	}
	return layerVeryFar
}

func layerProps(layer int) particleLayerProps {
	switch layer {
	case layerVeryNear:
		return particleLayerProps{minLength: 6, maxLength: 10, minSpeed: 24.0, maxSpeed: 30.0, opacity: 6, growTime: 0.12, delay: 0.04}
	case layerNear:
		return particleLayerProps{minLength: 4, maxLength: 8, minSpeed: 20.0, maxSpeed: 26.0, opacity: 5, growTime: 0.15, delay: 0.05}
	case layerMid:
		return particleLayerProps{minLength: 3, maxLength: 6, minSpeed: 16.0, maxSpeed: 20.0, opacity: 4, growTime: 0.20, delay: 0.08}
	case layerFar:
		return particleLayerProps{minLength: 2, maxLength: 4, minSpeed: 12.0, maxSpeed: 16.0, opacity: 3, growTime: 0.26, delay: 0.10}
	default:
		return particleLayerProps{minLength: 1, maxLength: 3, minSpeed: 9.0, maxSpeed: 13.0, opacity: 2, growTime: 0.32, delay: 0.15}
	}
}

func (ps *ParticleSystem) estimateLife(startY, targetVelY, delay float64) float64 {
	distance := float64(ps.height) - startY
	if distance < 0 {
		distance = 0
	}
	if targetVelY <= 0 || ps.baseSpeed <= 0 {
		return 0.9 + rand.Float64()*0.7
	}
	travelTime := delay + distance/(targetVelY*ps.baseSpeed)
	life := travelTime - 1.0
	if life < 0.6 {
		life = 0.6 + rand.Float64()*0.4
	}
	return life
}

func newParticle(x, y, velX, targetVelY float64, targetLength, layer int, props particleLayerProps, life float64) Particle {
	fadeTime := 0.18 + rand.Float64()*0.12
	if fadeTime > life*0.45 {
		fadeTime = life * 0.45
	}
	return Particle{
		X:            x,
		Y:            y,
		VelX:         velX,
		VelY:         0.0,
		TargetVelY:   targetVelY,
		Length:       1,
		TargetLength: targetLength,
		Age:          0.0,
		GrowTime:     props.growTime,
		Delay:        props.delay,
		Life:         life,
		FadeTime:     fadeTime,
		Opacity:      1,
		MaxOpacity:   props.opacity,
		Layer:        layer,
	}
}

func (ps *ParticleSystem) spawn() {
	if len(ps.particles) >= ps.maxSize {
		return
	}

	startX := float64(rand.Intn(ps.width))
	startY := float64(rand.Intn(5) - 5)
	velX := 0.0

	if ps.direction != 0 && rand.Intn(8) == 0 {
		if ps.direction > 0 {
			startX = -2
			startY = float64(rand.Intn(ps.height / 2))
			velX = 6.0 + rand.Float64()*2.0
		} else {
			startX = float64(ps.width + 2)
			startY = float64(rand.Intn(ps.height / 2))
			velX = -(6.0 + rand.Float64()*2.0)
		}
	} else if ps.direction != 0 {
		velX = float64(ps.direction) * (4.0 + rand.Float64()*2.0)
	}

	layer := ps.chooseLayer()
	props := layerProps(layer)
	minLength := props.minLength
	maxLength := props.maxLength
	layerEnergy := separating_frequencies.LayerEnergy(ps.separateFreq, layer, ps.pulse, ps.lowEnergy, ps.midEnergy, ps.highEnergy)
	energyLength := int(layerEnergy * ps.intensity * 4.0)
	finalLength := minLength + rand.Intn(maxLength-minLength+1) + energyLength
	if finalLength > maxLength+3 {
		finalLength = maxLength + 3
	}
	speedEnergy := separating_frequencies.LayerSpeedEnergy(ps.separateFreq, layer, ps.pulse, ps.lowEnergy, ps.midEnergy, ps.highEnergy)
	targetVelY := props.minSpeed + rand.Float64()*(props.maxSpeed-props.minSpeed) + speedEnergy*12.0
	life := ps.estimateLife(startY, targetVelY, props.delay)
	particle := newParticle(
		startX,
		startY,
		velX,
		targetVelY,
		finalLength,
		layer,
		props,
		life,
	)
	ps.particles = append(ps.particles, particle)
}

func (ps *ParticleSystem) spawnVisualizerDrops() {
	if len(ps.particles) >= ps.maxSize || !ps.enabled {
		return
	}

	sections := []struct {
		energy float64
		x0, x1 int
		speed  float64
		length int
	}{
		{ps.lowEnergy, 0, ps.width / 3, 18.0, 5},
		{ps.midEnergy, ps.width / 3, ps.width * 2 / 3, 20.0, 4},
		{ps.highEnergy, ps.width * 2 / 3, ps.width, 22.0, 3},
	}

	for idx, section := range sections {
		if section.energy < 0.08 {
			continue
		}
		spawnCount := int(section.energy * 3.0 * ps.intensity)
		if rand.Float64() < section.energy*0.6*ps.intensity {
			spawnCount++
		}
		for i := 0; i < spawnCount && len(ps.particles) < ps.maxSize; i++ {
			width := section.x1 - section.x0
			if width <= 0 {
				continue
			}
			x := float64(section.x0 + rand.Intn(width))
			targetVelY := section.speed + section.energy*12.0
			energyLength := int((ps.midEnergy*1.2 + ps.highEnergy*0.8 + ps.lowEnergy*0.5) * ps.intensity * 2.0)
			if energyLength > 3 {
				energyLength = 3
			}
			finalLength := section.length + rand.Intn(2) + energyLength
			if finalLength > section.length+3 {
				finalLength = section.length + 3
			}
			layer := ps.chooseLayer()
			if ps.separateFreq {
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
			fullVelY := targetVelY + ps.highEnergy*6.0*ps.intensity
			life := ps.estimateLife(float64(rand.Intn(5)-5), fullVelY, props.delay)
			particle := newParticle(
				x,
				float64(rand.Intn(5)-5),
				0.0,
				fullVelY,
				finalLength,
				layer,
				props,
				life,
			)
			ps.particles = append(ps.particles, particle)
		}
	}
}

func (ps *ParticleSystem) SpawnWashDrops(width, height int) {
	if !ps.enabled {
		return
	}
	centerX := width / 2
	positions := []int{centerX - 6, centerX - 3, centerX, centerX + 3, centerX + 6}
	for _, x := range positions {
		if len(ps.particles) >= ps.maxSize {
			break
		}
		if x < 0 || x >= width {
			continue
		}
		velX := 0.0
		if ps.direction != 0 {
			velX = float64(ps.direction) * 2.5
		}
		layer := ps.chooseLayer()
		props := layerProps(layer)
		targetVelY := 22.0
		life := ps.estimateLife(-3, targetVelY, props.delay)
		ps.particles = append(ps.particles, newParticle(
			float64(x),
			-3,
			velX,
			targetVelY,
			3,
			layer,
			props,
			life,
		))
	}
}

func (ps *ParticleSystem) SpawnMessageDrops(startX, row int, message string, hidden []bool) {
	if !ps.enabled {
		return
	}
	for i, ch := range message {
		if len(ps.particles) >= ps.maxSize {
			break
		}
		if ch == ' ' || (i < len(hidden) && hidden[i]) {
			continue
		}
		x := float64(startX + i)
		velX := 0.0
		if ps.direction != 0 {
			velX = float64(ps.direction) * 2.0
		}
		layer := ps.chooseLayer()
		props := layerProps(layer)
		fullVelY := 22.0
		life := ps.estimateLife(float64(row), fullVelY, props.delay)
		ps.particles = append(ps.particles, newParticle(
			x,
			float64(row),
			velX,
			fullVelY,
			3,
			layer,
			props,
			life,
		))
		if i < len(hidden) {
			hidden[i] = true
		}
	}
}
