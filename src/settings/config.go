package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DirectionMode string

const (
	DirectionRightToLeft DirectionMode = "right-to-left"
	DirectionLeftToRight DirectionMode = "left-to-right"
	DirectionStraight    DirectionMode = "straight"
	DirectionRandom      DirectionMode = "random"
)

type Config struct {
	FPS                    int           `json:"fps"`
	MaxParticles           int           `json:"max_particles"`
	RainSpeed              int           `json:"rain_speed"`
	RainEnabled            bool          `json:"rain_enabled"`
	MusicReactive          bool          `json:"music_reactive"`
	MusicReactiveIntensity int           `json:"music_reactive_intensity"`
	CoverAnimation         bool          `json:"cover_animation"`
	Direction              DirectionMode `json:"direction"`
	Player                 bool          `json:"player"`
}

func DefaultConfig() *Config {
	return &Config{
		FPS:                    60,
		MaxParticles:           220,
		RainSpeed:              100,
		RainEnabled:            true,
		MusicReactive:          false,
		MusicReactiveIntensity: 100,
		CoverAnimation:         false,
		Direction:              DirectionRandom,
		Player:                 false,
	}
}

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

func LoadConfig() (*Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), err
	}
	cfg.normalize()
	return cfg, nil
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

func (c *Config) normalize() {
	if c.FPS <= 0 {
		c.FPS = 60
	}
	if c.FPS > 240 {
		c.FPS = 240
	}
	if c.MaxParticles < 20 {
		c.MaxParticles = 20
	}
	if c.MaxParticles > 500 {
		c.MaxParticles = 500
	}
	if c.RainSpeed < 25 {
		c.RainSpeed = 25
	}
	if c.RainSpeed > 300 {
		c.RainSpeed = 300
	}
	if c.MusicReactiveIntensity < 20 {
		c.MusicReactiveIntensity = 20
	}
	if c.MusicReactiveIntensity > 200 {
		c.MusicReactiveIntensity = 200
	}
	switch c.Direction {
	case DirectionRightToLeft, DirectionLeftToRight, DirectionStraight, DirectionRandom:
		return
	default:
		c.Direction = DirectionRandom
	}
}

func (c *Config) Clone() *Config {
	return &Config{
		FPS:                    c.FPS,
		MaxParticles:           c.MaxParticles,
		RainSpeed:              c.RainSpeed,
		RainEnabled:            c.RainEnabled,
		MusicReactive:          c.MusicReactive,
		MusicReactiveIntensity: c.MusicReactiveIntensity,
		CoverAnimation:         c.CoverAnimation,
		Direction:              c.Direction,
		Player:                 c.Player,
	}
}

func (d DirectionMode) Label() string {
	switch d {
	case DirectionRightToLeft:
		return "right to left"
	case DirectionLeftToRight:
		return "left to right"
	case DirectionStraight:
		return "straight"
	case DirectionRandom:
		return "random"
	default:
		return string(d)
	}
}

func DirectionOptions() []DirectionMode {
	return []DirectionMode{
		DirectionRightToLeft,
		DirectionLeftToRight,
		DirectionStraight,
		DirectionRandom,
	}
}
