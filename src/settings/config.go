package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DirectionMode string
type PulseMode string

const (
	DirectionRightToLeft DirectionMode = "right-to-left"
	DirectionLeftToRight DirectionMode = "left-to-right"
	DirectionStraight    DirectionMode = "straight"
	DirectionRandom      DirectionMode = "random"

	PulseModeOff   PulseMode = "off"
	PulseModeRain  PulseMode = "rain"
	PulseModeCover PulseMode = "cover"
	PulseModeAll   PulseMode = "all"
)

type Config struct {
	FPS                    int           `json:"fps"`
	MaxParticles           int           `json:"max_particles"`
	RainSpeed              int           `json:"rain_speed"`
	PulseSpeed             int           `json:"pulse_speed"`
	PulseMode              PulseMode     `json:"pulse_mode"`
	RainEnabled            bool          `json:"rain_enabled"`
	RainPulse              int           `json:"rain_pulse"`
	RainPulseEnabled       bool          `json:"-"`
	SeparateFrequencies    bool          `json:"separate_frequencies"`
	Debug                  bool          `json:"debug"`
	MusicReactive          bool          `json:"music_reactive"`
	MusicReactiveIntensity int           `json:"music_reactive_intensity"`
	RainVisualizer         bool          `json:"rain_visualizer"`
	CoverAnimation         bool          `json:"-"`
	MusicPlayerAnimation   bool          `json:"music_player_animation"`
	MusicPlayerIntensity   int           `json:"music_player_intensity"`
	MusicPlayerInvert      bool          `json:"music_player_invert"`
	RainInFrontOfPlayer    bool          `json:"rain_in_front_of_player"`
	Direction              DirectionMode `json:"direction"`
	Player                 bool          `json:"player"`
}

type savedConfigJSON struct {
	FPS                    int           `json:"fps"`
	MaxParticles           int           `json:"max_particles"`
	RainSpeed              int           `json:"rain_speed"`
	PulseSpeed             int           `json:"pulse_speed"`
	PulseMode              PulseMode     `json:"pulse_mode"`
	RainEnabled            bool          `json:"rain_enabled"`
	RainPulse              int           `json:"rain_pulse"`
	SeparateFrequencies    bool          `json:"separate_frequencies"`
	Debug                  bool          `json:"debug"`
	MusicReactive          bool          `json:"music_reactive"`
	MusicReactiveIntensity int           `json:"music_reactive_intensity"`
	RainVisualizer         bool          `json:"rain_visualizer"`
	MusicPlayerAnimation   bool          `json:"music_player_animation"`
	MusicPlayerIntensity   int           `json:"music_player_intensity"`
	MusicPlayerInvert      bool          `json:"music_player_invert"`
	RainInFrontOfPlayer    bool          `json:"rain_in_front_of_player"`
	Direction              DirectionMode `json:"direction"`
	Player                 bool          `json:"player"`
}

type loadedConfigJSON struct {
	savedConfigJSON
	RainPulseEnabled bool `json:"rain_pulse_enabled"`
	CoverAnimation   bool `json:"cover_animation"`
}

func DefaultConfig() *Config {
	return &Config{
		FPS:                    60,
		MaxParticles:           220,
		RainSpeed:              100,
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
		RainInFrontOfPlayer:    true,
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

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(savedConfigJSON{
		FPS:                    c.FPS,
		MaxParticles:           c.MaxParticles,
		RainSpeed:              c.RainSpeed,
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
	if c.PulseSpeed < 25 {
		c.PulseSpeed = 25
	}
	if c.PulseSpeed > 300 {
		c.PulseSpeed = 300
	}
	c.normalizePulseMode()
	if c.RainPulse < 20 {
		c.RainPulse = 20
	}
	if c.RainPulse > 200 {
		c.RainPulse = 200
	}
	if c.MusicReactiveIntensity < 20 {
		c.MusicReactiveIntensity = 20
	}
	if c.MusicReactiveIntensity > 200 {
		c.MusicReactiveIntensity = 200
	}
	if c.MusicPlayerIntensity < 20 {
		c.MusicPlayerIntensity = 20
	}
	if c.MusicPlayerIntensity > 200 {
		c.MusicPlayerIntensity = 200
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
		PulseSpeed:             c.PulseSpeed,
		PulseMode:              c.PulseMode,
		RainEnabled:            c.RainEnabled,
		Debug:                  c.Debug,
		MusicReactive:          c.MusicReactive,
		MusicReactiveIntensity: c.MusicReactiveIntensity,
		RainPulse:              c.RainPulse,
		RainPulseEnabled:       c.RainPulseEnabled,
		SeparateFrequencies:    c.SeparateFrequencies,
		RainVisualizer:         c.RainVisualizer,
		CoverAnimation:         c.CoverAnimation,
		MusicPlayerAnimation:   c.MusicPlayerAnimation,
		MusicPlayerIntensity:   c.MusicPlayerIntensity,
		MusicPlayerInvert:      c.MusicPlayerInvert,
		RainInFrontOfPlayer:    c.RainInFrontOfPlayer,
		Direction:              c.Direction,
		Player:                 c.Player,
	}
}

func inferPulseMode(rainEnabled, coverEnabled bool) PulseMode {
	switch {
	case rainEnabled && coverEnabled:
		return PulseModeAll
	case rainEnabled:
		return PulseModeRain
	case coverEnabled:
		return PulseModeCover
	default:
		return PulseModeOff
	}
}

func (c *Config) normalizePulseMode() {
	if c.PulseMode == "" {
		c.PulseMode = inferPulseMode(c.RainPulseEnabled, c.CoverAnimation)
	}
	switch c.PulseMode {
	case PulseModeOff, PulseModeRain, PulseModeCover, PulseModeAll:
	default:
		c.PulseMode = inferPulseMode(c.RainPulseEnabled, c.CoverAnimation)
	}
	c.syncLegacyPulseFlags()
}

func (c *Config) syncLegacyPulseFlags() {
	c.RainPulseEnabled = c.PulseOnRain()
	c.CoverAnimation = c.PulseOnCover()
}

func (c *Config) SetPulseMode(mode PulseMode) {
	c.PulseMode = mode
	c.normalizePulseMode()
}

func (c *Config) PulseOnRain() bool {
	return c.PulseMode == PulseModeRain || c.PulseMode == PulseModeAll
}

func (c *Config) PulseOnCover() bool {
	return c.PulseMode == PulseModeCover || c.PulseMode == PulseModeAll
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

func (p PulseMode) Label() string {
	switch p {
	case PulseModeOff:
		return "off"
	case PulseModeRain:
		return "rain only"
	case PulseModeCover:
		return "cover only"
	case PulseModeAll:
		return "rain and cover"
	default:
		return string(p)
	}
}

func PulseModeOptions() []PulseMode {
	return []PulseMode{
		PulseModeOff,
		PulseModeRain,
		PulseModeCover,
		PulseModeAll,
	}
}
