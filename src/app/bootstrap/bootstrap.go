package bootstrap

import "github.com/gdamore/tcell/v2"

func NewScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset))
	screen.Clear()
	return screen, nil
}
