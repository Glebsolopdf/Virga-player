package config

func normalizeRainLayerMode(mode RainLayerMode, fallback RainLayerMode) RainLayerMode {
	switch mode {
	case RainLayerBehind, RainLayerBetween, RainLayerFront:
		return mode
	default:
		return fallback
	}
}

func (c *Config) normalizeRainLayerModes() {
	c.PlayerRainLayer = normalizeRainLayerMode(c.PlayerRainLayer, RainLayerBehind)
	c.LyricsRainLayer = normalizeRainLayerMode(c.LyricsRainLayer, RainLayerFront)
}

func (m RainLayerMode) Label() string {
	switch m {
	case RainLayerBehind:
		return "behind rain"
	case RainLayerBetween:
		return "between rain layers"
	case RainLayerFront:
		return "in front of rain"
	default:
		return string(m)
	}
}

func RainLayerOptions() []RainLayerMode {
	return []RainLayerMode{RainLayerBehind, RainLayerBetween, RainLayerFront}
}

func (c *Config) normalizeLyricsMode() {
	switch c.LyricsMode {
	case LyricsModeDisabled, LyricsModeRAMOnly, LyricsModeRAMWithAuto, LyricsModeRAMWithPrompt, LyricsModeDirectToDisk:
		return
	case LyricsModeLocal:
		c.LyricsMode = LyricsModeRAMOnly
	case LyricsModeAuto:
		c.LyricsMode = LyricsModeRAMWithAuto
	default:
		c.LyricsMode = LyricsModeDisabled
	}
}

func (m LyricsMode) Label() string {
	switch m {
	case LyricsModeDisabled:
		return "off"
	case LyricsModeRAMOnly:
		return "RAM only"
	case LyricsModeRAMWithAuto:
		return "RAM + auto-save"
	case LyricsModeRAMWithPrompt:
		return "RAM + save prompt"
	case LyricsModeDirectToDisk:
		return "direct to disk"
	default:
		return string(m)
	}
}

func LyricsModeOptions() []LyricsMode {
	return []LyricsMode{
		LyricsModeDisabled, LyricsModeRAMOnly, LyricsModeRAMWithAuto,
		LyricsModeRAMWithPrompt, LyricsModeDirectToDisk,
	}
}

func (c *Config) SetPulseMode(mode PulseMode) {
	c.PulseMode = mode
	c.normalizePulseMode()
}

func (c *Config) PulseOnRain() bool {
	return c.PulseMode == PulseModeRain || c.PulseMode == PulseModeAll
}

func (c *Config) PulseOnCover() bool {
	return c.PulseMode == PulseModeCover || c.PulseMode == PulseModeAll
}

func (d DirectionMode) Label() string {
	switch d {
	case DirectionRightToLeft:
		return "right to left"
	case DirectionLeftToRight:
		return "left to right"
	case DirectionStraight:
		return "straight"
	case DirectionRandom:
		return "random"
	default:
		return string(d)
	}
}

func DirectionOptions() []DirectionMode {
	return []DirectionMode{DirectionRightToLeft, DirectionLeftToRight, DirectionStraight, DirectionRandom}
}

func (p PulseMode) Label() string {
	switch p {
	case PulseModeOff:
		return "off"
	case PulseModeRain:
		return "rain only"
	case PulseModeCover:
		return "cover only"
	case PulseModeAll:
		return "rain and cover"
	default:
		return string(p)
	}
}

func PulseModeOptions() []PulseMode {
	return []PulseMode{PulseModeOff, PulseModeRain, PulseModeCover, PulseModeAll}
}
