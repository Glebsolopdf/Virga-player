package theme

import (
	"regexp"
	"strings"

	"virga-player/settings/theme/parser"

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

	varDecl := regexp.MustCompile(`--([a-z0-9-]+)\s*:\s*([^;]+);`)
	matches := varDecl.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		name := strings.TrimSpace(strings.ToLower(match[1]))
		value := strings.TrimSpace(match[2])

		if target, ok := colorVars[name]; ok {
			if parsed, valid := parser.ParseColor(value); valid {
				*target = parsed
			}
			continue
		}
		if target, ok := runeVars[name]; ok {
			if parsed, valid := parser.ParseRune(value); valid {
				*target = parsed
			}
		}
	}

	return result
}
