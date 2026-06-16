package welcome

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		m.frame++
		if m.animationsRunning() {
			return m, nextFrame()
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.result = &ConfigResult{Canceled: true}
			return m, tea.Quit
		case "up", "k":
			if m.stage == stageSelect && m.menuReady() {
				m.selected = (m.selected - 1 + len(m.items)) % len(m.items)
			}
		case "down", "j":
			if m.stage == stageSelect && m.menuReady() {
				m.selected = (m.selected + 1) % len(m.items)
			}
		case "enter":
			if m.menuReady() {
				if m.stage == stageSelect {
					m.chosen = m.items[m.selected].Config.Clone()
					m.stage = stageControls
					return m, nil
				}
				m.result = &ConfigResult{Config: m.chosen.Clone()}
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
