package page

import (
	"virga-player/settings"

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
			return true, true, false
		case "Reset style to default":
			_ = settings.ResetStyleToDefault()
			return true, true, false
		case "Delete Virga (remove config and PATH aliases)":
			p.ConfirmDelete = true
			return false, false, false
		}
		return true, true, false
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
	p.Modified = true
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
		case 2:
			p.Config.Debug = !p.Config.Debug
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
			p.Config.RainLifetime += delta * 10
			if p.Config.RainLifetime < 20 {
				p.Config.RainLifetime = 20
			}
			if p.Config.RainLifetime > 200 {
				p.Config.RainLifetime = 200
			}
		case 2:
			p.Config.PulseSpeed += delta * 10
			if p.Config.PulseSpeed < 25 {
				p.Config.PulseSpeed = 25
			}
			if p.Config.PulseSpeed > 300 {
				p.Config.PulseSpeed = 300
			}
		case 3:
			options := settings.PulseModeOptions()
			index := p.pulseModeIndex()
			index = (index + delta + len(options)) % len(options)
			p.Config.SetPulseMode(options[index])
		case 4:
			p.Config.RainEnabled = !p.Config.RainEnabled
		case 5:
			p.Config.RainPulse += delta * 10
			if p.Config.RainPulse < 20 {
				p.Config.RainPulse = 20
			}
			if p.Config.RainPulse > 200 {
				p.Config.RainPulse = 200
			}
		case 6:
			p.Config.SeparateFrequencies = !p.Config.SeparateFrequencies
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
		case 2:
			p.Config.RainVisualizer = !p.Config.RainVisualizer
		}
	case sectionVisual:
		switch p.Selected {
		case 0:
			p.Config.MusicPlayerAnimation = !p.Config.MusicPlayerAnimation
		case 1:
			p.Config.MusicPlayerIntensity += delta * 10
			if p.Config.MusicPlayerIntensity < 20 {
				p.Config.MusicPlayerIntensity = 20
			}
			if p.Config.MusicPlayerIntensity > 200 {
				p.Config.MusicPlayerIntensity = 200
			}
		case 2:
			p.Config.MusicPlayerInvert = !p.Config.MusicPlayerInvert
		case 3:
			p.Config.RainInFrontOfPlayer = !p.Config.RainInFrontOfPlayer
		case 4:
			op := settings.DirectionOptions()
			index := p.directionIndex()
			index = (index + delta + len(op)) % len(op)
			p.Config.Direction = op[index]
		case 5:
			p.Config.Player = !p.Config.Player
		}
	}
}

func (p *Page) pulseModeIndex() int {
	options := settings.PulseModeOptions()
	for i, mode := range options {
		if mode == p.Config.PulseMode {
			return i
		}
	}
	return 0
}

func (p *Page) directionIndex() int {
	op := settings.DirectionOptions()
	for i, mode := range op {
		if mode == p.Config.Direction {
			return i
		}
	}
	return 0
}
