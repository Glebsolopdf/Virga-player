package parser

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func ParseColor(value string) (tcell.Color, bool) {
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
