package julianday

import (
	"time"

	"github.com/shopspring/decimal"
)

// Convert converts a date to Julian Day
func Convert(date time.Time) decimal.Decimal {
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
