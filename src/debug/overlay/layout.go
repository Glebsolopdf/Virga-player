package overlay

import "github.com/gdamore/tcell/v2"

func BuildLayout(screen tcell.Screen) (Layout, bool) {
	w, h := screen.Size()
	if w < 50 || h < 12 {
		return Layout{}, false
	}
	boxW := 64
	if boxW > w-2 {
		boxW = w - 2
	}
	boxH := 24
	if boxH > h-2 {
		boxH = h - 2
	}
	x1 := w - boxW - 1
	y1 := 1
	return Layout{X1: x1, Y1: y1, X2: x1 + boxW, Y2: y1 + boxH, W: boxW, H: boxH}, true
}

func DrawFrame(screen tcell.Screen, lo Layout, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for y := lo.Y1; y <= lo.Y2; y++ {
		for x := lo.X1; x <= lo.X2; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
	for x := lo.X1; x <= lo.X2; x++ {
		screen.SetContent(x, lo.Y1, '-', nil, style)
		screen.SetContent(x, lo.Y2, '-', nil, style)
	}
	for y := lo.Y1; y <= lo.Y2; y++ {
		screen.SetContent(lo.X1, y, '|', nil, style)
		screen.SetContent(lo.X2, y, '|', nil, style)
	}
	screen.SetContent(lo.X1, lo.Y1, '+', nil, style)
	screen.SetContent(lo.X2, lo.Y1, '+', nil, style)
	screen.SetContent(lo.X1, lo.Y2, '+', nil, style)
	screen.SetContent(lo.X2, lo.Y2, '+', nil, style)
}
