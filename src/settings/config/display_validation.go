package config

func IsLyricsPositionAllowed(player, lyrics PositionPreset) bool {
	if lyrics == PositionBelowPlayer {
		return true
	}
	return player != lyrics
}

func ResolveLyricsPosition(player, lyrics PositionPreset) PositionPreset {
	if IsLyricsPositionAllowed(player, lyrics) {
		return lyrics
	}
	return PositionBelowPlayer
}
