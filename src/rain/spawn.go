package rain

import "math/rand"

func (ps *ParticleSystem) spawnInitial() {
	if !ps.enabled {
		return
	}
	for i := 0; i < ps.maxSize/10; i++ {
		ps.spawn()
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

	minLength := 2
	maxLength := 5
	finalLength := minLength + rand.Intn(maxLength-minLength+1)
	particle := Particle{
		X:            startX,
		Y:            startY,
		VelX:         velX,
		VelY:         0.0,
		TargetVelY:   18.0 + rand.Float64()*4.0,
		Length:       1,
		TargetLength: finalLength,
		Age:          0.0,
		GrowTime:     0.3,
		Delay:        0.1,
		Opacity:      3,
	}
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

	for _, section := range sections {
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
			finalLength := section.length + rand.Intn(2)
			particle := Particle{
				X:            x,
				Y:            float64(rand.Intn(5) - 5),
				VelX:         0.0,
				VelY:         0.0,
				TargetVelY:   targetVelY,
				Length:       1,
				TargetLength: finalLength,
				Age:          0.0,
				GrowTime:     0.2,
				Delay:        0.05,
				Opacity:      3,
			}
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
		ps.particles = append(ps.particles, Particle{
			X:            float64(x),
			Y:            -3,
			VelX:         velX,
			VelY:         0.0,
			TargetVelY:   22.0,
			Length:       1,
			TargetLength: 3,
			Age:          0.0,
			GrowTime:     0.3,
			Delay:        0.1,
			Opacity:      3,
		})
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
		ps.particles = append(ps.particles, Particle{
			X:            x,
			Y:            float64(row),
			VelX:         velX,
			VelY:         0.0,
			TargetVelY:   22.0,
			Length:       1,
			TargetLength: 3,
			Age:          0.0,
			GrowTime:     0.3,
			Delay:        0.1,
			Opacity:      3,
		})
		if i < len(hidden) {
			hidden[i] = true
		}
	}
}
