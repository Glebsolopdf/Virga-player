package artwork

import sixeldata "virga-player/app/artwork/sixeldata"

func loadFailed(imagePath string) bool {
	artworkLoadFailuresMu.RLock()
	defer artworkLoadFailuresMu.RUnlock()
	_, ok := artworkLoadFailures[imagePath]
	return ok
}

func markLoadFailed(imagePath string) {
	artworkLoadFailuresMu.Lock()
	artworkLoadFailures[imagePath] = struct{}{}
	artworkLoadFailuresMu.Unlock()
}

func (a *Artwork) beginCoverLoad() {
	if a.ImagePath == "" {
		return
	}

	if loadFailed(a.ImagePath) {
		return
	}

	a.mu.Lock()
	if a.loadStarted {
		a.mu.Unlock()
		return
	}
	a.loadStarted = true
	a.mu.Unlock()

	go func() {
		a.loadCoverImage()
		a.mu.RLock()
		img := a.CoverImg
		a.mu.RUnlock()
		if img == nil {
			markLoadFailed(a.ImagePath)
			return
		}
		a.computeAverageColor()
		if DetectSixelSupport() {
			if cached, ok := sixeldata.Get(a.ImagePath); ok {
				a.mu.Lock()
				a.SixelData = cached
				a.Mode = DisplaySixel
				a.mu.Unlock()
				return
			}
			a.mu.Lock()
			a.Mode = DisplayTextOnly
			a.mu.Unlock()
			a.prepareSixelDataAsync()
			return
		}

		a.mu.Lock()
		a.Mode = DisplayTextOnly
		a.mu.Unlock()
	}()
}
