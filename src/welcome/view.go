package welcome

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	banner := strings.Split(m.renderBanner(), "\n")
	block := strings.Split(m.renderContentBlock(), "\n")
	group := append(append([]string{}, banner...), "")
	group = append(group, block...)

	width := 0
	for _, line := range group {
		if w := lipgloss.Width(line); w > width {
			width = w
		}
	}

	for i, line := range group {
		group[i] = lipgloss.PlaceHorizontal(width, lipgloss.Center, line)
	}

	return lipgloss.Place(
		max(1, m.width),
		max(1, m.height),
		lipgloss.Center,
		lipgloss.Center,
		strings.Join(group, "\n"),
	)
}

func (m model) renderContentBlock() string {
	if m.stage == stageControls {
		return m.renderControls()
	}
	return m.renderMenu()
}
