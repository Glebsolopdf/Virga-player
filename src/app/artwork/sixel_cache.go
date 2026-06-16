package artwork

func (a *Artwork) shouldRenderSixelAt(x, y int, data []byte) bool {
	dataSize := len(data)

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.sixelRendered && a.lastSixelX == x && a.lastSixelY == y && a.lastSixelSize == dataSize {
		return false
	}

	a.sixelRendered = true
	a.lastSixelX = x
	a.lastSixelY = y
	a.lastSixelSize = dataSize
	return true
}
