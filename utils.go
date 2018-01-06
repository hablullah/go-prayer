package prayer

import (
	"math"
	"time"
)

func getJulianDay(date time.Time) float64 {
	date = date.UTC()
	second := date.Second()
	minute := date.Minute()
	hour := date.Hour()
	day := float64(date.Day()) + float64(hour*3600+minute*60+second)/86400.0
	month := int(date.Month())
	year := date.Year()

	// If year is before 4713 B.C, stop
	if year < -4712 {
		return 0
	}

	// If month <= 2, change year and month
	if month <= 2 {
		month += 12
		year--
	}

	// Check whether date is gregorian or julian
	var constant float64
	julianChanged := time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	if date.After(julianChanged) {
		temp := math.Floor(float64(year) / 100)
		constant = 2.0 + math.Floor(temp/4) - temp
	}

	// Calculate julian day
	julianDay := 1720994.5 +
		math.Floor(float64(year)*365.25) +
		math.Floor(float64(month+1)*30.6001) +
		constant + float64(day)

	return julianDay
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
