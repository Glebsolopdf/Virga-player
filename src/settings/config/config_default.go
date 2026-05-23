package config

func DefaultConfig() *Config {
	return &Config{
		FPS:                    60,
		MaxParticles:           220,
		RainSpeed:              100,
		RainLifetime:           100,
		PulseSpeed:             100,
		PulseMode:              PulseModeRain,
		RainEnabled:            true,
		RainPulse:              100,
		RainPulseEnabled:       true,
		SeparateFrequencies:    false,
		Debug:                  false,
		MusicReactive:          false,
		MusicReactiveIntensity: 100,
		RainVisualizer:         false,
		CoverAnimation:         false,
		MusicPlayerAnimation:   false,
		MusicPlayerIntensity:   100,
		MusicPlayerInvert:      false,
		RainInFrontOfPlayer:    true,
		Direction:              DirectionRandom,
		Player:                 false,
	}
}

func (c *Config) Clone() *Config {
	return &Config{
		FPS:                    c.FPS,
		MaxParticles:           c.MaxParticles,
		RainSpeed:              c.RainSpeed,
		PulseSpeed:             c.PulseSpeed,
		PulseMode:              c.PulseMode,
		RainEnabled:            c.RainEnabled,
		Debug:                  c.Debug,
		MusicReactive:          c.MusicReactive,
		MusicReactiveIntensity: c.MusicReactiveIntensity,
		RainPulse:              c.RainPulse,
		RainPulseEnabled:       c.RainPulseEnabled,
		SeparateFrequencies:    c.SeparateFrequencies,
		RainVisualizer:         c.RainVisualizer,
		CoverAnimation:         c.CoverAnimation,
		MusicPlayerAnimation:   c.MusicPlayerAnimation,
		MusicPlayerIntensity:   c.MusicPlayerIntensity,
		MusicPlayerInvert:      c.MusicPlayerInvert,
		RainInFrontOfPlayer:    c.RainInFrontOfPlayer,
		Direction:              c.Direction,
		Player:                 c.Player,
		RainLifetime:           c.RainLifetime,
	}
}
