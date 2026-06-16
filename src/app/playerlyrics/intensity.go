package playerlyrics

import "virga-player/audio"

func Intensity(bands audio.Bands) float64 {
	value := bands.Envelope*0.60 + bands.Low*0.30 + bands.Mid*0.20 + bands.High*0.10
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func appearFrames(pressure float64) int {
	frames := int(14 - pressure*12)
	if frames < 1 {
		return 1
	}
	return frames
}

func washFrames(pressure float64) int {
	frames := int(28 - pressure*12)
	if frames < 10 {
		return 10
	}
	return frames
}
