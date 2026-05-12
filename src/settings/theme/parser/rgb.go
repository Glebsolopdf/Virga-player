package parser

import "strings"

func parseRGBOrRGBA(input string, allowAlpha bool) (int, int, int, float64, bool) {
	trimmed := strings.TrimSpace(input)
	alpha := 1.0
	separatorIndex := strings.LastIndex(trimmed, "/")
	if separatorIndex != -1 {
		if !allowAlpha {
			return 0, 0, 0, 0, false
		}
		alphaStr := strings.TrimSpace(trimmed[separatorIndex+1:])
		trimmed = strings.TrimSpace(trimmed[:separatorIndex])
		parsedAlpha, ok := parseAlpha(alphaStr)
		if !ok {
			return 0, 0, 0, 0, false
		}
		alpha = parsedAlpha
	}

	pieces := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	filtered := make([]string, 0, len(pieces))
	for _, p := range pieces {
		if s := strings.TrimSpace(p); s != "" {
			filtered = append(filtered, s)
		}
	}

	if len(filtered) == 4 {
		if !allowAlpha {
			return 0, 0, 0, 0, false
		}
		parsedAlpha, okA := parseAlpha(filtered[3])
		if !okA {
			return 0, 0, 0, 0, false
		}
		alpha = parsedAlpha
		filtered = filtered[:3]
	}

	if len(filtered) != 3 {
		return 0, 0, 0, 0, false
	}

	r, okR := parseColorComponent(filtered[0])
	g, okG := parseColorComponent(filtered[1])
	b, okB := parseColorComponent(filtered[2])
	if !okR || !okG || !okB {
		return 0, 0, 0, 0, false
	}

	return r, g, b, alpha, true
}
