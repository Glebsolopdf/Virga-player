package settings

import (
	"fmt"

	"virga-player/renderer"

	"github.com/gdamore/tcell/v2"
)

type Page struct {
	Config        *Config
	Selected      int
	Section       int
	ConfirmDelete bool
}

type menuItem struct {
	label      string
	selectable bool
}

const (
	sectionNone    = -1
	sectionGeneral = iota
	sectionRain
	sectionAudio
	sectionVisual
	sectionReset
	sectionDanger
)

func NewPage(cfg *Config) *Page {
	p := &Page{Config: cfg, Section: sectionNone}
	p.Selected = p.firstSelectableIndex()
	return p
}

func (p *Page) topMenuItems() []menuItem {
	return []menuItem{
		{label: "General settings", selectable: true},
		{label: "Rain settings", selectable: true},
		{label: "Audio / reactive", selectable: true},
		{label: "Visual settings", selectable: true},
		{label: "Reset actions", selectable: true},
		{label: "Danger zone", selectable: true},
	}
}

func (p *Page) sectionMenuItems() []menuItem {
	switch p.Section {
	case sectionGeneral:
		return []menuItem{
			{label: fmt.Sprintf("Frame rate: %d", p.Config.FPS), selectable: true},
			{label: fmt.Sprintf("Max particles: %d", p.Config.MaxParticles), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionRain:
		return []menuItem{
			{label: fmt.Sprintf("Rain speed: %d%%", p.Config.RainSpeed), selectable: true},
			{label: fmt.Sprintf("Rain enabled: %v", p.Config.RainEnabled), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionAudio:
		return []menuItem{
			{label: fmt.Sprintf("Music reactive: %v", p.Config.MusicReactive), selectable: true},
			{label: fmt.Sprintf("Reactive intensity: %d%%", p.Config.MusicReactiveIntensity), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionVisual:
		return []menuItem{
			{label: fmt.Sprintf("Cover animation: %v", p.Config.CoverAnimation), selectable: true},
			{label: fmt.Sprintf("Rain direction: %s", p.Config.Direction.Label()), selectable: true},
			{label: fmt.Sprintf("Player mode: %v", p.Config.Player), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionReset:
		return []menuItem{
			{label: "Reset settings to default", selectable: true},
			{label: "Reset style to default", selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionDanger:
		return []menuItem{
			{label: "Delete Virga (remove config and PATH aliases)", selectable: true},
			{label: "Back", selectable: true},
		}
	default:
		return nil
	}
}

func (p *Page) menuItems() []menuItem {
	if p.Section == sectionNone {
		return p.topMenuItems()
	}
	return p.sectionMenuItems()
}

func (p *Page) firstSelectableIndex() int {
	for i, item := range p.menuItems() {
		if item.selectable {
			return i
		}
	}
	return 0
}

func (p *Page) moveSelection(delta int) {
	items := p.menuItems()
	count := len(items)
	if count == 0 {
		return
	}
	newIndex := p.Selected
	for {
		newIndex = (newIndex + delta + count) % count
		if items[newIndex].selectable {
			p.Selected = newIndex
			return
		}
	}
}

func (p *Page) selectSection(index int) {
	switch index {
	case 0:
		p.Section = sectionGeneral
	case 1:
		p.Section = sectionRain
	case 2:
		p.Section = sectionAudio
	case 3:
		p.Section = sectionVisual
	case 4:
		p.Section = sectionReset
	case 5:
		p.Section = sectionDanger
	}
	p.Selected = p.firstSelectableIndex()
}

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
		return true, false, false
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
			p.Config = DefaultConfig()
			return false, false, false
		case "Reset style to default":
			_ = ResetStyleToDefault()
			return false, false, false
		case "Delete Virga (remove config and PATH aliases)":
			p.ConfirmDelete = true
			return false, false, false
		default:
			return true, true, false
		}
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
		return true, false, false
	}
	return false, false, false
}

func (p *Page) adjust(delta int) {
	switch p.Section {
	case sectionGeneral:
		switch p.Selected {
		case 0:
			p.Config.FPS += delta * 5
			if p.Config.FPS < 15 {
				p.Config.FPS = 15
			}
			if p.Config.FPS > 240 {
				p.Config.FPS = 240
			}
		case 1:
			p.Config.MaxParticles += delta * 10
			if p.Config.MaxParticles < 20 {
				p.Config.MaxParticles = 20
			}
			if p.Config.MaxParticles > 500 {
				p.Config.MaxParticles = 500
			}
		}
	case sectionRain:
		switch p.Selected {
		case 0:
			p.Config.RainSpeed += delta * 5
			if p.Config.RainSpeed < 25 {
				p.Config.RainSpeed = 25
			}
			if p.Config.RainSpeed > 300 {
				p.Config.RainSpeed = 300
			}
		case 1:
			p.Config.RainEnabled = !p.Config.RainEnabled
		}
	case sectionAudio:
		switch p.Selected {
		case 0:
			p.Config.MusicReactive = !p.Config.MusicReactive
		case 1:
			p.Config.MusicReactiveIntensity += delta * 10
			if p.Config.MusicReactiveIntensity < 20 {
				p.Config.MusicReactiveIntensity = 20
			}
			if p.Config.MusicReactiveIntensity > 200 {
				p.Config.MusicReactiveIntensity = 200
			}
		}
	case sectionVisual:
		switch p.Selected {
		case 0:
			p.Config.CoverAnimation = !p.Config.CoverAnimation
		case 1:
			op := DirectionOptions()
			index := p.directionIndex()
			index = (index + delta + len(op)) % len(op)
			p.Config.Direction = op[index]
		case 2:
			p.Config.Player = !p.Config.Player
		}
	}
}

func (p *Page) directionIndex() int {
	op := DirectionOptions()
	for i, mode := range op {
		if mode == p.Config.Direction {
			return i
		}
	}
	return 0
}

func (p *Page) Render(screen tcell.Screen, renderer *renderer.Renderer, width, height int) {
	theme := CurrentTheme()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			renderer.DrawRune(x, y, ' ', tcell.ColorReset, theme.Background)
		}
	}

	title := "Virga Player Settings"
	renderer.DrawTextCentered(screen, 2, title, theme.SettingsTitle, theme.Background)

	subtitle := "Use arrows to move, Enter to open/confirm, Esc to go back/close"
	if p.Section == sectionNone {
		subtitle = "Use arrows to select a category and Enter to open it"
	}
	renderer.DrawTextCentered(screen, 4, subtitle, theme.SettingsHint, theme.Background)

	items := p.menuItems()
	startY := 8
	for i, item := range items {
		fg := theme.SettingsText
		bg := theme.Background
		if i == p.Selected {
			fg = theme.SettingsSelectedFg
			bg = theme.SettingsSelectedBg
		}
		renderer.DrawText(screen, 6, startY+i*2, item.label, fg, bg)
	}

	helpText := "Esc/s: back or exit"
	if p.ConfirmDelete {
		helpText = "Delete Virga confirmation: Enter/Y confirm, Esc/N cancel"
	}
	renderer.DrawTextCentered(screen, height-3, helpText, theme.SettingsHint, theme.Background)
}
