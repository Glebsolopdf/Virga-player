package config

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
	return []RainLayerMode{
		RainLayerBehind,
		RainLayerBetween,
		RainLayerFront,
	}
}

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
