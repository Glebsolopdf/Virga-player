package scene

type Building struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Lamp struct {
	X int
	Y int
}

type Scene struct {
	width     int
	height    int
	buildings []Building
	lamps     []Lamp
}
