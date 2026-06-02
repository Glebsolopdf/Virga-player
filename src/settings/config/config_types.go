package config

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

type loadedConfigJSON struct {
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
	PlayerRainLayer        RainLayerMode `json:"player_rain_layer"`
	RainInFrontOfPlayer    *bool         `json:"rain_in_front_of_player"`
	LyricsMode             LyricsMode    `json:"lyrics_mode"`
	LyricsVisible          bool          `json:"lyrics_visible"`
	LyricsRainLayer        RainLayerMode `json:"lyrics_rain_layer"`
	LyricsSaveToCache      bool          `json:"lyrics_save_to_cache"`
	LyricsAutoSaveAfterSec int           `json:"lyrics_auto_save_after_sec"`
	LyricsDoubleConfirm    *bool         `json:"lyrics_double_confirm"`
	LyricsTempDir          string        `json:"lyrics_temp_dir,omitempty"`
	LyricsPersistentDir    string        `json:"lyrics_persistent_dir,omitempty"`
	Direction              DirectionMode `json:"direction"`
	Player                 bool          `json:"player"`

	RainPulseEnabled bool `json:"rain_pulse_enabled"`
	CoverAnimation   bool `json:"cover_animation"`
}
