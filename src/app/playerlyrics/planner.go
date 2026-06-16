package playerlyrics

import "math/rand"

func nextPlacement(rng *rand.Rand, prev *dynamicLine, cue Cue, width, height int, player Rect, textW, textH int) Rect {
	if prev == nil {
		return randomRect(rng, width, height, textW, textH, player)
	}
	base := prev.Bounds
	candidates := placementOffsets(rng, effectiveGapMs(prev.Cue, cue))
	blocked := player.WithPadding(playerCollisionPadding)
	for _, shift := range candidates {
		rect := Rect{X: base.X + shift.X, Y: base.Y + shift.Y, W: textW, H: textH}
		if validPlacement(rect, width, height, blocked, base) {
			return rect
		}
	}
	return randomRect(rng, width, height, textW, textH, player)
}

type shift struct {
	X int
	Y int
}

func placementOffsets(rng *rand.Rand, gapMs int) []shift {
	steps := []int{2, 3}
	if gapMs >= 4000 {
		steps = []int{4, 5}
	}
	dirs := []shift{
		{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0},
		{X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: -1},
	}
	out := make([]shift, 0, len(steps)*len(dirs))
	for _, step := range steps {
		for _, dir := range dirs {
			out = append(out, shift{X: dir.X * step, Y: dir.Y * step})
		}
	}
	rng.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

func validPlacement(rect Rect, width, height int, blocked, prev Rect) bool {
	if rect.X < 0 || rect.Y < 0 || rect.X+rect.W > width || rect.Y+rect.H > height {
		return false
	}
	if rect.Intersects(blocked) {
		return false
	}
	return rect.X != prev.X || rect.Y != prev.Y
}

func effectiveGapMs(prev, current Cue) int {
	gapMs := current.AtMillis - prev.AtMillis
	if gapMs < 0 {
		return 0
	}
	if gapMs > 4500 {
		return 4500
	}
	return gapMs
}
