package rain

import (
	"math"
	"math/rand"

	"virga-player/rain/separating_frequencies"
)

const (
	maxRainMotionStep = 1.0 / 30.0
	maxRainSpawnStep  = 1.0 / 45.0
)

func (ps *ParticleSystem) HitMessage(message string, startX, row int, hidden []bool) {
	messageRunes := []rune(message)
	for _, p := range ps.particles {
		x := int(p.X)
		idx := x - startX
		if idx < 0 || idx >= len(messageRunes) || idx >= len(hidden) || hidden[idx] || messageRunes[idx] == ' ' {
			continue
		}
		topY := int(p.Y)
		bottomY := topY + p.Length - 1
		if row >= topY && row <= bottomY {
			hidden[idx] = true
		}
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	if !ps.enabled {
		ps.particles = nil
		return
	}

	motionDT := dt
	if motionDT > maxRainMotionStep {
		motionDT = maxRainMotionStep
	}
	spawnDT := dt
	if spawnDT > maxRainSpawnStep {
		spawnDT = maxRainSpawnStep
	}

	if ps.musicOn {
		if ps.silenced {
			ps.speedMul = ps.baseSpeed
			ps.spawnMul = 0
		} else {
			if ps.visualizer {
				ps.spawnVisualizerDrops()
			}
			ps.speedMul = ps.baseSpeed * ps.energyMul
		}
	}

	defaultSpeedMul := ps.speedMul
	if !ps.spawnPaused && len(ps.particles) < ps.maxSize {
		expected := ps.spawnRate * ps.spawnMul * spawnDT
		attempts := int(expected)
		if rand.Float64() < expected-float64(attempts) {
			attempts++
		}
		for i := 0; i < attempts; i++ {
			if rand.Float64() < ps.spawnChance {
				ps.spawn()
			}
		}
	}

	alive := ps.particles[:0]
	for i := range ps.particles {
		p := ps.particles[i]
		p.Age += dt

		if p.Age >= p.Life {
			continue
		}

		if p.Age >= p.Delay {
			p.VelY = p.TargetVelY
			partSpeedMul := defaultSpeedMul
			if ps.separateFreq && !ps.silenced {
				speedEnergy := separating_frequencies.LayerSpeedEnergy(ps.separateFreq, p.Layer, ps.pulse, ps.lowEnergy, ps.midEnergy, ps.highEnergy)
				partSpeedMul = ps.baseSpeed * separating_frequencies.Clamp(speedEnergy*3.2*ps.intensity, 0.12, 7.0)
			}
			p.X += p.VelX * motionDT * partSpeedMul
			p.Y += p.VelY * motionDT * partSpeedMul
		}

		if p.Age < p.GrowTime {
			progress := p.Age / p.GrowTime
			if progress > 1 {
				progress = 1
			}
			p.Length = 1 + int(progress*float64(p.TargetLength-1))
			if p.Length < 1 {
				p.Length = 1
			}
		} else {
			p.Length = p.TargetLength
		}

		opacityScale := 1.0
		if p.GrowTime > 0 {
			opacityScale = p.Age / p.GrowTime
			if opacityScale > 1 {
				opacityScale = 1
			}
		}
		if p.FadeTime > 0 && p.Age > p.Life-p.FadeTime {
			fadeProgress := (p.Age - (p.Life - p.FadeTime)) / p.FadeTime
			if fadeProgress > 1 {
				fadeProgress = 1
			}
			opacityScale *= 1 - fadeProgress
		}
		p.Opacity = int(math.Ceil(float64(p.MaxOpacity) * opacityScale))
		if p.Opacity < 1 {
			p.Opacity = 1
		}

		if p.X < -5 || p.X >= float64(ps.width+5) {
			continue
		}

		if p.Y < float64(ps.height) && p.Opacity > 0 {
			alive = append(alive, p)
		}
	}
	ps.particles = alive
}
