package page

import (
	"fmt"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/version"

	"github.com/gdamore/tcell/v2"
)

func (p *Page) Render(screen tcell.Screen, renderer *renderer.Renderer, width, height int) {
	theme := settings.CurrentTheme()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			renderer.DrawRune(x, y, ' ', tcell.ColorReset, theme.Background)
		}
	}

	title := "Virga Player Settings"
	renderer.DrawTextCentered(screen, 2, title, theme.SettingsTitle, theme.Background)

	subtitle := "Use arrows to select a category and Enter to open it"
	switch p.Section {
	case sectionNone:
		// keep default subtitle
	default:
		subtitle = "Use Left/Right to change values, Enter to save and exit, Esc to go back"
	}
	renderer.DrawTextCentered(screen, 4, subtitle, theme.SettingsHint, theme.Background)

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
			renderer.DrawTextCentered(screen, startY+i*2, item.label, fg, bg)
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
			renderer.DrawTextCentered(screen, startY+i*2, item.label, fg, bg)
		}
	}

	if p.Section == sectionNone {
		dir := settings.ConfigDirPath()
		renderer.DrawTextCentered(screen, height-6, fmt.Sprintf("Player directory: %s", dir), theme.SettingsHint, theme.Background)
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
	renderer.DrawTextCentered(screen, height-4, helpText, theme.SettingsHint, theme.Background)
	footerText := fmt.Sprintf("%s | %s", version.AppVersion, version.GitHubURL)
	renderer.DrawTextCentered(screen, height-2, footerText, theme.SettingsHint, theme.Background)
}
