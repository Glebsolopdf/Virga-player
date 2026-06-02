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
	LyricsCurrent       tcell.Color
	LyricsInactive      tcell.Color
	LyricsPulse         tcell.Color
	LyricsBackground    tcell.Color
	LyricsGap           int
	TimelineBracket     tcell.Color
	TimelinePlayed      tcell.Color
	TimelineCurrent     tcell.Color
	TimelineRemaining   tcell.Color
	RainHead            tcell.Color
	RainTail            tcell.Color
	RainVeryNear        tcell.Color
	RainNear            tcell.Color
	RainMid             tcell.Color
	RainFar             tcell.Color
	RainVeryFar         tcell.Color
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
	ArtworkBlockRune    rune
}

var currentTheme = DefaultTheme()

func DefaultTheme() Theme {
	return Theme{
		Background:          tcell.ColorDefault,
		MessageText:         tcell.ColorWhite,
		TrackTitle:          tcell.ColorWhite,
		TrackArtist:         tcell.ColorGreen,
		TrackAlbum:          tcell.ColorYellow,
		TrackTime:           tcell.ColorGray,
		LyricsCurrent:       tcell.ColorWhite,
		LyricsInactive:      tcell.ColorLightGray,
		LyricsPulse:         tcell.ColorLightCyan,
		LyricsBackground:    tcell.ColorDefault,
		LyricsGap:           2,
		TimelineBracket:     tcell.ColorSilver,
		TimelinePlayed:      tcell.ColorGreen,
		TimelineCurrent:     tcell.ColorGreen,
		TimelineRemaining:   tcell.ColorGray,
		RainHead:            tcell.ColorWhite,
		RainTail:            tcell.ColorDefault,
		RainVeryNear:        tcell.ColorWhite,
		RainNear:            tcell.ColorLightCyan,
		RainMid:             tcell.ColorWhite,
		RainFar:             tcell.ColorLightGray,
		RainVeryFar:         tcell.ColorDarkGray,
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
		ArtworkBlockRune:    '▀',
	}
}

func SetCurrentTheme(theme Theme) {
	currentTheme = theme
}

func CurrentTheme() Theme {
	return currentTheme
}
