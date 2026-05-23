package overlay

import "strings"

func WrapText(text string, width int) []string {
	if width <= 0 {
		return nil
	}
	parts := strings.Split(strings.TrimSpace(text), "\n")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line == "" {
			out = append(out, "")
			continue
		}
		for len([]rune(line)) > width {
			cut := lastSpaceWithin(line, width)
			if cut <= 0 {
				cut = width
			}
			chunk := strings.TrimSpace(string([]rune(line)[:cut]))
			if chunk != "" {
				out = append(out, chunk)
			}
			line = strings.TrimSpace(string([]rune(line)[cut:]))
		}
		out = append(out, line)
	}
	if len(out) == 0 {
		return []string{""}
	}
	return out
}

func lastSpaceWithin(s string, width int) int {
	r := []rune(s)
	if len(r) <= width {
		return len(r)
	}
	for i := width; i > 0; i-- {
		if r[i-1] == ' ' || r[i-1] == '\t' {
			return i - 1
		}
	}
	return width
}
