package welcome

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"virga-player/settings"
)

func (m model) renderControls() string {
	state := m.menuAnim.state(m.frame)
	innerWidth := fitWidth(m.width-10, 56, 84)
	textFade := easeOut(state.textProgress)

	titleColor := fadeColor(
		colorful.LinearRgb(0.24, 0.25, 0.28),
		colorful.LinearRgb(0.82, 0.84, 0.88),
		textFade,
	)
	bodyColor := fadeColor(
		colorful.LinearRgb(0.20, 0.21, 0.24),
		colorful.LinearRgb(0.68, 0.70, 0.74),
		textFade,
	)

	content := []string{
		revealText(centerText("Control guide", innerWidth), textFade, titleColor, true),
		"",
	}

	lines := []string{
		`While you are on the main screen, press "S" on an English keyboard layout to open the settings menu.`,
		"Press ESC to exit Virga.",
		fmt.Sprintf("Launch commands: %s or %s", "virga", "virgaplayer"),
		fmt.Sprintf("Data directory: %s", settings.ConfigDirPath()),
		fmt.Sprintf("Style file: %s", settings.StylePath()),
		fmt.Sprintf("Config file: %s", settings.ConfigPath()),
		fmt.Sprintf("Lyrics directory: %s", m.selectedLyricsDir()),
	}
	for i, line := range lines {
		wrapped := wrapWithPrefix(line, "  ", "  ", innerWidth)
		for _, part := range wrapped {
			content = append(content, revealText(part, textFade, bodyColor, false))
		}
		if i < len(lines)-1 {
			content = append(content, "")
		}
	}

	content = append(content, "")
	content = append(content, revealText(centerText("Press Enter to start Virga", innerWidth), textFade, bodyColor, false))
	return renderAnimatedBox(content, innerWidth, state, float64(m.frame), borderStartTopRight)
}

func (m model) selectedLyricsDir() string {
	if m.chosen != nil && m.chosen.LyricsPersistentDir != "" {
		return m.chosen.LyricsPersistentDir
	}
	return settings.DefaultConfig().LyricsPersistentDir
}
