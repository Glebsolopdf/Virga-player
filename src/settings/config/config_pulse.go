package config

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
	return []DirectionMode{
		DirectionRightToLeft,
		DirectionLeftToRight,
		DirectionStraight,
		DirectionRandom,
	}
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
	return []PulseMode{
		PulseModeOff,
		PulseModeRain,
		PulseModeCover,
		PulseModeAll,
	}
}
