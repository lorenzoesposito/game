package utils

import "strconv"

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

func StringToFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	}
	return -1
}
