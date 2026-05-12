package scene

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

func (s *Scene) Resize(width, height int) {
	s.width = width
	s.height = height
	s.generateBuildings()
	s.generateLamps()
}
