package scene

import (
	"github.com/gdamore/tcell/v2"
)

type Building struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Scene struct {
	width     int
	height    int
	buildings []Building
	lamps     []Lamp
}

type Lamp struct {
	X int
	Y int
}

func NewScene(width, height int) *Scene {
	s := &Scene{
		width:  width,
		height: height,
	}
	s.generateBuildings()
	s.generateLamps()
	return s
}

func (s *Scene) generateBuildings() {
	s.buildings = []Building{}
	buildingWidth := 12
	spacing := 3
	x := 0

	for x < s.width {
		height := 5 + int(x%3)*2
		y := s.height - height - 1
		s.buildings = append(s.buildings, Building{
			X:      x,
			Y:      y,
			Width:  buildingWidth,
			Height: height,
		})
		x += buildingWidth + spacing
	}
}
func (s *Scene) generateLamps() {
	s.lamps = []Lamp{}
	spacing := 20
	lampY := s.height - 2

	for x := 10; x < s.width; x += spacing {
		s.lamps = append(s.lamps, Lamp{X: x, Y: lampY})
	}
}

func (s *Scene) Draw(screen tcell.Screen) {
	for _, building := range s.buildings {
		s.drawBuilding(screen, building)
	}

	for x := 0; x < s.width; x++ {
		screen.SetContent(x, s.height-1, '─', nil, tcell.StyleDefault.
			Foreground(tcell.ColorGray).
			Background(tcell.ColorBlack))
	}

	for _, lamp := range s.lamps {
		s.drawLamp(screen, lamp)
	}
}

func (s *Scene) drawBuilding(screen tcell.Screen, b Building) {
	for y := b.Y; y < b.Y+b.Height; y++ {
		for x := b.X; x < b.X+b.Width; x++ {
			if y == b.Y || y == b.Y+b.Height-1 {
				screen.SetContent(x, y, '─', nil, tcell.StyleDefault.
					Foreground(tcell.ColorDarkGray).
					Background(tcell.ColorBlack))
			} else if x == b.X || x == b.X+b.Width-1 {
				screen.SetContent(x, y, '│', nil, tcell.StyleDefault.
					Foreground(tcell.ColorDarkGray).
					Background(tcell.ColorBlack))
			} else {
				screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.
					Background(tcell.ColorBlack))
			}
		}
	}

	for y := b.Y + 1; y < b.Y+b.Height-1; y += 2 {
		for x := b.X + 2; x < b.X+b.Width-1; x += 3 {
			if (x+y)%3 == 0 {
				screen.SetContent(x, y, '█', nil, tcell.StyleDefault.
					Foreground(tcell.ColorMaroon).
					Background(tcell.ColorBlack))
			} else {
				screen.SetContent(x, y, '·', nil, tcell.StyleDefault.
					Foreground(tcell.ColorGray).
					Background(tcell.ColorBlack))
			}
		}
	}
}

func (s *Scene) drawLamp(screen tcell.Screen, l Lamp) {

	screen.SetContent(l.X, l.Y-2, '│', nil, tcell.StyleDefault.
		Foreground(tcell.ColorGray).
		Background(tcell.ColorBlack))
	screen.SetContent(l.X, l.Y-1, '│', nil, tcell.StyleDefault.
		Foreground(tcell.ColorGray).
		Background(tcell.ColorBlack))

	screen.SetContent(l.X, l.Y-3, '◆', nil, tcell.StyleDefault.
		Foreground(tcell.ColorMaroon).
		Background(tcell.ColorBlack))

	if l.X-1 >= 0 {
		screen.SetContent(l.X-1, l.Y-3, '·', nil, tcell.StyleDefault.
			Foreground(tcell.ColorGray).
			Background(tcell.ColorBlack))
	}
	if l.X+1 < s.width {
		screen.SetContent(l.X+1, l.Y-3, '·', nil, tcell.StyleDefault.
			Foreground(tcell.ColorGray).
			Background(tcell.ColorBlack))
	}
}

func (s *Scene) Resize(width, height int) {
	s.width = width
	s.height = height
	s.generateBuildings()
	s.generateLamps()
}
