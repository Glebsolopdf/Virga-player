package welcome

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func (m model) renderMenu() string {
	state := m.menuAnim.state(m.frame)
	innerWidth := fitWidth(m.width-10, 56, 84)
	textFade := easeOut(state.textProgress)
	content := []string{
		renderMenuHeader(innerWidth, textFade),
		"",
	}

	for i, item := range m.items {
		content = append(content, m.renderMenuItem(item, i == m.selected, innerWidth, textFade)...)
		if i < len(m.items)-1 {
			content = append(content, "")
		}
	}
	content = append(content, "")
	content = append(content, renderMenuFooter(innerWidth, textFade))
	return renderAnimatedBox(content, innerWidth, state, float64(m.frame), borderStartTopRight)
}

func renderMenuHeader(width int, fade float64) string {
	color := fadeColor(
		colorful.LinearRgb(0.24, 0.25, 0.28),
		colorful.LinearRgb(0.82, 0.84, 0.88),
		fade,
	)
	return revealText(centerText("Choose your startup mode", width), fade, color, true)
}

func renderMenuFooter(width int, fade float64) string {
	color := fadeColor(
		colorful.LinearRgb(0.22, 0.23, 0.26),
		colorful.LinearRgb(0.62, 0.64, 0.69),
		fade,
	)
	return revealText(centerText("Enter to continue  •  Esc to exit", width), fade, color, false)
}

func (m model) renderMenuItem(item menuItem, selected bool, width int, fade float64) []string {
	cursor := "  "
	titleTarget := colorful.LinearRgb(0.86, 0.88, 0.92)
	if selected && m.menuReady() {
		cursor = "> "
		titleTarget = colorful.LinearRgb(0.78, 0.93, 1)
	}

	titleColor := fadeColor(colorful.LinearRgb(0.24, 0.25, 0.27), titleTarget, fade)
	descColor := fadeColor(colorful.LinearRgb(0.20, 0.21, 0.24), colorful.LinearRgb(0.64, 0.67, 0.72), fade)

	firstPrefix := " " + cursor
	nextPrefix := strings.Repeat(" ", lipgloss.Width(firstPrefix))
	lines := wrapWithPrefix(item.Title, firstPrefix, nextPrefix, width)
	for i, line := range lines {
		lines[i] = revealText(line, fade, titleColor, true)
	}

	descLines := wrapWithPrefix(item.Description, nextPrefix, nextPrefix, width)
	for i, line := range descLines {
		descLines[i] = revealText(line, fade, descColor, false)
	}

	return append(lines, descLines...)
}

func wrapWithPrefix(text, firstPrefix, nextPrefix string, width int) []string {
	limit := width - lipgloss.Width(firstPrefix)
	if limit < 12 {
		limit = 12
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{firstPrefix}
	}

	lines := []string{}
	current := firstPrefix
	currentWidth := 0
	prefix := firstPrefix

	for _, word := range words {
		for lipgloss.Width(word) > limit {
			head, tail := cutToWidth(word, limit)
			if currentWidth > 0 {
				lines = append(lines, current)
				prefix = nextPrefix
				current = prefix
				currentWidth = 0
			}
			lines = append(lines, prefix+head)
			word = tail
			prefix = nextPrefix
			current = prefix
		}

		wordWidth := lipgloss.Width(word)
		space := 0
		if currentWidth > 0 {
			space = 1
		}
		if currentWidth+space+wordWidth > limit {
			lines = append(lines, current)
			prefix = nextPrefix
			current = prefix + word
			currentWidth = wordWidth
			continue
		}
		if currentWidth > 0 {
			current += " "
		}
		current += word
		currentWidth += space + wordWidth
	}

	lines = append(lines, current)
	return lines
}

func cutToWidth(text string, limit int) (string, string) {
	runes := []rune(text)
	width := 0
	for i, r := range runes {
		runeWidth := lipgloss.Width(string(r))
		if width+runeWidth > limit {
			return string(runes[:i]), string(runes[i:])
		}
		width += runeWidth
	}
	return text, ""
}
