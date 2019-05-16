package prayer

import (
	"math"
	"time"

	"github.com/shopspring/decimal"
)

func getJulianDay(date time.Time) float64 {
	// Convert to UTC
	date = date.UTC()

	// Prepare variables for calculating
	Y := int64(date.Year())
	M := int64(date.Month())
	D := int64(date.Day())
	H := int64(date.Hour())
	m := int64(date.Minute())
	s := int64(date.Second())

	// If year is before 4713 B.C, stop
	if Y < -4712 {
		return 0
	}

	// If date is in blank days, stop
	endOfJulian := time.Date(1582, 10, 4, 23, 59, 59, 0, time.UTC)
	startOfGregorian := time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	if date.After(endOfJulian) && date.Before(startOfGregorian) {
		return 0
	}

	// If month <= 2, change year and month
	if M <= 2 {
		M += 12
		Y--
	}

	// Check whether date is gregorian or julian
	constant := decimal.Zero
	if date.After(endOfJulian) {
		temp := decimal.New(Y, -2).Floor()
		constant = decimal.New(2, 0).
			Add(temp.Div(decimal.New(4, 0)).Floor()).
			Sub(temp)
	}

	// Calculate julian day
	yearToDays := decimal.New(Y, 0).
		Mul(decimal.NewFromFloat(365.25)).
		Floor()

	monthToDays := decimal.New(M+1, 0).
		Mul(decimal.NewFromFloat(30.6001)).
		Floor()

	timeToSeconds := H*3600 + m*60 + s
	timeToDays := decimal.New(timeToSeconds, 0).
		Div(decimal.New(86400, 0))

	jd, _ := decimal.NewFromFloat(1720994.5).
		Add(yearToDays).
		Add(monthToDays).
		Add(constant).
		Add(decimal.New(D, 0)).
		Add(timeToDays).
		Float64()

	return jd
}

func getSunDeclination(jd float64) float64 {
	T := 2 * math.Pi * (jd - 2451545.0) / 365.25

	return 0.37877 +
		23.264*sin(57.297*T-79.547) +
		0.3812*sin(2*57.297*T-82.682) +
		0.17132*sin(3*57.297*T-59.722)
}

func getEquationOfTime(jd float64) float64 {
	U := (jd - 2451545.0) / 36525.0
	L0 := 280.46607 + 36000.7698*U

	return (-(1789+237*U)*sin(L0) -
		(7146-62*U)*cos(L0) +
		(9934-14*U)*sin(2*L0) -
		(29+5*U)*cos(2*L0) +
		(74+10*U)*sin(3*L0) +
		(320-4*U)*cos(3*L0) -
		212*sin(4*L0)) / 1000
}

func sin(d float64) float64 {
	return math.Sin(d * math.Pi / 180.0)
}

func cos(d float64) float64 {
	return math.Cos(d * math.Pi / 180.0)
}

func tan(d float64) float64 {
	return math.Tan(d * math.Pi / 180.0)
}

func acos(d float64) float64 {
	return math.Acos(d) * 180.0 / math.Pi
}

func arccot(d float64) float64 {
	return (math.Pi/2.0 - math.Atan(d)) * 180.0 / math.Pi
}

func round(d float64, decimalPlaces int) float64 {
	return math.Floor(d*math.Pow10(decimalPlaces)+0.5) /
		math.Pow10(decimalPlaces)
}
