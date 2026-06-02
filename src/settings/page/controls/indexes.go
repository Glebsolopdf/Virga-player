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

func LyricsModeIndex(cfg *settings.Config) int {
	options := settings.LyricsModeOptions()
	for i, mode := range options {
		if mode == cfg.LyricsMode {
			return i
		}
	}
	return 0
}

func RainLayerIndex(current settings.RainLayerMode) int {
	options := settings.RainLayerOptions()
	for i, mode := range options {
		if mode == current {
			return i
		}
	}
	return 0
}
