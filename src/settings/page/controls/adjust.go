package controls

import "virga-player/settings"

func General(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.FPS
		cfg.FPS = clampInt(cfg.FPS+delta*5, 15, 240)
		return cfg.FPS != prev
	case 1:
		prev := cfg.MaxParticles
		cfg.MaxParticles = clampInt(cfg.MaxParticles+delta*10, 20, 500)
		return cfg.MaxParticles != prev
	case 2:
		prev := cfg.Debug
		cfg.Debug = !cfg.Debug
		return cfg.Debug != prev
	}
	return false
}

func Rain(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.RainSpeed
		cfg.RainSpeed = clampInt(cfg.RainSpeed+delta*5, 25, 300)
		return cfg.RainSpeed != prev
	case 1:
		prev := cfg.RainLifetime
		cfg.RainLifetime = clampInt(cfg.RainLifetime+delta*10, 20, 200)
		return cfg.RainLifetime != prev
	case 2:
		prev := cfg.PulseSpeed
		cfg.PulseSpeed = clampInt(cfg.PulseSpeed+delta*10, 25, 300)
		return cfg.PulseSpeed != prev
	case 3:
		prev := cfg.PulseMode
		options := settings.PulseModeOptions()
		index := PulseModeIndex(cfg)
		index = (index + delta + len(options)) % len(options)
		cfg.SetPulseMode(options[index])
		return cfg.PulseMode != prev
	case 4:
		prev := cfg.RainEnabled
		cfg.RainEnabled = !cfg.RainEnabled
		return cfg.RainEnabled != prev
	case 5:
		prev := cfg.RainPulse
		cfg.RainPulse = clampInt(cfg.RainPulse+delta*10, 20, 200)
		return cfg.RainPulse != prev
	case 6:
		prev := cfg.SeparateFrequencies
		cfg.SeparateFrequencies = !cfg.SeparateFrequencies
		return cfg.SeparateFrequencies != prev
	}
	return false
}

func Audio(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.MusicReactive
		cfg.MusicReactive = !cfg.MusicReactive
		return cfg.MusicReactive != prev
	case 1:
		prev := cfg.MusicReactiveIntensity
		cfg.MusicReactiveIntensity = clampInt(cfg.MusicReactiveIntensity+delta*10, 20, 200)
		return cfg.MusicReactiveIntensity != prev
	case 2:
		prev := cfg.RainVisualizer
		cfg.RainVisualizer = !cfg.RainVisualizer
		return cfg.RainVisualizer != prev
	}
	return false
}

func Visual(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.MusicPlayerAnimation
		cfg.MusicPlayerAnimation = !cfg.MusicPlayerAnimation
		return cfg.MusicPlayerAnimation != prev
	case 1:
		prev := cfg.MusicPlayerIntensity
		cfg.MusicPlayerIntensity = clampInt(cfg.MusicPlayerIntensity+delta*10, 20, 200)
		return cfg.MusicPlayerIntensity != prev
	case 2:
		prev := cfg.MusicPlayerInvert
		cfg.MusicPlayerInvert = !cfg.MusicPlayerInvert
		return cfg.MusicPlayerInvert != prev
	case 3:
		prev := cfg.PlayerRainLayer
		options := settings.RainLayerOptions()
		index := RainLayerIndex(cfg.PlayerRainLayer)
		index = (index + delta + len(options)) % len(options)
		cfg.PlayerRainLayer = options[index]
		return cfg.PlayerRainLayer != prev
	case 4:
		prev := cfg.Direction
		options := settings.DirectionOptions()
		index := DirectionIndex(cfg)
		index = (index + delta + len(options)) % len(options)
		cfg.Direction = options[index]
		return cfg.Direction != prev
	case 5:
		prev := cfg.Player
		cfg.Player = !cfg.Player
		return cfg.Player != prev
	}
	return false
}

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
		index = (index + delta + len(options)) % len(options)
		cfg.LyricsMode = options[index]
		return cfg.LyricsMode != prev
	case 2:
		prev := cfg.LyricsAutoSaveAfterSec
		cfg.LyricsAutoSaveAfterSec = clampInt(cfg.LyricsAutoSaveAfterSec+delta*5, 5, 600)
		return cfg.LyricsAutoSaveAfterSec != prev
	case 3:
		prev := cfg.LyricsRainLayer
		options := settings.RainLayerOptions()
		index := RainLayerIndex(cfg.LyricsRainLayer)
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

func Notifications(cfg *settings.Config, selected, delta int) bool {
	switch selected {
	case 0:
		prev := cfg.NotificationsEnabled
		cfg.NotificationsEnabled = !cfg.NotificationsEnabled
		return cfg.NotificationsEnabled != prev
	case 1:
		prev := cfg.NotifyUnreadToast
		cfg.NotifyUnreadToast = !cfg.NotifyUnreadToast
		return cfg.NotifyUnreadToast != prev
	}
	return false
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
