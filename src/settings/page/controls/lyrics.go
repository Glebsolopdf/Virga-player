package controls

import "virga-player/settings"

func Lyrics(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.LyricsVisible
		cfg.LyricsVisible = !cfg.LyricsVisible
		return cfg.LyricsVisible != prev
	case 1:
		prev := cfg.LyricsMode
		options := settings.LyricsModeOptions()
		index := LyricsModeIndex(cfg)
		cfg.LyricsMode = options[(index+delta+len(options))%len(options)]
		return cfg.LyricsMode != prev
	case 2:
		prev := cfg.LyricsDisplayMode
		options := settings.LyricsDisplayModeOptions()
		index := LyricsDisplayIndex(cfg.LyricsDisplayMode)
		cfg.LyricsDisplayMode = options[(index+delta+len(options))%len(options)]
		return cfg.LyricsDisplayMode != prev
	case 3:
		prev := cfg.LyricsPosition
		cfg.LyricsPosition = nextLyricsPosition(cfg.LyricsPosition, cfg.PlayerPosition, delta)
		return cfg.LyricsPosition != prev
	case 4:
		prev := cfg.LyricsAutoSaveAfterSec
		cfg.LyricsAutoSaveAfterSec = clampInt(cfg.LyricsAutoSaveAfterSec+delta*5, 5, 600)
		return cfg.LyricsAutoSaveAfterSec != prev
	case 5:
		prev := cfg.LyricsRainLayer
		options := settings.RainLayerOptions()
		index := RainLayerIndex(cfg.LyricsRainLayer)
		cfg.LyricsRainLayer = options[(index+delta+len(options))%len(options)]
		return cfg.LyricsRainLayer != prev
	case 6:
		prev := cfg.LyricsDoubleConfirm
		cfg.LyricsDoubleConfirm = !cfg.LyricsDoubleConfirm
		return cfg.LyricsDoubleConfirm != prev
	}
	return false
}

func nextLyricsPosition(current, player settings.PositionPreset, delta int) settings.PositionPreset {
	options := settings.PositionPresetOptions()
	index := PositionIndex(current)
	for range options {
		index = (index + delta + len(options)) % len(options)
		if settings.IsLyricsPositionAllowed(player, options[index]) {
			return options[index]
		}
	}
	return current
}
