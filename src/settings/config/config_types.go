package config

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
	RainLifetime           int           `json:"rain_lifetime"`
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
	RainLifetime           int           `json:"rain_lifetime"`
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
