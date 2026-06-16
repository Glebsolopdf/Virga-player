package artwork

import (
	imageops "virga-player/app/artwork/imageops"
	textrender "virga-player/app/artwork/textrender"
)

func (a *Artwork) loadCoverImage() {
	cov, err := imageops.LoadNormalized(a.ImagePath)
	if err != nil {
		return
	}
	a.mu.Lock()
	a.CoverImg = cov
	a.textRenderCache = textrender.Cache{}
	a.mu.Unlock()
}
