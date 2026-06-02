package page

import (
	"virga-player/notification"
	"virga-player/settings"
)

type Page struct {
	Config                 *settings.Config
	Selected               int
	Section                int
	ConfirmDelete          bool
	Modified               bool
	NotificationsSupported bool
	Notifications          []notification.Item
	OnOpenNotifications    func() []notification.Item
}

const (
	sectionNone    = -1
	sectionGeneral = iota
	sectionRain
	sectionAudio
	sectionVisual
	sectionLyrics
	sectionNotifications
	sectionReset
	sectionDanger
)

func NewPage(cfg *settings.Config) *Page {
	p := &Page{Config: cfg, Section: sectionNone, NotificationsSupported: true}
	p.Selected = p.firstSelectableIndex()
	return p
}

func (p *Page) SetNotificationsSupported(supported bool) {
	p.NotificationsSupported = supported
}

func (p *Page) SetNotifications(items []notification.Item, onOpen func() []notification.Item) {
	p.Notifications = append([]notification.Item(nil), items...)
	p.OnOpenNotifications = onOpen
}

type menuItem struct {
	label      string
	selectable bool
}
