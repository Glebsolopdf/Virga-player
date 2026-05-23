package app

import "github.com/gdamore/tcell/v2"

func (a *App) handleEvent(event tcell.Event) bool {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if a.debug != nil && a.debug.HandleEvent(ev) {
			return false
		}
		if a.settingsOpen {
			exit, save, deleteVirga := a.settingsPage.HandleKey(ev)
			if exit {
				return a.closeSettings(save, deleteVirga)
			}
			return false
		}
		if ev.Rune() == 's' {
			a.openSettings()
			return false
		}
		if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
			a.animEngine.Stop()
			return true
		}
	case *tcell.EventResize:
		a.width, a.height = a.screen.Size()
		a.particleSystem.Resize(a.width, a.height)
		a.state.Resize(a.width, a.height)
		a.screen.Sync()
	case *tcell.EventMouse:
		if a.debug != nil && a.debug.HandleEvent(ev) {
			return false
		}
	}
	return false
}
