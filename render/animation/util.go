package animation

func Rescale(fromMin, fromMax, toMin, toMax, v float64) float64 {
	if fromMax == fromMin {
		return toMax
	}

	return toMin + (v-fromMin)/(fromMax-fromMin)*(toMax-toMin)
}

func Lerp(from, to, t float64) float64 {
	return from + (to-from)*t
}
