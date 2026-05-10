package theme

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func parseThemeCSS(input string, base Theme) Theme {
	result := base

	colorVars := map[string]*tcell.Color{
		"bg":                   &result.Background,
		"message-text":         &result.MessageText,
		"track-title":          &result.TrackTitle,
		"track-artist":         &result.TrackArtist,
		"track-album":          &result.TrackAlbum,
		"track-time":           &result.TrackTime,
		"timeline-bracket":     &result.TimelineBracket,
		"timeline-played":      &result.TimelinePlayed,
		"timeline-current":     &result.TimelineCurrent,
		"timeline-remaining":   &result.TimelineRemaining,
		"rain-head":            &result.RainHead,
		"rain-tail":            &result.RainTail,
		"settings-title":       &result.SettingsTitle,
		"settings-hint":        &result.SettingsHint,
		"settings-text":        &result.SettingsText,
		"settings-selected-fg": &result.SettingsSelectedFg,
		"settings-selected-bg": &result.SettingsSelectedBg,
		"settings-danger":      &result.SettingsDanger,
		"settings-danger-bg":   &result.SettingsDangerBg,
	}

	runeVars := map[string]*rune{
		"timeline-char-left":    &result.TimelineLeftRune,
		"timeline-char-right":   &result.TimelineRightRune,
		"timeline-char-played":  &result.TimelinePlayedRune,
		"timeline-char-current": &result.TimelineCurrentRune,
		"timeline-char-empty":   &result.TimelineEmptyRune,
		"rain-char-body":        &result.RainBodyRune,
		"rain-char-head":        &result.RainHeadRune,
		"rain-char-left":        &result.RainLeftRune,
		"rain-char-right":       &result.RainRightRune,
	}

	varDecl := regexp.MustCompile(`--([a-z0-9-]+)\\s*:\\s*([^;]+);`)
	matches := varDecl.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		name := strings.TrimSpace(strings.ToLower(match[1]))
		value := strings.TrimSpace(match[2])

		if target, ok := colorVars[name]; ok {
			if parsed, valid := parseColor(value); valid {
				*target = parsed
			}
			continue
		}
		if target, ok := runeVars[name]; ok {
			if parsed, valid := parseRune(value); valid {
				*target = parsed
			}
		}
	}

	return result
}

func parseColor(value string) (tcell.Color, bool) {
	v := strings.TrimSpace(strings.Trim(value, `"'`))
	if v == "" {
		return tcell.ColorDefault, false
	}

	vLower := strings.ToLower(v)
	if vLower == "default" || vLower == "reset" || vLower == "transparent" {
		return tcell.ColorDefault, true
	}

	if strings.HasPrefix(vLower, "rgba(") && strings.HasSuffix(vLower, ")") {
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(vLower, "rgba("), ")"))
		r, g, b, a, ok := parseRGBOrRGBA(inner, true)
		if ok {
			if a <= 0 {
				return tcell.ColorDefault, true
			}
			return tcell.NewRGBColor(int32(r), int32(g), int32(b)), true
		}
		return tcell.ColorDefault, false
	}

	if strings.HasPrefix(vLower, "rgb(") && strings.HasSuffix(vLower, ")") {
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(vLower, "rgb("), ")"))
		r, g, b, _, ok := parseRGBOrRGBA(inner, false)
		if ok {
			return tcell.NewRGBColor(int32(r), int32(g), int32(b)), true
		}
		return tcell.ColorDefault, false
	}

	if strings.HasPrefix(v, "#") && len(v) == 9 {
		r, errR := strconv.ParseInt(v[1:3], 16, 32)
		g, errG := strconv.ParseInt(v[3:5], 16, 32)
		b, errB := strconv.ParseInt(v[5:7], 16, 32)
		a, errA := strconv.ParseInt(v[7:9], 16, 32)
		if errR == nil && errG == nil && errB == nil && errA == nil {
			if a == 0 {
				return tcell.ColorDefault, true
			}
			return tcell.NewRGBColor(int32(r), int32(g), int32(b)), true
		}
		return tcell.ColorDefault, false
	}

	if strings.HasPrefix(v, "#") && len(v) == 7 {
		r, errR := strconv.ParseInt(v[1:3], 16, 32)
		g, errG := strconv.ParseInt(v[3:5], 16, 32)
		b, errB := strconv.ParseInt(v[5:7], 16, 32)
		if errR == nil && errG == nil && errB == nil {
			return tcell.NewRGBColor(int32(r), int32(g), int32(b)), true
		}
		return tcell.ColorDefault, false
	}

	named := map[string]tcell.Color{
		"black":   tcell.ColorBlack,
		"white":   tcell.ColorWhite,
		"gray":    tcell.ColorGray,
		"grey":    tcell.ColorGray,
		"silver":  tcell.ColorSilver,
		"red":     tcell.ColorRed,
		"maroon":  tcell.ColorMaroon,
		"green":   tcell.ColorGreen,
		"yellow":  tcell.ColorYellow,
		"blue":    tcell.ColorBlue,
		"navy":    tcell.ColorNavy,
		"teal":    tcell.ColorTeal,
		"aqua":    tcell.ColorAqua,
		"fuchsia": tcell.ColorFuchsia,
		"purple":  tcell.ColorPurple,
		"olive":   tcell.ColorOlive,
	}
	if c, ok := named[strings.ToLower(v)]; ok {
		return c, true
	}

	return tcell.ColorDefault, false
}

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

func parseRune(value string) (rune, bool) {
	v := strings.TrimSpace(value)
	v = strings.Trim(v, `"'`)
	v = strings.ReplaceAll(v, `\\`, `\`)
	v = strings.ReplaceAll(v, `\t`, "\t")
	v = strings.ReplaceAll(v, `\n`, "\n")
	if v == "" {
		return 0, false
	}
	r := []rune(v)
	if len(r) == 0 {
		return 0, false
	}
	return r[0], true
}
