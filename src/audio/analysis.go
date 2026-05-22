package audio

import "math"

func bandEnergy(samples []float64, freqs []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	sum := 0.0
	for _, f := range freqs {
		sum += goertzelPower(samples, sampleRate, f)
	}
	return sum / float64(len(freqs))
}

func goertzelPower(samples []float64, rate int, freq float64) float64 {
	n := len(samples)
	if n == 0 {
		return 0
	}
	k := int(0.5 + (float64(n)*freq)/float64(rate))
	omega := 2.0 * math.Pi * float64(k) / float64(n)
	coeff := 2.0 * math.Cos(omega)
	s0, s1, s2 := 0.0, 0.0, 0.0
	for _, x := range samples {
		s0 = x + coeff*s1 - s2
		s2 = s1
		s1 = s0
	}
	power := s1*s1 + s2*s2 - coeff*s1*s2
	if power < 0 {
		return 0
	}
	return power / float64(n)
}
