package prayer

import "math"

func sin(degree float64) float64 {
	return math.Sin(rad(degree))
}

func cos(degree float64) float64 {
	return math.Cos(rad(degree))
}

func tan(degree float64) float64 {
	return math.Tan(rad(degree))
}

func asin(sinValue float64) float64 {
	return degree(math.Asin(sinValue))
}

func acos(cosValue float64) float64 {
	return degree(math.Acos(cosValue))
}

func atan2(y, x float64) float64 {
	return degree(math.Atan2(y, x))
}

func acot(cotValue float64) float64 {
	return degree(math.Atan(1 / cotValue))
}

func rad(degree float64) float64 {
	return degree * math.Pi / 180
}

func degree(rad float64) float64 {
	return rad * 180 / math.Pi
}
