package page

import (
	"fmt"

	"virga-player/settings"
)

func (p *Page) topMenuItems() []menuItem {
	return []menuItem{
		{label: "General settings", selectable: true},
		{label: "Rain settings", selectable: true},
		{label: "Audio / reactive", selectable: true},
		{label: "Visual settings", selectable: true},
		{label: "Lyrics settings", selectable: true},
		{label: "Notifications", selectable: true},
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
			{label: fmt.Sprintf("Debug mode: %v", p.Config.Debug), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionRain:
		return []menuItem{
			{label: fmt.Sprintf("Rain speed: %d%%", p.Config.RainSpeed), selectable: true},
			{label: fmt.Sprintf("Rain duration: %d%%", p.Config.RainLifetime), selectable: true},
			{label: fmt.Sprintf("Pulse speed: %d%%", p.Config.PulseSpeed), selectable: true},
			{label: fmt.Sprintf("Pulse target: %s", p.Config.PulseMode.Label()), selectable: true},
			{label: fmt.Sprintf("Rain enabled: %v", p.Config.RainEnabled), selectable: true},
			{label: fmt.Sprintf("Rain pulse: %d%%", p.Config.RainPulse), selectable: true},
			{label: fmt.Sprintf("Separate frequencies: %v", p.Config.SeparateFrequencies), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionAudio:
		return []menuItem{
			{label: fmt.Sprintf("Music reactive: %v", p.Config.MusicReactive), selectable: true},
			{label: fmt.Sprintf("Reactive intensity: %d%%", p.Config.MusicReactiveIntensity), selectable: true},
			{label: fmt.Sprintf("Rain visualizer: %v", p.Config.RainVisualizer), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionVisual:
		return []menuItem{
			{label: fmt.Sprintf("Player music animation: %v", p.Config.MusicPlayerAnimation), selectable: true},
			{label: fmt.Sprintf("Music intensity: %d%%", p.Config.MusicPlayerIntensity), selectable: true},
			{label: fmt.Sprintf("Invert music motion: %v", p.Config.MusicPlayerInvert), selectable: true},
			{label: fmt.Sprintf("Player rain layer: %s", p.Config.PlayerRainLayer.Label()), selectable: true},
			{label: fmt.Sprintf("Rain direction: %s", p.Config.Direction.Label()), selectable: true},
			{label: fmt.Sprintf("Player mode: %v", p.Config.Player), selectable: true},
			{label: "Back", selectable: true},
		}
	case sectionLyrics:
		promptHint := "Prompt mode: request stays active for current song after hint fades"
		doubleConfirmHint := "Prompt mode: confirm with Y or Enter"
		if p.Config.LyricsDoubleConfirm {
			doubleConfirmHint = "Prompt mode: confirm with Y, then Y/Enter within 10s"
		}
		if p.Config.LyricsMode != settings.LyricsModeRAMWithPrompt {
			promptHint = "Prompt mode details appear when RAM + save prompt is selected"
			doubleConfirmHint = ""
		}

		items := []menuItem{
			{label: fmt.Sprintf("Show lyrics in player: %v", p.Config.LyricsVisible), selectable: true},
			{label: fmt.Sprintf("Lyrics mode: %s", p.Config.LyricsMode.Label()), selectable: true},
			{label: fmt.Sprintf("Auto-save / prompt delay: %ds", p.Config.LyricsAutoSaveAfterSec), selectable: true},
			{label: fmt.Sprintf("Lyrics rain layer: %s", p.Config.LyricsRainLayer.Label()), selectable: true},
			{label: fmt.Sprintf("Prompt double confirmation: %v", p.Config.LyricsDoubleConfirm), selectable: true},
			{label: "Source: LRCLIB API", selectable: false},
			{label: "Auto and prompt modes require internet access", selectable: false},
			{label: "Direct to disk: use for stable playlists or favorite albums", selectable: false},
			{label: promptHint, selectable: false},
		}
		if doubleConfirmHint != "" {
			items = append(items, menuItem{label: doubleConfirmHint, selectable: false})
		}
		items = append(items, menuItem{label: "Back", selectable: true})
		return items
	case sectionNotifications:
		items := []menuItem{
			{label: fmt.Sprintf("Receive notifications: %v", p.Config.NotificationsEnabled), selectable: true},
			{label: fmt.Sprintf("Сообщать о непрочитанных уведомлениях: %v", p.Config.NotifyUnreadToast), selectable: true},
		}
		items = append(items, menuItem{label: "Back", selectable: true})
		return items
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
		p.Section = sectionLyrics
	case 5:
		p.Section = sectionNotifications
		if p.OnOpenNotifications != nil {
			p.Notifications = p.OnOpenNotifications()
		}
	case 6:
		p.Section = sectionReset
	case 7:
		p.Section = sectionDanger
	}
	p.Selected = p.firstSelectableIndex()
}
