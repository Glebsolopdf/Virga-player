package artwork

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

func (a *Artwork) renderSixel(screen tcell.Screen) {
	theme := settings.CurrentTheme()
	if a.getCoverImg() == nil {
		a.renderTextOnly(screen)
		return
	}

	if a.SixelData == nil {
		if !a.prepareSixelData() {
			a.renderTextOnly(screen)
			return
		}
	}

	screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.TrackTitle))

	w, h := screen.Size()
	centerX := w / 2
	centerY := h / 2

	imgY := centerY - 12
	imgX := centerX - 20
	if imgX < 0 {
		imgX = 0
	}
	if imgY < 0 {
		imgY = 0
	}

	fmt.Printf("\x1B[%d;%dH", imgY+1, imgX+1)
	os.Stdout.Write(a.SixelData)

	a.drawText(screen, centerX-len(a.Title)/2, centerY+3, a.Title, theme.TrackTitle)
	a.drawText(screen, centerX-len(a.Artist)/2, centerY+5, a.Artist, theme.TrackArtist)
	a.drawText(screen, centerX-len(a.Album)/2, centerY+7, a.Album, theme.TrackAlbum)

	if a.Duration > 0 {
		a.drawTimeline(screen, centerX, centerY+9, 28)
		timeStr := formatTime(a.Elapsed) + " / " + formatTime(a.Duration)
		a.drawText(screen, centerX-len(timeStr)/2, centerY+11, timeStr, theme.TrackTime)
	}
}

func (a *Artwork) prepareSixelData() bool {
	img := a.getCoverImg()
	if img == nil {
		return false
	}

	var pngData bytes.Buffer
	if err := imageEncodePNG(&pngData, img); err != nil {
		return false
	}

	sixelCmd := exec.Command(
		"convert",
		"png:-",
		"-filter", "Lanczos",
		"-resize", "256x256^",
		"-gravity", "center",
		"-extent", "256x256",
		"sixel:-",
	)
	sixelCmd.Stdin = &pngData
	output, err := sixelCmd.Output()
	if err != nil || len(output) == 0 {
		return false
	}

	a.SixelData = output
	return true
}
