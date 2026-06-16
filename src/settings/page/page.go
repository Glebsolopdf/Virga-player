package page

import (
	"virga-player/settings"
)

type Page struct {
	Config        *settings.Config
	Selected      int
	Section       int
	ConfirmDelete bool
	Modified      bool
}

const (
	sectionNone    = -1
	sectionGeneral = iota
	sectionRain
	sectionAudio
	sectionVisual
	sectionLyrics
	sectionReset
	sectionDanger
)

func NewPage(cfg *settings.Config) *Page {
	p := &Page{Config: cfg, Section: sectionNone}
	p.Selected = p.firstSelectableIndex()
	return p
}

// notifications removed

type menuItem struct {
	label      string
	selectable bool
}
