package separating_frequencies

func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func LayerEnergy(enabled bool, layer int, pulse, low, mid, high float64) float64 {
	if !enabled {
		return Clamp(high*0.64+mid*0.24+low*0.08+pulse*0.04, 0, 1)
	}

	switch layer {
	case 0:
		return Clamp(low*0.90+mid*0.07+high*0.02+pulse*0.04, 0, 1)
	case 1:
		return Clamp(low*0.08+mid*0.84+high*0.08+pulse*0.03, 0, 1)
	case 2:
		return Clamp(low*0.03+mid*0.20+high*0.76+pulse*0.03, 0, 1)
	default:
		return Clamp(low*0.02+mid*0.16+high*0.79+pulse*0.03, 0, 1)
	}
}

func LayerSpeedEnergy(enabled bool, layer int, pulse, low, mid, high float64) float64 {
	if !enabled {
		return Clamp(high*0.42+mid*0.32+low*0.20+pulse*0.06, 0, 1)
	}

	switch layer {
	case 0:
		return Clamp(low*0.93+mid*0.04+high*0.03+pulse*0.05, 0, 1)
	case 1:
		return Clamp(low*0.05+mid*0.84+high*0.10+pulse*0.04, 0, 1)
	case 2:
		return Clamp(low*0.02+mid*0.16+high*0.78+pulse*0.04, 0, 1)
	default:
		return Clamp(low*0.01+mid*0.12+high*0.83+pulse*0.04, 0, 1)
	}
}
