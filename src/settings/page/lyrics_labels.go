package page

import (
	"fmt"

	"virga-player/settings"
)

func lyricsWarning(cfg *settings.Config) string {
	if settings.IsLyricsPositionAllowed(cfg.PlayerPosition, cfg.LyricsPosition) {
		return "Static lyrics can be placed below the player or in a separate screen zone"
	}
	return fmt.Sprintf("Warning: %s is blocked because player already uses it", cfg.LyricsPosition.Label())
}
