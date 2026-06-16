package manager

import "github.com/gdamore/tcell/v2"

func (m *Manager) HandleEvent(ev tcell.Event) bool {
	if !m.Enabled() {
		return false
	}
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Rune() {
		case 'c', 'C':
			m.copy()
			return true
		case 'k', 'K':
			m.save()
			return true
		}
	case *tcell.EventMouse:
		x, y := e.Position()
		if e.Buttons()&tcell.Button1 == 0 {
			return false
		}
		m.mu.RLock()
		copyBtn := m.copyBtn
		saveBtn := m.saveBtn
		m.mu.RUnlock()
		if copyBtn.contains(x, y) {
			m.copy()
			return true
		}
		if saveBtn.contains(x, y) {
			m.save()
			return true
		}
	}
	return false
}
