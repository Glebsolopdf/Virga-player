package welcome

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type point struct {
	x int
	y int
}

type borderStart int

const (
	borderStartTopLeft borderStart = iota
	borderStartTopRight
)

func rectanglePerimeter(width, height int, start borderStart) []point {
	if start == borderStartTopRight {
		return rectanglePerimeterTopRight(width, height)
	}
	return rectanglePerimeterTopLeft(width, height)
}

func rectanglePerimeterTopLeft(width, height int) []point {
	points := make([]point, 0, width*2+height*2-4)
	for x := 0; x < width; x++ {
		points = append(points, point{x: x, y: 0})
	}
	for y := 1; y < height; y++ {
		points = append(points, point{x: width - 1, y: y})
	}
	for x := width - 2; x >= 0; x-- {
		points = append(points, point{x: x, y: height - 1})
	}
	for y := height - 2; y >= 1; y-- {
		points = append(points, point{x: 0, y: y})
	}
	return points
}

func rectanglePerimeterTopRight(width, height int) []point {
	points := make([]point, 0, width*2+height*2-4)
	for x := width - 1; x >= 0; x-- {
		points = append(points, point{x: x, y: 0})
	}
	for y := 1; y < height; y++ {
		points = append(points, point{x: 0, y: y})
	}
	for x := 1; x < width; x++ {
		points = append(points, point{x: x, y: height - 1})
	}
	for y := height - 2; y >= 1; y-- {
		points = append(points, point{x: width - 1, y: y})
	}
	return points
}

func padLines(lines []string, width int) []string {
	padded := make([]string, len(lines))
	for i, line := range lines {
		diff := width - lipgloss.Width(line)
		if diff < 0 {
			diff = 0
		}
		padded[i] = line + strings.Repeat(" ", diff)
	}
	return padded
}

func centerText(text string, width int) string {
	gap := width - lipgloss.Width(text)
	if gap <= 0 {
		return text
	}
	left := gap / 2
	right := gap - left
	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", left), text, strings.Repeat(" ", right))
}

func fallbackCell(value string) string {
	if value == "" {
		return " "
	}
	return value
}

func fitWidth(available, minWidth, maxWidth int) int {
	available = max(12, available)
	if available < minWidth {
		return available
	}
	if available > maxWidth {
		return maxWidth
	}
	return available
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
