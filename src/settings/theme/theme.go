package theme

import (
	"github.com/gdamore/tcell/v2"
)

type Theme struct {
	Background          tcell.Color
	MessageText         tcell.Color
	TrackTitle          tcell.Color
	TrackArtist         tcell.Color
	TrackAlbum          tcell.Color
	TrackTime           tcell.Color
	TimelineBracket     tcell.Color
	TimelinePlayed      tcell.Color
	TimelineCurrent     tcell.Color
	TimelineRemaining   tcell.Color
	RainHead            tcell.Color
	RainTail            tcell.Color
	SettingsTitle       tcell.Color
	SettingsHint        tcell.Color
	SettingsText        tcell.Color
	SettingsSelectedFg  tcell.Color
	SettingsSelectedBg  tcell.Color
	SettingsDanger      tcell.Color
	SettingsDangerBg    tcell.Color
	TimelineLeftRune    rune
	TimelineRightRune   rune
	TimelinePlayedRune  rune
	TimelineCurrentRune rune
	TimelineEmptyRune   rune
	RainBodyRune        rune
	RainHeadRune        rune
	RainLeftRune        rune
	RainRightRune       rune
}

var currentTheme = DefaultTheme()
var currentThemeFileContent string

func DefaultTheme() Theme {
	return Theme{
		Background:          tcell.ColorDefault,
		MessageText:         tcell.ColorWhite,
		TrackTitle:          tcell.ColorWhite,
		TrackArtist:         tcell.ColorGreen,
		TrackAlbum:          tcell.ColorYellow,
		TrackTime:           tcell.ColorGray,
		TimelineBracket:     tcell.ColorSilver,
		TimelinePlayed:      tcell.ColorGreen,
		TimelineCurrent:     tcell.ColorGreen,
		TimelineRemaining:   tcell.ColorGray,
		RainHead:            tcell.ColorWhite,
		RainTail:            tcell.ColorDefault,
		SettingsTitle:       tcell.ColorWhite,
		SettingsHint:        tcell.ColorGray,
		SettingsText:        tcell.ColorWhite,
		SettingsSelectedFg:  tcell.ColorBlack,
		SettingsSelectedBg:  tcell.ColorWhite,
		SettingsDanger:      tcell.ColorRed,
		SettingsDangerBg:    tcell.ColorMaroon,
		TimelineLeftRune:    '[',
		TimelineRightRune:   ']',
		TimelinePlayedRune:  '█',
		TimelineCurrentRune: '▌',
		TimelineEmptyRune:   '░',
		RainBodyRune:        '│',
		RainHeadRune:        '•',
		RainLeftRune:        '/',
		RainRightRune:       '\\',
	}
}

func SetCurrentTheme(theme Theme) {
	currentTheme = theme
}

func CurrentTheme() Theme {
	return currentTheme
}
