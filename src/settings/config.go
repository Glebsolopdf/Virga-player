package settings

import cfg "virga-player/settings/config"

type Config = cfg.Config
type DirectionMode = cfg.DirectionMode
type LyricsMode = cfg.LyricsMode
type PulseMode = cfg.PulseMode
type RainLayerMode = cfg.RainLayerMode

const (
	DirectionRightToLeft = cfg.DirectionRightToLeft
	DirectionLeftToRight = cfg.DirectionLeftToRight
	DirectionStraight    = cfg.DirectionStraight
	DirectionRandom      = cfg.DirectionRandom

	LyricsModeDisabled      = cfg.LyricsModeDisabled
	LyricsModeRAMOnly       = cfg.LyricsModeRAMOnly
	LyricsModeRAMWithAuto   = cfg.LyricsModeRAMWithAuto
	LyricsModeRAMWithPrompt = cfg.LyricsModeRAMWithPrompt
	LyricsModeDirectToDisk  = cfg.LyricsModeDirectToDisk

	// Legacy aliases kept to preserve compatibility with older code paths.
	LyricsModeLocal = cfg.LyricsModeLocal
	LyricsModeAuto  = cfg.LyricsModeAuto

	PulseModeOff   = cfg.PulseModeOff
	PulseModeRain  = cfg.PulseModeRain
	PulseModeCover = cfg.PulseModeCover
	PulseModeAll   = cfg.PulseModeAll

	RainLayerBehind  = cfg.RainLayerBehind
	RainLayerBetween = cfg.RainLayerBetween
	RainLayerFront   = cfg.RainLayerFront
)

func DefaultConfig() *Config {
	return cfg.DefaultConfig()
}

func LoadOrCreateConfig() (*Config, bool, error) {
	return cfg.LoadOrCreateConfig()
}

func SaveConfig(cfgObj *Config) error {
	return cfg.SaveConfig(cfgObj)
}

func ConfigDirPath() string {
	return cfg.ConfigDirPath()
}

func DirectionOptions() []DirectionMode {
	return cfg.DirectionOptions()
}

func LyricsModeOptions() []LyricsMode {
	return cfg.LyricsModeOptions()
}

func PulseModeOptions() []PulseMode {
	return cfg.PulseModeOptions()
}

func RainLayerOptions() []RainLayerMode {
	return cfg.RainLayerOptions()
}
