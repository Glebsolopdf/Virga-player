package config

import (
	"path/filepath"

	"virga-player/lyricsearch"
)

type DirectionMode string
type LyricsMode string
type PulseMode string
type RainLayerMode string

const (
	DirectionRightToLeft DirectionMode = "right-to-left"
	DirectionLeftToRight DirectionMode = "left-to-right"
	DirectionStraight    DirectionMode = "straight"
	DirectionRandom      DirectionMode = "random"

	LyricsModeDisabled      LyricsMode = "disabled"
	LyricsModeRAMOnly       LyricsMode = "ram-only"
	LyricsModeRAMWithAuto   LyricsMode = "ram-with-auto-save"
	LyricsModeRAMWithPrompt LyricsMode = "ram-with-save-prompt"
	LyricsModeDirectToDisk  LyricsMode = "direct-to-disk"

	// Legacy values retained for config migration compatibility.
	LyricsModeLocal LyricsMode = "local"
	LyricsModeAuto  LyricsMode = "auto"

	PulseModeOff   PulseMode = "off"
	PulseModeRain  PulseMode = "rain"
	PulseModeCover PulseMode = "cover"
	PulseModeAll   PulseMode = "all"

	RainLayerBehind  RainLayerMode = "behind"
	RainLayerBetween RainLayerMode = "between"
	RainLayerFront   RainLayerMode = "front"
)

type Config struct {
	FPS                    int           `json:"fps"`
	MaxParticles           int           `json:"max_particles"`
	RainSpeed              int           `json:"rain_speed"`
	RainLifetime           int           `json:"rain_lifetime"`
	PulseSpeed             int           `json:"pulse_speed"`
	PulseMode              PulseMode     `json:"pulse_mode"`
	RainEnabled            bool          `json:"rain_enabled"`
	RainPulse              int           `json:"rain_pulse"`
	RainPulseEnabled       bool          `json:"rain_pulse_enabled"`
	SeparateFrequencies    bool          `json:"separate_frequencies"`
	Debug                  bool          `json:"debug"`
	MusicReactive          bool          `json:"music_reactive"`
	MusicReactiveIntensity int           `json:"music_reactive_intensity"`
	RainVisualizer         bool          `json:"rain_visualizer"`
	CoverAnimation         bool          `json:"cover_animation"`
	MusicPlayerAnimation   bool          `json:"music_player_animation"`
	MusicPlayerIntensity   int           `json:"music_player_intensity"`
	MusicPlayerInvert      bool          `json:"music_player_invert"`
	PlayerRainLayer        RainLayerMode `json:"player_rain_layer"`
	LyricsMode             LyricsMode    `json:"lyrics_mode"`
	LyricsVisible          bool          `json:"lyrics_visible"`
	LyricsRainLayer        RainLayerMode `json:"lyrics_rain_layer"`
	LyricsSaveToCache      bool          `json:"lyrics_save_to_cache"`
	LyricsAutoSaveAfterSec int           `json:"lyrics_auto_save_after_sec"`
	LyricsDoubleConfirm    bool          `json:"lyrics_double_confirm"`
	LyricsTempDir          string        `json:"lyrics_temp_dir,omitempty"`
	LyricsPersistentDir    string        `json:"lyrics_persistent_dir,omitempty"`
	Direction              DirectionMode `json:"direction"`
	Player                 bool          `json:"player"`
}

func DefaultConfig() *Config {
	return &Config{
		FPS:                    60,
		MaxParticles:           220,
		RainSpeed:              100,
		RainLifetime:           100,
		PulseSpeed:             100,
		PulseMode:              PulseModeRain,
		RainEnabled:            true,
		RainPulse:              100,
		RainPulseEnabled:       true,
		SeparateFrequencies:    false,
		Debug:                  false,
		MusicReactive:          false,
		MusicReactiveIntensity: 100,
		RainVisualizer:         false,
		CoverAnimation:         false,
		MusicPlayerAnimation:   false,
		MusicPlayerIntensity:   100,
		MusicPlayerInvert:      false,
		PlayerRainLayer:        RainLayerBehind,
		LyricsMode:             LyricsModeDisabled,
		LyricsVisible:          true,
		LyricsRainLayer:        RainLayerFront,
		LyricsSaveToCache:      true,
		LyricsAutoSaveAfterSec: 30,
		LyricsDoubleConfirm:    true,
		LyricsTempDir:          filepath.Join("/tmp", "virgaplayerlyrics"),
		LyricsPersistentDir:    lyricsearch.DefaultPersistentDir(),
		Direction:              DirectionRandom,
		Player:                 false,
	}
}

func (c *Config) Clone() *Config {
	clone := *c
	return &clone
}


