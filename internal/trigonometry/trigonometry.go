package trigonometry

import (
	"math"

	"github.com/shopspring/decimal"
)

var (
	decPi  = decimal.NewFromFloat(math.Pi)
	dec180 = decimal.New(180, 0)
)

// Sin calculate sin value from the specified degree.
func Sin(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Sin(rad))
}

// Cos calculate cos value from the specified degree.
func Cos(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Cos(rad))
}

// Tan calculate tan value from the specified degree.
func Tan(d decimal.Decimal) decimal.Decimal {
	rad, _ := d.Mul(decPi).Div(dec180).Float64()
	return decimal.NewFromFloat(math.Tan(rad))
}

// Acos calculate degree from specified cos value.
func Acos(cosValue decimal.Decimal) decimal.Decimal {
	fl, _ := cosValue.Float64()

	// Prevent NaN value
	if fl < -1 {
		fl = -1
	}
	if fl > 1 {
		fl = 1
	}

	return decimal.NewFromFloat(math.Acos(fl)).
		Mul(dec180).
		Div(decPi)
}

// Acot calculate degree from specified cotangent value.
func Acot(cotValue decimal.Decimal) decimal.Decimal {
	fl, _ := cotValue.Float64()
	return decPi.Div(decimal.New(2, 0)).
		Sub(decimal.NewFromFloat(math.Atan(fl))).
		Mul(dec180).
		Div(decPi)
}
