package page

import "virga-player/settings"

func generalAdjust(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.FPS
		cfg.FPS = max(15, min(cfg.FPS+delta*5, 240))
		return cfg.FPS != prev
	case 1:
		prev := cfg.MaxParticles
		cfg.MaxParticles = max(20, min(cfg.MaxParticles+delta*10, 500))
		return cfg.MaxParticles != prev
	case 2:
		prev := cfg.Debug
		cfg.Debug = !cfg.Debug
		return cfg.Debug != prev
	}
	return false
}

func rainAdjust(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.RainSpeed
		cfg.RainSpeed = max(25, min(cfg.RainSpeed+delta*5, 300))
		return cfg.RainSpeed != prev
	case 1:
		prev := cfg.RainLifetime
		cfg.RainLifetime = max(20, min(cfg.RainLifetime+delta*10, 200))
		return cfg.RainLifetime != prev
	case 2:
		prev := cfg.Direction
		options := settings.DirectionOptions()
		index := directionIndex(cfg)
		index = (index + delta + len(options)) % len(options)
		cfg.Direction = options[index]
		return cfg.Direction != prev
	case 3:
		prev := cfg.PulseSpeed
		cfg.PulseSpeed = max(25, min(cfg.PulseSpeed+delta*10, 300))
		return cfg.PulseSpeed != prev
	case 4:
		prev := cfg.PulseMode
		options := settings.PulseModeOptions()
		index := pulseModeIndex(cfg)
		index = (index + delta + len(options)) % len(options)
		cfg.SetPulseMode(options[index])
		return cfg.PulseMode != prev
	case 5:
		prev := cfg.RainEnabled
		cfg.RainEnabled = !cfg.RainEnabled
		return cfg.RainEnabled != prev
	case 6:
		prev := cfg.RainPulse
		cfg.RainPulse = max(20, min(cfg.RainPulse+delta*10, 200))
		return cfg.RainPulse != prev
	}
	return false
}

func audioAdjust(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.MusicReactive
		cfg.MusicReactive = !cfg.MusicReactive
		return cfg.MusicReactive != prev
	case 1:
		prev := cfg.MusicReactiveIntensity
		cfg.MusicReactiveIntensity = max(20, min(cfg.MusicReactiveIntensity+delta*10, 200))
		return cfg.MusicReactiveIntensity != prev
	case 2:
		prev := cfg.SeparateFrequencies
		cfg.SeparateFrequencies = !cfg.SeparateFrequencies
		return cfg.SeparateFrequencies != prev
	case 3:
		prev := cfg.RainVisualizer
		cfg.RainVisualizer = !cfg.RainVisualizer
		return cfg.RainVisualizer != prev
	}
	return false
}

func visualAdjust(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.MusicPlayerAnimation
		cfg.MusicPlayerAnimation = !cfg.MusicPlayerAnimation
		return cfg.MusicPlayerAnimation != prev
	case 1:
		prev := cfg.MusicPlayerIntensity
		cfg.MusicPlayerIntensity = max(20, min(cfg.MusicPlayerIntensity+delta*10, 200))
		return cfg.MusicPlayerIntensity != prev
	case 2:
		prev := cfg.MusicPlayerInvert
		cfg.MusicPlayerInvert = !cfg.MusicPlayerInvert
		return cfg.MusicPlayerInvert != prev
	case 3:
		prev := cfg.PlayerRainLayer
		options := settings.RainLayerOptions()
		index := rainLayerIndex(cfg.PlayerRainLayer)
		index = (index + delta + len(options)) % len(options)
		cfg.PlayerRainLayer = options[index]
		return cfg.PlayerRainLayer != prev
	case 4:
		prev := cfg.Player
		cfg.Player = !cfg.Player
		return cfg.Player != prev
	}
	return false
}

func lyricsAdjust(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.LyricsVisible
		cfg.LyricsVisible = !cfg.LyricsVisible
		return cfg.LyricsVisible != prev
	case 1:
		prev := cfg.LyricsMode
		options := settings.LyricsModeOptions()
		index := lyricsModeIndex(cfg)
		index = (index + delta + len(options)) % len(options)
		cfg.LyricsMode = options[index]
		return cfg.LyricsMode != prev
	case 2:
		prev := cfg.LyricsAutoSaveAfterSec
		cfg.LyricsAutoSaveAfterSec = max(5, min(cfg.LyricsAutoSaveAfterSec+delta*5, 600))
		return cfg.LyricsAutoSaveAfterSec != prev
	case 3:
		prev := cfg.LyricsRainLayer
		options := settings.RainLayerOptions()
		index := rainLayerIndex(cfg.LyricsRainLayer)
		index = (index + delta + len(options)) % len(options)
		cfg.LyricsRainLayer = options[index]
		return cfg.LyricsRainLayer != prev
	case 4:
		prev := cfg.LyricsDoubleConfirm
		cfg.LyricsDoubleConfirm = !cfg.LyricsDoubleConfirm
		return cfg.LyricsDoubleConfirm != prev
	}
	return false
}
