package notification

import (
	"time"
	"unicode/utf8"

	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

type Toast struct {
	message   string
	expiresAt time.Time
	now       func() time.Time
}

func NewToast() *Toast {
	return &Toast{now: time.Now}
}

func (t *Toast) Show(message string, duration time.Duration) {
	t.message = message
	t.expiresAt = t.now().Add(duration)
}

func (t *Toast) Hide() {
	t.message = ""
	t.expiresAt = time.Time{}
}

func (t *Toast) Visible() bool {
	return t.message != "" && t.now().Before(t.expiresAt)
}

func (t *Toast) Render(screen tcell.Screen) {
	if !t.Visible() {
		return
	}

	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}

	theme := settings.CurrentTheme()
	style := tcell.StyleDefault.Foreground(theme.SettingsSelectedFg).Background(theme.SettingsSelectedBg)
	y := height - 1
	for x := 0; x < width; x++ {
		screen.SetContent(x, y, ' ', nil, style)
	}

	text := fitToWidth(t.message, width-2)
	x := 1
	if textWidth := utf8.RuneCountInString(text); textWidth < width {
		x = (width - textWidth) / 2
	}
	for offset, ch := range []rune(text) {
		screen.SetContent(x+offset, y, ch, nil, style)
	}
	if !t.now().Before(t.expiresAt) {
		t.Hide()
	}
}

func fitToWidth(text string, width int) string {
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
