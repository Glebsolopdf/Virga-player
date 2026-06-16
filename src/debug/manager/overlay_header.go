package manager

import (
	"virga-player/debug/overlay"

	"github.com/gdamore/tcell/v2"
)

func (m *Manager) drawHeader(screen tcell.Screen, lo overlay.Layout, accent, bg tcell.Color) {
	overlay.DrawText(screen, lo.X1+2, lo.Y1, " DEBUG LOGS (last 20) ", accent, bg)

	copyLabel := "[Copy:C]"
	saveLabel := "[Save:K]"
	copyX := lo.X2 - len(copyLabel) - len(saveLabel) - 4
	saveX := lo.X2 - len(saveLabel) - 2
	if copyX < lo.X1+2 {
		copyX = lo.X1 + 2
		saveX = copyX + len(copyLabel) + 2
	}
	overlay.DrawText(screen, copyX, lo.Y1, copyLabel, tcell.ColorBlack, tcell.ColorGreen)
	overlay.DrawText(screen, saveX, lo.Y1, saveLabel, tcell.ColorBlack, tcell.ColorYellow)

	m.mu.Lock()
	m.copyBtn = rect{x1: copyX, y1: lo.Y1, x2: copyX + len(copyLabel) - 1, y2: lo.Y1}
	m.saveBtn = rect{x1: saveX, y1: lo.Y1, x2: saveX + len(saveLabel) - 1, y2: lo.Y1}
	m.lastOverlayW = lo.W
	m.mu.Unlock()
}
