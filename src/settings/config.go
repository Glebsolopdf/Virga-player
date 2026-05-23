package settings

import cfg "virga-player/settings/config"

type Config = cfg.Config
type DirectionMode = cfg.DirectionMode
type PulseMode = cfg.PulseMode

const (
	DirectionRightToLeft = cfg.DirectionRightToLeft
	DirectionLeftToRight = cfg.DirectionLeftToRight
	DirectionStraight    = cfg.DirectionStraight
	DirectionRandom      = cfg.DirectionRandom

	PulseModeOff   = cfg.PulseModeOff
	PulseModeRain  = cfg.PulseModeRain
	PulseModeCover = cfg.PulseModeCover
	PulseModeAll   = cfg.PulseModeAll
)

func DefaultConfig() *Config {
	return cfg.DefaultConfig()
}

func LoadOrCreateConfig() (*Config, bool, error) {
	return cfg.LoadOrCreateConfig()
}

func LoadConfig() (*Config, error) {
	return cfg.LoadConfig()
}

func SaveConfig(cfgObj *Config) error {
	return cfg.SaveConfig(cfgObj)
}

func ConfigPath() string {
	return cfg.ConfigPath()
}

func ConfigDirPath() string {
	return cfg.ConfigDirPath()
}

func DirectionOptions() []DirectionMode {
	return cfg.DirectionOptions()
}

func PulseModeOptions() []PulseMode {
	return cfg.PulseModeOptions()
}
