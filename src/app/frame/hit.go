package frame

import "fmt"

func (f Frame) hitPlayer() {
	if f.Player == nil || !f.PlayerEnabled {
		return
	}

	titleY := f.Player.TextY
	artistY := f.Player.TextY + 1
	albumY := f.Player.TextY + 2
	durationY := f.Player.TextY + 3
	titleEndX := f.Player.TextX + len(f.Player.Title)
	artistEndX := f.Player.TextX + len(f.Player.Artist)
	albumEndX := f.Player.TextX + len(f.Player.Album)
	progress := fmt.Sprintf("%d:%02d / %d:%02d", f.Player.Elapsed/60, f.Player.Elapsed%60, f.Player.Duration/60, f.Player.Duration%60)
	durationEndX := f.Player.TextX + len(progress)

	for _, part := range f.ParticleSystem.GetParticles() {
		dropX := int(part.X)
		topY := int(part.Y)
		bottomY := topY + part.Length - 1

		// Mark full overlap against the rectangular cover area.
		if dropX >= f.Player.CoverX && dropX < f.Player.CoverX+f.Player.CoverW {
			overlapStart := topY
			if overlapStart < f.Player.CoverY {
				overlapStart = f.Player.CoverY
			}
			overlapEnd := bottomY
			coverBottom := f.Player.CoverY + f.Player.CoverH - 1
			if overlapEnd > coverBottom {
				overlapEnd = coverBottom
			}
			if overlapStart <= overlapEnd {
				for y := overlapStart; y <= overlapEnd; y++ {
					f.Player.MarkHit("cover", dropX, y)
				}
			}
		}

		if dropX >= f.Player.TextX && dropX < titleEndX && titleY >= topY && titleY <= bottomY {
			f.Player.MarkHit("title", dropX, titleY)
		}
		if dropX >= f.Player.TextX && dropX < artistEndX && artistY >= topY && artistY <= bottomY {
			f.Player.MarkHit("artist", dropX, artistY)
		}
		if dropX >= f.Player.TextX && dropX < albumEndX && albumY >= topY && albumY <= bottomY {
			f.Player.MarkHit("album", dropX, albumY)
		}
		if dropX >= f.Player.TextX && dropX < durationEndX && durationY >= topY && durationY <= bottomY {
			f.Player.MarkHit("duration", dropX, durationY)
		}
	}
}
