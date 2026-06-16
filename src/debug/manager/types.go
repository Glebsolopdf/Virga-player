package manager

type level string

const (
	levelInfo  level = "INFO"
	levelWarn  level = "WARN"
	levelError level = "ERROR"
	levelDebug level = "DEBUG"
)

type rect struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func (r rect) contains(x, y int) bool {
	return x >= r.x1 && x <= r.x2 && y >= r.y1 && y <= r.y2
}
