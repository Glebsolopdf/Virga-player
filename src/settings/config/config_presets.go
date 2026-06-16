package config

func PlayerConfig() *Config {
	cfg := DefaultConfig()
	cfg.FPS = 70
	cfg.MaxParticles = 240
	cfg.RainSpeed = 100
	cfg.RainLifetime = 100
	cfg.PulseSpeed = 100
	cfg.SetPulseMode(PulseModeRain)
	cfg.RainEnabled = true
	cfg.RainPulse = 100
	cfg.SeparateFrequencies = true
	cfg.Player = true
	cfg.MusicReactive = true
	cfg.MusicReactiveIntensity = 100
	cfg.RainVisualizer = false
	cfg.MusicPlayerAnimation = false
	cfg.MusicPlayerIntensity = 100
	cfg.MusicPlayerInvert = false
	cfg.PlayerPosition = PositionCenter
	cfg.PlayerRainLayer = RainLayerBetween
	cfg.LyricsMode = LyricsModeRAMWithPrompt
	cfg.LyricsDisplayMode = LyricsDisplayStatic
	cfg.LyricsPosition = PositionBelowPlayer
	cfg.LyricsVisible = true
	cfg.LyricsRainLayer = RainLayerBetween
	cfg.LyricsSaveToCache = true
	cfg.LyricsAutoSaveAfterSec = 15
	cfg.LyricsDoubleConfirm = true
	cfg.Direction = DirectionStraight
	return cfg
}
