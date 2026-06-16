package sixeldata

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"os/exec"
	"time"
)

const buildTimeout = 2 * time.Second

func Build(img image.Image) ([]byte, bool) {
	if img == nil {
		return nil, false
	}

	var pngData bytes.Buffer
	if err := png.Encode(&pngData, img); err != nil {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	defer cancel()

	sixelCmd := exec.CommandContext(
		ctx,
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
	if ctx.Err() != nil || err != nil || len(output) == 0 {
		return nil, false
	}
	return output, true
}
