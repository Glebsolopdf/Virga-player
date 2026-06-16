package welcome

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func revealText(text string, progress float64, color string, bold bool) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Background(lipgloss.Color("#000000")).
		Bold(bold)
	total := countVisibleRunes(text)
	if total == 0 {
		return text
	}

	revealed := int(math.Round(progress * float64(total)))
	if progress > 0 && revealed == 0 {
		revealed = 1
	}

	var out strings.Builder
	seen := 0
	for _, r := range text {
		switch {
		case r == ' ':
			out.WriteRune(r)
		case seen < revealed:
			out.WriteString(style.Render(string(r)))
			seen++
		default:
			out.WriteRune(' ')
		}
	}
	return out.String()
}

func countVisibleRunes(text string) int {
	count := 0
	for _, r := range text {
		if r != ' ' {
			count++
		}
	}
	return count
}
