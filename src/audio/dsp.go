package audio

import "math"

func decodePCM16(buf []byte) []float64 {
	n := len(buf) / 2
	samples := make([]float64, 0, n)
	for i := 0; i+1 < len(buf); i += 2 {
		raw := int16(buf[i]) | int16(buf[i+1])<<8
		samples = append(samples, float64(raw)/32768.0)
	}
	return samples
}

func rms(samples []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range samples {
		sum += s * s
	}
	return math.Sqrt(sum / float64(len(samples)))
}

func smooth(prev, next, alpha float64) float64 {
	return prev + alpha*(next-prev)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
