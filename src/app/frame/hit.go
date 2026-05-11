package frame

import "fmt"

func (f Frame) hitPlayer() {
	if f.Player == nil || !f.PlayerEnabled {
		return
	}

	for _, part := range f.ParticleSystem.GetParticles() {
		for i := 0; i < part.Length; i++ {
			dropX := int(part.X)
			dropY := int(part.Y) + i

			// Hit album cover
			if dropX >= f.Player.CoverX && dropX < f.Player.CoverX+f.Player.CoverW &&
				dropY >= f.Player.CoverY && dropY < f.Player.CoverY+f.Player.CoverH {
				f.Player.MarkHit("cover", dropX, dropY)
			}

			// Hit title
			titleY := f.Player.TextY
			if dropY == titleY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Title) {
				f.Player.MarkHit("title", dropX, titleY)
			}

			// Hit artist
			artistY := f.Player.TextY + 1
			if dropY == artistY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Artist) {
				f.Player.MarkHit("artist", dropX, artistY)
			}

			// Hit album
			albumY := f.Player.TextY + 2
			if dropY == albumY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(f.Player.Album) {
				f.Player.MarkHit("album", dropX, albumY)
			}

			// Hit duration
			durationY := f.Player.TextY + 3
			progress := fmt.Sprintf("%d:%02d / %d:%02d", f.Player.Elapsed/60, f.Player.Elapsed%60, f.Player.Duration/60, f.Player.Duration%60)
			if dropY == durationY && dropX >= f.Player.TextX && dropX < f.Player.TextX+len(progress) {
				f.Player.MarkHit("duration", dropX, durationY)
			}
		}
	}
}
