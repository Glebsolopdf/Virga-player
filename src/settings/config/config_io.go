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

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), false, err
	}
	cfg.normalize()
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
		RainInFrontOfPlayer:    c.RainInFrontOfPlayer,
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
	c.RainInFrontOfPlayer = raw.RainInFrontOfPlayer
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
