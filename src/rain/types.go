package rain

import (
	"math/rand"

	"virga-player/settings"
)

const (
	layerVeryNear = iota
	layerNear
	layerMid
	layerFar
	layerVeryFar
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
	Life         float64
	FadeTime     float64
	Opacity      int
	MaxOpacity   int
	Layer        int
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
	visualizer  bool
	intensity   float64
	lifeMul     float64
	enabled     bool

	lowEnergy     float64
	midEnergy     float64
	highEnergy    float64
	bassPhase     int
	bassTimer     float64
	prevLow       float64
	lastEnvelope  float64
	energyMul     float64
	pulse         float64
	pulseBias     float64
	pulseSpeed    float64
	pulseTarget   float64
	lastPulseKey  float64
	beatTimer     float64
	beatInterval  float64
	adaptiveSpeed float64
	pulseActive   bool
	pulseAttack   bool
	pulseEnabled  bool
	separateFreq  bool
	silenced      bool
	spawnPaused   bool
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(width, height int, cfg *settings.Config) *ParticleSystem {
	ps := &ParticleSystem{
		width:       width,
		height:      height,
		spawnChance: 0.35,
		particles:   make([]Particle, 0, cfg.MaxParticles+20),
	}
	ps.ApplyConfig(cfg)
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
