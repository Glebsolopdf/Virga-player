package page

import (
	"fmt"
	"unicode/utf8"

	"virga-player/settings"
	"virga-player/version"

	"github.com/gdamore/tcell/v2"
)

func drawTextCentered(screen tcell.Screen, y int, text string, fg, bg tcell.Color) {
	w, _ := screen.Size()
	x := (w - utf8.RuneCountInString(text)) / 2
	if x < 0 {
		x = 0
	}
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for i, ch := range []rune(text) {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func (p *Page) Render(screen tcell.Screen, width, height int) {
	theme := settings.CurrentTheme()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorReset).Background(theme.Background))
		}
	}

	title := "Virga Player Settings"
	drawTextCentered(screen, 2, title, theme.SettingsTitle, theme.Background)

	subtitle := "Use arrows to select a category and Enter to open it"
	switch p.Section {
	case sectionNone:
	default:
		subtitle = "Use Left/Right to change values, Enter to save and exit, Esc to go back"
	}
	drawTextCentered(screen, 4, subtitle, theme.SettingsHint, theme.Background)

	if p.Section == sectionNone {
		items := p.menuItems()
		startY := 8
		for i, item := range items {
			fg := theme.SettingsText
			bg := theme.Background
			if !item.selectable {
				fg = theme.SettingsHint
			}
			if i == p.Selected {
				fg = theme.SettingsSelectedFg
				bg = theme.SettingsSelectedBg
			}
			drawTextCentered(screen, startY+i*2, item.label, fg, bg)
		}
	} else {
		items := p.sectionMenuItems()
		startY := 8
		for i, item := range items {
			fg := theme.SettingsText
			bg := theme.Background
			if !item.selectable {
				fg = theme.SettingsHint
			}
			if i == p.Selected {
				fg = theme.SettingsSelectedFg
				bg = theme.SettingsSelectedBg
			}
			drawTextCentered(screen, startY+i*2, item.label, fg, bg)
		}
	}

	if p.Section == sectionNone {
		dir := settings.ConfigDirPath()
		drawTextCentered(screen, height-6, fmt.Sprintf("Player directory: %s", dir), theme.SettingsHint, theme.Background)
	}

	helpText := "Esc: back/cancel  |  s: save and exit"
	switch p.Section {
	case sectionNone:
		helpText = "Esc: cancel and exit  |  s: save and exit"
	}
	if p.Config.Debug {
		helpText += "  |  C: copy logs  K: save logs"
	}
	if p.ConfirmDelete {
		helpText = "Delete Virga confirmation: Enter/Y confirm, Esc/N cancel"
	}
	drawTextCentered(screen, height-4, helpText, theme.SettingsHint, theme.Background)
	footerText := fmt.Sprintf("%s | %s", version.AppVersion, version.GitHubURL)
	drawTextCentered(screen, height-2, footerText, theme.SettingsHint, theme.Background)
}
