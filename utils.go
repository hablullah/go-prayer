package prayer

import (
	"time"

	"github.com/shopspring/decimal"
)

func getJulianDay(date time.Time) decimal.Decimal {
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
		return decimal.Zero
	}

	// If date is in blank days, stop
	endOfJulian := time.Date(1582, 10, 4, 23, 59, 59, 0, time.UTC)
	startOfGregorian := time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	if date.After(endOfJulian) && date.Before(startOfGregorian) {
		return decimal.Zero
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

	return decimal.NewFromFloat(1720994.5).
		Add(yearToDays).
		Add(monthToDays).
		Add(constant).
		Add(decimal.New(D, 0)).
		Add(timeToDays)
}

func getSunDeclination(jd decimal.Decimal) decimal.Decimal {
	T := decimal.New(2, 0).
		Mul(decPi).
		Mul(jd.Sub(decimal.New(2451545, 0))).
		Div(decimal.NewFromFloat(365.25))

	A := T.Mul(decimal.NewFromFloat(57.297)).
		Sub(decimal.NewFromFloat(79.547))

	B := T.Mul(decimal.New(2, 0)).
		Mul(decimal.NewFromFloat(57.297)).
		Sub(decimal.NewFromFloat(82.682))

	C := T.Mul(decimal.New(3, 0)).
		Mul(decimal.NewFromFloat(57.297)).
		Sub(decimal.NewFromFloat(59.722))

	return decimal.NewFromFloat(0.37877).
		Add(sin(A).Mul(decimal.NewFromFloat(23.264))).
		Add(sin(B).Mul(decimal.NewFromFloat(0.3812))).
		Add(sin(C).Mul(decimal.NewFromFloat(0.17132)))
}

func getEquationOfTime(jd decimal.Decimal) decimal.Decimal {
	U := jd.
		Sub(decimal.New(2451545, 0)).
		Div(decimal.New(36525, 0))

	L0 := U.
		Mul(decimal.NewFromFloat(36000.7698)).
		Add(decimal.NewFromFloat(280.46607))

	A := decimal.New(-1789, 0).
		Sub(U.Mul(decimal.New(237, 0))).
		Mul(sin(L0))

	B := decimal.New(7146, 0).
		Sub(U.Mul(decimal.New(62, 0))).
		Mul(cos(L0))

	C := decimal.New(9934, 0).
		Sub(U.Mul(decimal.New(14, 0))).
		Mul(sin(L0.Mul(decimal.New(2, 0))))

	D := decimal.New(29, 0).
		Add(U.Mul(decimal.New(5, 0))).
		Mul(cos(L0.Mul(decimal.New(2, 0))))

	E := decimal.New(74, 0).
		Add(U.Mul(decimal.New(10, 0))).
		Mul(sin(L0.Mul(decimal.New(3, 0))))

	F := decimal.New(320, 0).
		Sub(U.Mul(decimal.New(4, 0))).
		Mul(cos(L0.Mul(decimal.New(3, 0))))

	G := decimal.New(212, 0).
		Mul(sin(L0.Mul(decimal.New(4, 0))))

	return A.Sub(B).Add(C).Sub(D).Add(E).Add(F).Sub(G).
		Div(decimal.New(1000, 0))
}

func getTimezone(date time.Time) int64 {
	_, utcOffset := date.Zone()
	return decimal.New(int64(utcOffset), 0).
		Div(decimal.New(3600, 0)).
		IntPart()
}

func getTransitTime(timezone int64, longitude float64, eoT decimal.Decimal) decimal.Decimal {
	decTimezone := decimal.New(timezone, 0)
	decLongitude := decimal.NewFromFloat(longitude)

	return decimal.New(12, 0).
		Add(decTimezone).
		Sub(decLongitude.Div(decimal.New(15, 0))).
		Sub(eoT.Div(decimal.New(60, 0)))
}

func getHourAngle(latitude float64, sunAlt, sunDecli decimal.Decimal) decimal.Decimal {
	decLatitude := decimal.NewFromFloat(latitude)
	cosHourAngle := sin(sunAlt).
		Sub(sin(decLatitude).Mul(sin(sunDecli))).
		Div(cos(decLatitude).Mul(cos(sunDecli)))

	return acos(cosHourAngle)
}
