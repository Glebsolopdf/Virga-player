package overlay

import "github.com/gdamore/tcell/v2"

type Layout struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
	W  int
	H  int
}

type Line struct {
	Text  string
	Color tcell.Color
}
