package config

import (
	"path/filepath"

	"virga-player/lyricsearch"
)

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
	if c.RainLifetime < 20 {
		c.RainLifetime = 20
	}
	if c.RainLifetime > 200 {
		c.RainLifetime = 200
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
	c.normalizeLyricsMode()
	if c.LyricsAutoSaveAfterSec <= 0 {
		c.LyricsAutoSaveAfterSec = 30
	}
	if c.LyricsAutoSaveAfterSec > 600 {
		c.LyricsAutoSaveAfterSec = 600
	}
	if c.LyricsTempDir == "" {
		c.LyricsTempDir = filepath.Join("/tmp", "virgaplayerlyrics")
	}
	if c.LyricsPersistentDir == "" {
		c.LyricsPersistentDir = lyricsearch.DefaultPersistentDir()
	}
	c.normalizeRainLayerModes()
	switch c.Direction {
	case DirectionRightToLeft, DirectionLeftToRight, DirectionStraight, DirectionRandom:
		return
	default:
		c.Direction = DirectionRandom
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
