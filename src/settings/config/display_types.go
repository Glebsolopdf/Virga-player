package config

type PositionPreset string
type LyricsDisplayMode string

const (
	PositionTopRight    PositionPreset = "top-right"
	PositionBottomRight PositionPreset = "bottom-right"
	PositionTopCenter   PositionPreset = "top-center"
	PositionCenter      PositionPreset = "center"
	PositionBelowPlayer PositionPreset = "below-player"

	LyricsDisplayStatic  LyricsDisplayMode = "static"
	LyricsDisplayDynamic LyricsDisplayMode = "dynamic"
)

func PositionPresetOptions() []PositionPreset {
	return []PositionPreset{
		PositionTopRight,
		PositionBottomRight,
		PositionTopCenter,
		PositionCenter,
		PositionBelowPlayer,
	}
}

func LyricsDisplayModeOptions() []LyricsDisplayMode {
	return []LyricsDisplayMode{LyricsDisplayStatic, LyricsDisplayDynamic}
}

func (p PositionPreset) Label() string {
	switch p {
	case PositionTopRight:
		return "top-right"
	case PositionBottomRight:
		return "bottom-right"
	case PositionTopCenter:
		return "top-center"
	case PositionCenter:
		return "center"
	case PositionBelowPlayer:
		return "below-player"
	default:
		return string(p)
	}
}

func (m LyricsDisplayMode) Label() string {
	switch m {
	case LyricsDisplayStatic:
		return "static"
	case LyricsDisplayDynamic:
		return "dynamic"
	default:
		return string(m)
	}
}
