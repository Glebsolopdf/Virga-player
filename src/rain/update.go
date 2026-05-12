package rain

import "math/rand"

func (ps *ParticleSystem) HitMessage(message string, startX, row int, hidden []bool) {
	for _, p := range ps.particles {
		x := int(p.X)
		idx := x - startX
		if idx < 0 || idx >= len(message) || idx >= len(hidden) || hidden[idx] || message[idx] == ' ' {
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

	moveDT := dt * ps.speedMul
	if !ps.spawnPaused && len(ps.particles) < ps.maxSize {
		expected := ps.spawnRate * ps.spawnMul * dt
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

		if p.Age >= p.Delay {
			p.VelY = p.TargetVelY
			p.X += p.VelX * moveDT
			p.Y += p.VelY * moveDT
		}

		if p.Age < p.GrowTime {
			progress := p.Age / p.GrowTime
			p.Length = 1 + int(progress*float64(p.TargetLength-1))
			if p.Length < 1 {
				p.Length = 1
			}
		} else {
			p.Length = p.TargetLength
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
