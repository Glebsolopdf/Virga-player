package rain

import (
	"math/rand"

	"virga-player/settings"
)

type Particle struct {
	X            float64
	Y            float64
	VelX         float64
	VelY         float64
	TargetVelY   float64
	Length       int
	TargetLength int
	Age          float64
	GrowTime     float64
	Delay        float64
	Opacity      int
}

type ParticleSystem struct {
	particles   []Particle
	width       int
	height      int
	maxSize     int
	spawnRate   float64
	baseSpawn   float64
	spawnChance float64
	spawnMul    float64
	direction   int
	baseSpeed   float64
	speedMul    float64
	musicOn     bool
	intensity   float64
	enabled     bool

	bassPhase   int
	bassTimer   float64
	prevLow     float64
	energyMul   float64
	silenced    bool
	spawnPaused bool
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(width, height int, cfg *settings.Config) *ParticleSystem {
	direction := directionFromConfig(cfg.Direction)
	ps := &ParticleSystem{
		width:       width,
		height:      height,
		maxSize:     cfg.MaxParticles,
		spawnRate:   150.0,
		baseSpawn:   150.0,
		spawnChance: 0.35,
		spawnMul:    1.0,
		particles:   make([]Particle, 0, cfg.MaxParticles+20),
		direction:   direction,
		baseSpeed:   float64(cfg.RainSpeed) / 100.0,
		speedMul:    float64(cfg.RainSpeed) / 100.0,
		musicOn:     cfg.MusicReactive,
		intensity:   float64(cfg.MusicReactiveIntensity) / 100.0,
		enabled:     cfg.RainEnabled,
	}
	if ps.enabled {
		ps.spawnInitial()
	}
	return ps
}

func directionFromConfig(mode settings.DirectionMode) int {
	switch mode {
	case settings.DirectionRightToLeft:
		return -1
	case settings.DirectionLeftToRight:
		return 1
	case settings.DirectionStraight:
		return 0
	case settings.DirectionRandom:
		return rand.Intn(3) - 1
	default:
		return rand.Intn(3) - 1
	}
}
