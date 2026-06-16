package playerlyrics

import "virga-player/settings"

const playerCollisionPadding = 2

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Point struct {
	X int
	Y int
}

func Anchor(width, height, itemW, itemH int, preset settings.PositionPreset) Point {
	marginX := 4
	marginY := 2
	switch preset {
	case settings.PositionTopRight:
		return Point{X: width - itemW - marginX, Y: marginY}
	case settings.PositionBottomRight:
		return Point{X: width - itemW - marginX, Y: height - itemH - marginY}
	case settings.PositionTopCenter:
		return Point{X: (width - itemW) / 2, Y: marginY}
	default:
		return Point{X: (width - itemW) / 2, Y: (height - itemH) / 2}
	}
}

func StaticRect(width, height, itemW, itemH int, player Rect, preset settings.PositionPreset) Rect {
	if preset == settings.PositionBelowPlayer {
		x := player.X + (player.W-itemW)/2
		y := player.Y + player.H + 2
		return Rect{X: x, Y: y, W: itemW, H: itemH}.Clamp(width, height)
	}
	pt := Anchor(width, height, itemW, itemH, preset)
	return Rect{X: pt.X, Y: pt.Y, W: itemW, H: itemH}.Clamp(width, height)
}

func PlayerBounds(width, height int, preset settings.PositionPreset) Rect {
	body := Rect{W: 56, H: 24}
	pt := Anchor(width, height, body.W, body.H, preset)
	body.X = pt.X
	body.Y = pt.Y
	return body.Clamp(width, height)
}

func (r Rect) WithPadding(pad int) Rect {
	return Rect{X: r.X - pad, Y: r.Y - pad, W: r.W + pad*2, H: r.H + pad*2}
}

func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.W &&
		r.X+r.W > other.X &&
		r.Y < other.Y+other.H &&
		r.Y+r.H > other.Y
}

func (r Rect) Clamp(width, height int) Rect {
	if r.W > width {
		r.W = width
	}
	if r.H > height {
		r.H = height
	}
	if r.X < 0 {
		r.X = 0
	}
	if r.Y < 0 {
		r.Y = 0
	}
	if r.X+r.W > width {
		r.X = width - r.W
	}
	if r.Y+r.H > height {
		r.Y = height - r.H
	}
	return r
}
