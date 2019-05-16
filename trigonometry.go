package prayer

import (
	"math"

	"github.com/shopspring/decimal"
)

var (
	decPi  = decimal.NewFromFloat(math.Pi)
	dec180 = decimal.New(180, 0)
)

func sin(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Sin(rad))
}

func cos(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Cos(rad))
}

func tan(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Tan(rad))
}

func acos(cosValue decimal.Decimal) decimal.Decimal {
	fl, _ := cosValue.Float64()
	return decimal.NewFromFloat(math.Acos(fl)).
		Mul(dec180).
		Div(decPi)
}

func acot(cotValue decimal.Decimal) decimal.Decimal {
	fl, _ := cotValue.Float64()
	return decPi.Div(decimal.New(2, 0)).
		Sub(decimal.NewFromFloat(math.Atan(fl))).
		Mul(dec180).
		Div(decPi)
}
