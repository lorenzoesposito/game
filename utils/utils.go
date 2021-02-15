package utils

func GetAxis(neg, pos bool) float64 {
	ax := float64(0)
	if neg {
		ax--
	}
	if pos {
		ax++
	}
	return ax
}
