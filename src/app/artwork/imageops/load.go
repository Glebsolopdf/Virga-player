package imageops

import (
	"context"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoadNormalized(imagePath string) (image.Image, error) {
	img, err := load(imagePath)
	if err != nil {
		return nil, err
	}
	return NormalizeSquare(trimPadding(img)), nil
}

func load(imagePath string) (image.Image, error) {
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		client := http.Client{Timeout: 4 * time.Second}
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, imagePath, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, io.EOF
		}
		img, _, err := image.Decode(io.LimitReader(resp.Body, 8*1024*1024))
		return img, err
	}

	f, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}
