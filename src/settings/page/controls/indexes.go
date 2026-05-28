package controls

import "virga-player/settings"

func PulseModeIndex(cfg *settings.Config) int {
	options := settings.PulseModeOptions()
	for i, mode := range options {
		if mode == cfg.PulseMode {
			return i
		}
	}
	return 0
}

func DirectionIndex(cfg *settings.Config) int {
	options := settings.DirectionOptions()
	for i, mode := range options {
		if mode == cfg.Direction {
			return i
		}
	}
	return 0
}
