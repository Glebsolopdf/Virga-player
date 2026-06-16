package page

import (
	"fmt"

	"virga-player/settings"
)

func (p *Page) lyricsMenuItems() []menuItem {
	promptHint := "Prompt mode: request stays active for current song after hint fades"
	doubleConfirmHint := "Prompt mode: confirm with Y or Enter"
	if p.Config.LyricsDoubleConfirm {
		doubleConfirmHint = "Prompt mode: confirm with Y, then Y/Enter within 10s"
	}
	if p.Config.LyricsMode != settings.LyricsModeRAMWithPrompt {
		promptHint = "Prompt mode details appear when RAM + save prompt is selected"
		doubleConfirmHint = ""
	}

	items := []menuItem{
		{label: fmt.Sprintf("Show lyrics in player: %v", p.Config.LyricsVisible), selectable: true},
		{label: fmt.Sprintf("Lyrics mode: %s", p.Config.LyricsMode.Label()), selectable: true},
		{label: fmt.Sprintf("Lyrics display: %s", p.Config.LyricsDisplayMode.Label()), selectable: true},
		{label: fmt.Sprintf("Static lyrics position: %s", p.Config.LyricsPosition.Label()), selectable: true},
		{label: fmt.Sprintf("Auto-save / prompt delay: %ds", p.Config.LyricsAutoSaveAfterSec), selectable: true},
		{label: fmt.Sprintf("Lyrics rain layer: %s", p.Config.LyricsRainLayer.Label()), selectable: true},
		{label: fmt.Sprintf("Prompt double confirmation: %v", p.Config.LyricsDoubleConfirm), selectable: true},
		{label: lyricsWarning(p.Config), selectable: false},
		{label: "Dynamic mode shifts each new line by 2 or 3 cells in a new direction", selectable: false},
		{label: "Source: LRCLIB API", selectable: false},
		{label: "Auto and prompt modes require internet access", selectable: false},
		{label: "Direct to disk: use for stable playlists or favorite albums", selectable: false},
		{label: promptHint, selectable: false},
	}
	if doubleConfirmHint != "" {
		items = append(items, menuItem{label: doubleConfirmHint, selectable: false})
	}
	return append(items, menuItem{label: "Back", selectable: true})
}
