package welcome

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderAnimatedBox(content []string, innerWidth int, state blockState, phase float64, start borderStart) string {
	lines := padLines(content, innerWidth)
	width := innerWidth + 2
	height := len(lines) + 2
	perimeter := rectanglePerimeter(width, height, start)
	visible := int(math.Round(state.borderProgress * float64(len(perimeter))))
	if state.borderProgress > 0 && visible == 0 {
		visible = 1
	}

	border := make(map[point]string, visible)
	for i := 0; i < visible && i < len(perimeter); i++ {
		pt := perimeter[i]
		color := animatedBorderColor(i, len(perimeter), state.settleProgress, phase)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		border[pt] = style.Render(borderRune(pt.x, pt.y, width, height))
	}

	rows := make([]string, height)
	for y := 0; y < height; y++ {
		if y == 0 || y == height-1 {
			rows[y] = borderRow(border, width, y)
			continue
		}
		left := fallbackCell(border[point{x: 0, y: y}])
		right := fallbackCell(border[point{x: width - 1, y: y}])
		rows[y] = left + lines[y-1] + right
	}
	return strings.Join(rows, "\n")
}

func borderRow(border map[point]string, width, y int) string {
	var row strings.Builder
	for x := 0; x < width; x++ {
		row.WriteString(fallbackCell(border[point{x: x, y: y}]))
	}
	return row.String()
}

func borderRune(x, y, width, height int) string {
	switch {
	case x == 0 && y == 0:
		return "╭"
	case x == width-1 && y == 0:
		return "╮"
	case x == width-1 && y == height-1:
		return "╯"
	case x == 0 && y == height-1:
		return "╰"
	case y == 0 || y == height-1:
		return "─"
	default:
		return "│"
	}
}
