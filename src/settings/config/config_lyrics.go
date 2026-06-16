package config

func (c *Config) normalizeLyricsMode() {
	switch c.LyricsMode {
	case LyricsModeDisabled, LyricsModeRAMOnly, LyricsModeRAMWithAuto, LyricsModeRAMWithPrompt, LyricsModeDirectToDisk:
		return
	case LyricsModeLocal:
		c.LyricsMode = LyricsModeRAMOnly
		return
	case LyricsModeAuto:
		c.LyricsMode = LyricsModeRAMWithAuto
		return
	default:
		c.LyricsMode = LyricsModeDisabled
	}
}

func (m LyricsMode) Label() string {
	switch m {
	case LyricsModeDisabled:
		return "off"
	case LyricsModeRAMOnly:
		return "RAM only"
	case LyricsModeRAMWithAuto:
		return "RAM + auto-save"
	case LyricsModeRAMWithPrompt:
		return "RAM + save prompt"
	case LyricsModeDirectToDisk:
		return "direct to disk"
	default:
		return string(m)
	}
}

func LyricsModeOptions() []LyricsMode {
	return []LyricsMode{
		LyricsModeDisabled,
		LyricsModeRAMOnly,
		LyricsModeRAMWithAuto,
		LyricsModeRAMWithPrompt,
		LyricsModeDirectToDisk,
	}
}
