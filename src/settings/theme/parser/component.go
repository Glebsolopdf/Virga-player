package parser

import (
	"strconv"
	"strings"
)

func parseColorComponent(input string) (int, bool) {
	v := strings.TrimSpace(input)
	if strings.HasSuffix(v, "%") {
		numStr := strings.TrimSpace(strings.TrimSuffix(v, "%"))
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil || f < 0 || f > 100 {
			return 0, false
		}
		return int((f / 100.0) * 255.0), true
	}

	n, err := strconv.Atoi(v)
	if err != nil || n < 0 || n > 255 {
		return 0, false
	}
	return n, true
}

func parseAlpha(input string) (float64, bool) {
	v := strings.TrimSpace(input)
	if strings.HasSuffix(v, "%") {
		numStr := strings.TrimSpace(strings.TrimSuffix(v, "%"))
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil || f < 0 || f > 100 {
			return 0, false
		}
		return f / 100.0, true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f < 0 || f > 1 {
		return 0, false
	}
	return f, true
}
