package prayer

import (
	"math"
	"time"

	trig "github.com/RadhiFadlillah/go-prayer/internal/trigonometry"
	"github.com/shopspring/decimal"
)

// TimeCalculator is calculator that used to calculate the prayer times.
type TimeCalculator struct {
	Elevation        float64
	Latitude         decimal.Decimal
	Longitude        decimal.Decimal
	AsrCoefficient   decimal.Decimal
	FajrAngle        decimal.Decimal
	IshaAngle        decimal.Decimal
	MaghribDuration  time.Duration
	PreciseToSeconds bool

	date           time.Time
	transitTime    decimal.Decimal
	sunDeclination decimal.Decimal
}

// SetDate set current date.
func (tc *TimeCalculator) SetDate(date time.Time) {
	// Make sure date is at 12 local time
	tc.date = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		12, 0, 0, 0,
		date.Location())

	// Calculate Julian Day
	jd := getJulianDay(tc.date)

	// Calculate timezone
	_, utcOffset := tc.date.Zone()
	timezone := decimal.New(int64(utcOffset), 0).
		Div(decimal.New(3600, 0))

	// Calculate transit time
	eoT := tc.getEquationOfTime(jd)
	tc.transitTime = decimal.New(12, 0).
		Add(timezone).
		Sub(tc.Longitude.Div(decimal.New(15, 0))).
		Sub(eoT.Div(decimal.New(60, 0)))

	// Calculate sun declination
	tc.sunDeclination = tc.getSunDeclination(jd)
}

// GetFajrTime returns time of Fajr.
func (tc TimeCalculator) GetFajrTime() time.Time {
	sunAltitude := tc.FajrAngle.Neg()
	hourAngle := tc.getHourAngle(sunAltitude)
	hourAngle = hourAngle.Div(decimal.New(15, 0))

	hours := tc.transitTime.Sub(hourAngle)
	return tc.hoursToTime(hours)
}

// GetSunriseTime returns time of sunrise.
func (tc TimeCalculator) GetSunriseTime() time.Time {
	sqrtElevation := decimal.NewFromFloat(math.Sqrt(tc.Elevation))
	A := decimal.New(-5, 0).Div(decimal.New(6, 0)) // -0.833333
	B := decimal.NewFromFloat(0.0347).Mul(sqrtElevation)
	sunAltitude := A.Sub(B)

	hourAngle := tc.getHourAngle(sunAltitude)
	hourAngle = hourAngle.Div(decimal.New(15, 0))

	hours := tc.transitTime.Sub(hourAngle)
	return tc.hoursToTime(hours)
}

// GetZuhrTime returns time of Zuhr.
func (tc TimeCalculator) GetZuhrTime() time.Time {
	twoMinutes := decimal.New(2, 0).Div(decimal.New(60, 0))
	hours := tc.transitTime.Add(twoMinutes)
	return tc.hoursToTime(hours)
}

// GetAsrTime returns time of Asr.
func (tc TimeCalculator) GetAsrTime() time.Time {
	A := trig.Tan(tc.sunDeclination.Sub(tc.Latitude).Abs())
	sunAltitude := trig.Acot(tc.AsrCoefficient.Add(A))

	hourAngle := tc.getHourAngle(sunAltitude)
	hourAngle = hourAngle.Div(decimal.New(15, 0))

	hours := tc.transitTime.Add(hourAngle)
	return tc.hoursToTime(hours)
}

// GetMaghribTime returns time of Maghrib.
func (tc TimeCalculator) GetMaghribTime() time.Time {
	sqrtElevation := decimal.NewFromFloat(math.Sqrt(tc.Elevation))
	A := decimal.New(-5, 0).Div(decimal.New(6, 0)) // -0.833333
	B := decimal.NewFromFloat(0.0347).Mul(sqrtElevation)
	sunAltitude := A.Sub(B)

	hourAngle := tc.getHourAngle(sunAltitude)
	hourAngle = hourAngle.Div(decimal.New(15, 0))

	hours := tc.transitTime.Add(hourAngle)
	return tc.hoursToTime(hours)
}

