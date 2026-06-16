package manager

import (
	"time"

	"virga-player/debug/overlay"
	"virga-player/debug/storage"

	"github.com/gdamore/tcell/v2"
)

func (m *Manager) DrawOverlay(screen tcell.Screen) {
	if !m.Enabled() {
		return
	}
	lo, ok := overlay.BuildLayout(screen)
	if !ok {
		return
	}
	lines := m.buf.Last(20)

	fg := tcell.ColorWhite
	bg := tcell.ColorBlack
	muted := tcell.ColorSilver
	accent := tcell.ColorAqua
	errorColor := tcell.ColorRed

	overlay.DrawFrame(screen, lo, fg, bg)
	m.drawHeader(screen, lo, accent, bg)

	status, statusAt := m.statusSnapshot()
	contentTop := lo.Y1 + 2
	contentW := lo.W - 3
	if contentW < 8 {
		return
	}

	footer := "debug active"
	if status != "" && time.Since(statusAt) < 8*time.Second {
		footer = status
	}
	stats := m.runtimeSummary()
	tailLines := buildTailLines(footer, stats, contentW, accent)
	contentBottom := lo.Y2 - len(tailLines) - 1
	available := contentBottom - contentTop + 1
	if available < 1 {
		available = 0
	}
	renderLines := buildLogRenderLines(lines, contentW, muted, errorColor)
	if available > 0 && len(renderLines) > available {
		renderLines = renderLines[len(renderLines)-available:]
	}
	overlay.DrawLines(screen, lo.X1+2, contentTop, contentW, renderLines, bg)
	startTailY := lo.Y2 - len(tailLines)
	overlay.DrawLines(screen, lo.X1+2, startTailY, contentW, tailLines, bg)
}

func (m *Manager) statusSnapshot() (string, time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status, m.statusAt
}

func buildLogRenderLines(items []storage.Entry, width int, muted, errorColor tcell.Color) []overlay.Line {
	out := make([]overlay.Line, 0, len(items))
	for _, item := range items {
		full := item.At.Format("15:04:05") + " [" + item.Level + "] " + item.Text
		col := muted
		if item.Level == string(levelError) {
			col = errorColor
		} else if item.Level == string(levelWarn) {
			col = tcell.ColorOrange
		}
		for _, wrapped := range overlay.WrapText(full, width) {
			out = append(out, overlay.Line{Text: wrapped, Color: col})
		}
	}
	return out
}

func buildTailLines(footer, stats string, width int, accent tcell.Color) []overlay.Line {
	out := make([]overlay.Line, 0, 4)
	for _, s := range overlay.WrapText(footer, width) {
		out = append(out, overlay.Line{Text: s, Color: accent})
	}
	for _, s := range overlay.WrapText(stats, width) {
		out = append(out, overlay.Line{Text: s, Color: tcell.ColorLightCyan})
	}
	if len(out) > 4 {
		return out[len(out)-4:]
	}
	return out
}
