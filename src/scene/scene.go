package scene

func NewScene(width, height int) *Scene {
	s := &Scene{
		width:  width,
		height: height,
	}
	s.generateBuildings()
	s.generateLamps()
	return s
}
