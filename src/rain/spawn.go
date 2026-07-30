package rain

import "math/rand"

func (ps *ParticleSystem) spawnInitial() {
	for i := 0; i < initialCount(ps.spawnCfg()); i++ {
		ps.spawn()
	}
}

func newParticle(plan particlePlan) Particle {
	fadeTime := 0.24 + rand.Float64()*0.16
	if fadeTime > plan.Life*0.45 {
		fadeTime = plan.Life * 0.45
	}
	return Particle{
		X:            plan.X,
		Y:            plan.Y,
		VelX:         plan.VelX,
		VelY:         0.0,
		TargetVelY:   plan.TargetVelY,
		Length:       1,
		TargetLength: plan.TargetLength,
		Age:          0.0,
		GrowTime:     plan.GrowTime,
		Delay:        plan.Delay,
		Life:         plan.Life,
		FadeTime:     fadeTime,
		Opacity:      1,
		MaxOpacity:   plan.MaxOpacity,
		Layer:        plan.Layer,
	}
}

func (ps *ParticleSystem) spawn() {
	plan, ok := planSpawn(ps.spawnCfg())
	if !ok {
		return
	}
	ps.particles = append(ps.particles, newParticle(plan))
}

func (ps *ParticleSystem) spawnVisualizerDrops() {
	for _, plan := range planVisualizerSpawns(ps.spawnCfg()) {
		if len(ps.particles) >= ps.maxSize {
			break
		}
		ps.particles = append(ps.particles, newParticle(plan))
	}
}

func (ps *ParticleSystem) SpawnMessageDrops(startX, row int, message string, hidden []bool) {
	plans, hiddenIndices := planMessageSpawns(ps.spawnCfg(), startX, row, message, hidden)
	for _, idx := range hiddenIndices {
		if idx < len(hidden) {
			hidden[idx] = true
		}
	}
	for _, plan := range plans {
		if len(ps.particles) >= ps.maxSize {
			break
		}
		ps.particles = append(ps.particles, newParticle(plan))
	}
}

func (ps *ParticleSystem) spawnCfg() spawnState {
	return spawnState{
		Width:         ps.width,
		Height:        ps.height,
		MaxSize:       ps.maxSize,
		ParticleCount: len(ps.particles),
		Direction:     ps.direction,
		SeparateFreq:  ps.separateFreq,
		MusicOn:       ps.musicOn,
		Enabled:       ps.enabled,
		LowEnergy:     ps.lowEnergy,
		MidEnergy:     ps.midEnergy,
		HighEnergy:    ps.highEnergy,
		Pulse:         ps.pulse,
		Intensity:     ps.intensity,
		BaseSpeed:     ps.baseSpeed,
		LifeMul:       ps.lifeMul,
	}
}
