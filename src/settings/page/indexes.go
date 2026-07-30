package page

import "virga-player/settings"

func pulseModeIndex(cfg *settings.Config) int {
	options := settings.PulseModeOptions()
	for i, mode := range options {
		if mode == cfg.PulseMode {
			return i
		}
	}
	return 0
}

func directionIndex(cfg *settings.Config) int {
	options := settings.DirectionOptions()
	for i, mode := range options {
		if mode == cfg.Direction {
			return i
		}
	}
	return 0
}

func lyricsModeIndex(cfg *settings.Config) int {
	options := settings.LyricsModeOptions()
	for i, mode := range options {
		if mode == cfg.LyricsMode {
			return i
		}
	}
	return 0
}

func rainLayerIndex(current settings.RainLayerMode) int {
	options := settings.RainLayerOptions()
	for i, mode := range options {
		if mode == current {
			return i
		}
	}
	return 0
}
