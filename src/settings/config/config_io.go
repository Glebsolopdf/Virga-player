package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadOrCreateConfig() (*Config, bool, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			return cfg, true, SaveConfig(cfg)
		}
		return DefaultConfig(), false, err
	}

	missingLyricsFields := true
	missingRainLayerFields := true
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		_, hasLyricsMode := raw["lyrics_mode"]
		_, hasLyricsVisible := raw["lyrics_visible"]
		_, hasLyricsAutoSaveAfterSec := raw["lyrics_auto_save_after_sec"]
		_, hasLyricsDoubleConfirm := raw["lyrics_double_confirm"]
		missingLyricsFields = !hasLyricsMode || !hasLyricsVisible || !hasLyricsAutoSaveAfterSec || !hasLyricsDoubleConfirm
		_, hasPlayerRainLayer := raw["player_rain_layer"]
		_, hasLyricsRainLayer := raw["lyrics_rain_layer"]
		missingRainLayerFields = !hasPlayerRainLayer || !hasLyricsRainLayer

	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), false, err
	}
	cfg.normalize()
	if missingLyricsFields || missingRainLayerFields {
		if err := SaveConfig(cfg); err != nil {
			return cfg, false, err
		}
	}
	return cfg, false, nil
}

func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(savedConfigJSON{
		FPS:                    c.FPS,
		MaxParticles:           c.MaxParticles,
		RainSpeed:              c.RainSpeed,
		RainLifetime:           c.RainLifetime,
		PulseSpeed:             c.PulseSpeed,
		PulseMode:              c.PulseMode,
		RainEnabled:            c.RainEnabled,
		RainPulse:              c.RainPulse,
		SeparateFrequencies:    c.SeparateFrequencies,
		Debug:                  c.Debug,
		MusicReactive:          c.MusicReactive,
		MusicReactiveIntensity: c.MusicReactiveIntensity,
		RainVisualizer:         c.RainVisualizer,
		MusicPlayerAnimation:   c.MusicPlayerAnimation,
		MusicPlayerIntensity:   c.MusicPlayerIntensity,
		MusicPlayerInvert:      c.MusicPlayerInvert,
		PlayerRainLayer:        c.PlayerRainLayer,
		LyricsMode:             c.LyricsMode,
		LyricsVisible:          c.LyricsVisible,
		LyricsRainLayer:        c.LyricsRainLayer,
		LyricsSaveToCache:      c.LyricsSaveToCache,
		LyricsAutoSaveAfterSec: c.LyricsAutoSaveAfterSec,
		LyricsDoubleConfirm:    c.LyricsDoubleConfirm,
		LyricsTempDir:          c.LyricsTempDir,
		LyricsPersistentDir:    c.LyricsPersistentDir,
		Direction:              c.Direction,
		Player:                 c.Player,
	})
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var raw loadedConfigJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.FPS = raw.FPS
	c.MaxParticles = raw.MaxParticles
	c.RainSpeed = raw.RainSpeed
	c.RainLifetime = raw.RainLifetime
	c.PulseSpeed = raw.PulseSpeed
	c.PulseMode = raw.PulseMode
	c.RainEnabled = raw.RainEnabled
	c.RainPulse = raw.RainPulse
	c.RainPulseEnabled = raw.RainPulseEnabled
	c.SeparateFrequencies = raw.SeparateFrequencies
	c.Debug = raw.Debug
	c.MusicReactive = raw.MusicReactive
	c.MusicReactiveIntensity = raw.MusicReactiveIntensity
	c.RainVisualizer = raw.RainVisualizer
	c.CoverAnimation = raw.CoverAnimation
	c.MusicPlayerAnimation = raw.MusicPlayerAnimation
	c.MusicPlayerIntensity = raw.MusicPlayerIntensity
	c.MusicPlayerInvert = raw.MusicPlayerInvert
	c.PlayerRainLayer = raw.PlayerRainLayer
	if c.PlayerRainLayer == "" && raw.RainInFrontOfPlayer != nil {
		if *raw.RainInFrontOfPlayer {
			c.PlayerRainLayer = RainLayerBehind
		} else {
			c.PlayerRainLayer = RainLayerFront
		}
	}
	c.LyricsMode = raw.LyricsMode
	c.LyricsVisible = raw.LyricsVisible
	c.LyricsRainLayer = raw.LyricsRainLayer
	c.LyricsSaveToCache = raw.LyricsSaveToCache
	c.LyricsAutoSaveAfterSec = raw.LyricsAutoSaveAfterSec
	if raw.LyricsDoubleConfirm == nil {
		c.LyricsDoubleConfirm = true
	} else {
		c.LyricsDoubleConfirm = *raw.LyricsDoubleConfirm
	}
	c.LyricsTempDir = raw.LyricsTempDir
	c.LyricsPersistentDir = raw.LyricsPersistentDir
	c.Direction = raw.Direction
	c.Player = raw.Player

	return nil
}

func ConfigPath() string {
	return filepath.Join(ConfigDirPath(), "config.json")
}

func ConfigDirPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join("/tmp", "virga-player")
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "virga-player")
}
