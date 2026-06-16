package playerlyrics

import "math/rand"

func randomRect(rng *rand.Rand, width, height, textWidth, textHeight int, player Rect) Rect {
	rect := Rect{W: clampInt(textWidth, 8, maxInt(8, width-2)), H: clampInt(textHeight, 1, maxInt(1, height-1))}
	if height < 3 {
		rect.X = maxInt(0, (width-rect.W)/2)
		return rect
	}

	blocked := player.WithPadding(playerCollisionPadding)
	for attempt := 0; attempt < 48; attempt++ {
		rect.X = rng.Intn(maxInt(1, width-rect.W))
		rect.Y = rng.Intn(maxInt(1, height-rect.H))
		if !rect.Intersects(blocked) {
			return rect.Clamp(width, height)
		}
	}
	return fallbackRect(width, height, rect, blocked)
}

func fallbackRect(width, height int, rect, blocked Rect) Rect {
	positions := []Rect{
		{X: 1, Y: 1, W: rect.W, H: rect.H},
		{X: width - rect.W - 1, Y: 1, W: rect.W, H: rect.H},
		{X: 1, Y: height - rect.H - 1, W: rect.W, H: rect.H},
		{X: width - rect.W - 1, Y: height - rect.H - 1, W: rect.W, H: rect.H},
	}
	for _, candidate := range positions {
		candidate = candidate.Clamp(width, height)
		if !candidate.Intersects(blocked) {
			return candidate
		}
	}
	return Rect{X: 0, Y: 0, W: rect.W, H: rect.H}.Clamp(width, height)
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