// GetIshaTime returns time of Isha.
func (tc TimeCalculator) GetIshaTime() time.Time {
	if tc.MaghribDuration != 0 {
		maghrib := tc.GetMaghribTime()
		return maghrib.Add(tc.MaghribDuration)
	}

	sunAltitude := tc.IshaAngle.Neg()
	hourAngle := tc.getHourAngle(sunAltitude)
	hourAngle = hourAngle.Div(decimal.New(15, 0))

	hours := tc.transitTime.Add(hourAngle)
	return tc.hoursToTime(hours)
}

func (tc TimeCalculator) getHourAngle(sunAltitude decimal.Decimal) decimal.Decimal {
	sinSunAltitude := trig.Sin(sunAltitude)
	sinLatitude := trig.Sin(tc.Latitude)
	cosLatitude := trig.Cos(tc.Latitude)
	sinSunDeclination := trig.Sin(tc.sunDeclination)
	cosSunDeclination := trig.Cos(tc.sunDeclination)

	cosHourAngle := sinSunAltitude.
		Sub(sinLatitude.Mul(sinSunDeclination)).
		Div(cosLatitude.Mul(cosSunDeclination))

	return trig.Acos(cosHourAngle)
}

func (tc TimeCalculator) getEquationOfTime(jd decimal.Decimal) decimal.Decimal {
	U := jd.
		Sub(decimal.New(2451545, 0)).
		Div(decimal.New(36525, 0))

	L0 := U.
		Mul(decimal.NewFromFloat(36000.7698)).
		Add(decimal.NewFromFloat(280.46607))

	A := decimal.New(-1789, 0).
		Sub(U.Mul(decimal.New(237, 0))).
		Mul(trig.Sin(L0))

	B := decimal.New(7146, 0).
		Sub(U.Mul(decimal.New(62, 0))).
		Mul(trig.Cos(L0))

	C := decimal.New(9934, 0).
		Sub(U.Mul(decimal.New(14, 0))).
		Mul(trig.Sin(L0.Mul(decimal.New(2, 0))))

	D := decimal.New(29, 0).
		Add(U.Mul(decimal.New(5, 0))).
		Mul(trig.Cos(L0.Mul(decimal.New(2, 0))))

	E := decimal.New(74, 0).
		Add(U.Mul(decimal.New(10, 0))).
		Mul(trig.Sin(L0.Mul(decimal.New(3, 0))))

	F := decimal.New(320, 0).
		Sub(U.Mul(decimal.New(4, 0))).
		Mul(trig.Cos(L0.Mul(decimal.New(3, 0))))

	G := decimal.New(212, 0).
		Mul(trig.Sin(L0.Mul(decimal.New(4, 0))))

	return A.Sub(B).Add(C).Sub(D).Add(E).Add(F).Sub(G).
		Div(decimal.New(1000, 0))
}

func (tc TimeCalculator) getSunDeclination(jd decimal.Decimal) decimal.Decimal {
	decPi := decimal.NewFromFloat(math.Pi)

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
		Add(trig.Sin(A).Mul(decimal.NewFromFloat(23.264))).
		Add(trig.Sin(B).Mul(decimal.NewFromFloat(0.3812))).
		Add(trig.Sin(C).Mul(decimal.NewFromFloat(0.17132)))
}

func (tc TimeCalculator) hoursToTime(hours decimal.Decimal) time.Time {
	var minutes, seconds int

	if tc.PreciseToSeconds {
		seconds = int(hours.Mul(decimal.New(3600, 0)).Ceil().IntPart())
	} else {
		minutes = int(hours.Mul(decimal.New(60, 0)).Ceil().IntPart())
	}

	return time.Date(
		tc.date.Year(),
		tc.date.Month(),
		tc.date.Day(),
		0, minutes, seconds, 0,
		tc.date.Location())
}
