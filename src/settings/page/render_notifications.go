package page

import (
	"strings"

	"virga-player/renderer"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (p *Page) renderNotifications(screen tcell.Screen, renderer *renderer.Renderer, width, height int) {
	theme := settings.CurrentTheme()
	if width <= 0 || height <= 0 {
		return
	}

	contentX := 4
	if width < 16 {
		contentX = 1
	}
	contentWidth := width - (contentX * 2)
	if contentWidth < 1 {
		contentWidth = 1
	}
	backY := height - 7
	if backY < 8 {
		backY = height - 3
	}
	contentBottom := backY - 2
	if contentBottom < 8 {
		contentBottom = height - 8
	}

	y := 8
	toggleText := "Receive notifications: false"
	if p.Config != nil && p.Config.NotificationsEnabled {
		toggleText = "Receive notifications: true"
	}
	toggleFg := theme.SettingsText
	toggleBg := theme.Background
	if p.Selected == 0 {
		toggleFg = theme.SettingsSelectedFg
		toggleBg = theme.SettingsSelectedBg
	}
	renderer.DrawTextCentered(screen, y-2, toggleText, toggleFg, toggleBg)

	toastText := "Notify about unread notifications: false"
	if p.Config != nil && p.Config.NotifyUnreadToast {
		toastText = "Notify about unread notifications: true"
	}
	toastFg := theme.SettingsText
	toastBg := theme.Background
	if p.Selected == 1 {
		toastFg = theme.SettingsSelectedFg
		toastBg = theme.SettingsSelectedBg
	}
	renderer.DrawTextCentered(screen, y, toastText, toastFg, toastBg)

	truncated := false
	if len(p.Notifications) == 0 {
		renderer.DrawTextCentered(screen, y+3, "No notifications yet.", theme.SettingsHint, theme.Background)
	} else {
		y += 2
		for _, item := range p.Notifications {
			if y > contentBottom {
				truncated = true
				break
			}

			timestamp := trimToWidth(item.CreatedAt.Local().Format("2006-01-02 15:04"), contentWidth)
			renderer.DrawText(screen, contentX, y, timestamp, theme.SettingsHint, theme.Background)
			y++

			title := item.Title
			if !item.IsRead() {
				title = "[NEW] " + title
			}
			titleColor := theme.SettingsText
			if !item.IsRead() {
				titleColor = theme.SettingsTitle
			}
			for _, line := range wrapLines(title, contentWidth) {
				if y > contentBottom {
					truncated = true
					break
				}
				renderer.DrawText(screen, contentX, y, line, titleColor, theme.Background)
				y++
			}
			if truncated {
				break
			}

			for _, line := range wrapLines(item.Body, contentWidth) {
				if y > contentBottom {
					truncated = true
					break
				}
				renderer.DrawText(screen, contentX, y, line, theme.SettingsText, theme.Background)
				y++
			}
			if truncated {
				break
			}

			if item.Version != "" {
				if y > contentBottom {
					truncated = true
					break
				}
				versionLine := trimToWidth("Version: "+item.Version, contentWidth)
				renderer.DrawText(screen, contentX, y, versionLine, theme.SettingsHint, theme.Background)
				y++
			}

			y++
		}
	}

	if truncated {
		renderer.DrawTextCentered(screen, height-9, "Older notifications are hidden on smaller terminals.", theme.SettingsHint, theme.Background)
	}

	backText := "Back"
	backFg := theme.SettingsSelectedFg
	backBg := theme.SettingsSelectedBg
	if p.Selected != 2 {
		backFg = theme.SettingsText
		backBg = theme.Background
	}
	renderer.DrawTextCentered(screen, backY, backText, backFg, backBg)
}

func wrapLines(text string, width int) []string {
	if width <= 0 {
		return nil
	}
	parts := strings.Split(strings.TrimSpace(text), "\n")
	lines := make([]string, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line == "" {
			lines = append(lines, "")
			continue
		}
		for runeLen(line) > width {
			cut := lastSpaceWithin(line, width)
			if cut <= 0 {
				cut = width
			}
			chunk := strings.TrimSpace(string([]rune(line)[:cut]))
			if chunk != "" {
				lines = append(lines, chunk)
			}
			line = strings.TrimSpace(string([]rune(line)[cut:]))
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func trimToWidth(text string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

func lastSpaceWithin(text string, width int) int {
	runes := []rune(text)
	if len(runes) <= width {
		return len(runes)
	}
	for index := width; index > 0; index-- {
		if runes[index-1] == ' ' || runes[index-1] == '\t' {
			return index - 1
		}
	}
	return width
}

func runeLen(text string) int {
	return len([]rune(text))
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
