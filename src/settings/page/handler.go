package page

import (
	"virga-player/settings"
	pagecontrols "virga-player/settings/page/controls"

	"github.com/gdamore/tcell/v2"
)

func (p *Page) HandleKey(ev *tcell.EventKey) (exit bool, save bool, deleteVirga bool) {
	if p.ConfirmDelete {
		switch ev.Key() {
		case tcell.KeyEscape:
			p.ConfirmDelete = false
			return false, false, false
		case tcell.KeyEnter:
			return true, false, true
		}

		switch ev.Rune() {
		case 'y', 'Y':
			return true, false, true
		case 'n', 'N', 's':
			p.ConfirmDelete = false
			return false, false, false
		}
	}

	switch ev.Key() {
	case tcell.KeyEscape:
		if p.Section != sectionNone {
			p.Section = sectionNone
			p.Selected = p.firstSelectableIndex()
			return false, false, false
		}
		return true, p.Modified, false
	case tcell.KeyEnter:
		if p.Section == sectionNone {
			p.selectSection(p.Selected)
			return false, false, false
		}
		items := p.menuItems()
		switch items[p.Selected].label {
		case "Back":
			p.Section = sectionNone
			p.Selected = p.firstSelectableIndex()
			return false, false, false
		case "Reset settings to default":
			p.Config = settings.DefaultConfig()
			p.Modified = true
			return true, true, false
		case "Reset style to default":
			_ = settings.ResetStyleToDefault()
			return true, p.Modified, false
		case "Delete Virga (remove config and PATH aliases)":
			p.ConfirmDelete = true
			return false, false, false
		}
		return true, p.Modified, false
	case tcell.KeyUp:
		p.ConfirmDelete = false
		p.moveSelection(-1)
	case tcell.KeyDown:
		p.ConfirmDelete = false
		p.moveSelection(1)
	case tcell.KeyLeft:
		if p.Section != sectionNone {
			p.adjust(-1)
		}
	case tcell.KeyRight:
		if p.Section != sectionNone {
			p.adjust(1)
		}
	}
	if ev.Rune() == 's' {
		if p.Section != sectionNone {
			p.Section = sectionNone
			p.Selected = p.firstSelectableIndex()
			return false, false, false
		}
		return true, p.Modified, false
	}
	return false, false, false
}

func (p *Page) adjust(delta int) {
	switch p.Section {
	case sectionGeneral:
		p.Modified = pagecontrols.General(p.Config, p.Selected, delta) || p.Modified
	case sectionRain:
		p.Modified = pagecontrols.Rain(p.Config, p.Selected, delta) || p.Modified
	case sectionAudio:
		p.Modified = pagecontrols.Audio(p.Config, p.Selected, delta) || p.Modified
	case sectionVisual:
		p.Modified = pagecontrols.Visual(p.Config, p.Selected, delta) || p.Modified
	}
}
