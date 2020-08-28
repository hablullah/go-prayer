package prayer

import (
	"math"
	"time"

	trig "github.com/RadhiFadlillah/go-prayer/internal/trigonometry"
	"github.com/shopspring/decimal"
)

func (calc Calculator) getHourAngle(sunAltitude, sunDeclination decimal.Decimal) decimal.Decimal {
	sinSunAltitude := trig.Sin(sunAltitude)
	sinLatitude := trig.Sin(calc.latitude)
	cosLatitude := trig.Cos(calc.latitude)
	sinSunDeclination := trig.Sin(sunDeclination)
	cosSunDeclination := trig.Cos(sunDeclination)

	cosHourAngle := sinSunAltitude.
		Sub(sinLatitude.Mul(sinSunDeclination)).
		Div(cosLatitude.Mul(cosSunDeclination))

	return trig.Acos(cosHourAngle)
}

func (calc Calculator) getEquationOfTime(jd decimal.Decimal) decimal.Decimal {
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

func (calc Calculator) getTransitTime(jd decimal.Decimal) decimal.Decimal {
	eoT := calc.getEquationOfTime(jd)
	return decimal.New(12, 0).
		Add(calc.timezone).
		Sub(calc.longitude.Div(decimal.New(15, 0))).
		Sub(eoT.Div(decimal.New(60, 0)))
}

func (calc Calculator) getSunDeclination(jd decimal.Decimal) decimal.Decimal {
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

func (calc Calculator) getSunAltitude(target Target, jd decimal.Decimal) decimal.Decimal {
	switch target {
	case Fajr:
		return calc.fajrAngle.Neg()

	case Sunrise, Maghrib:
		sqrtElevation := decimal.NewFromFloat(math.Sqrt(calc.Elevation))
		A := decimal.New(-5, 0).Div(decimal.New(6, 0)) // -0.833333
		B := decimal.NewFromFloat(0.0347).Mul(sqrtElevation)
		if calc.IgnoreElevation {
			return A
		} else {
			return A.Sub(B)
		}

	case Asr:
		sunDeclination := calc.getSunDeclination(jd)
		A := trig.Tan(sunDeclination.Sub(calc.latitude).Abs())
		B := calc.asrCoefficient.Add(A)
		return trig.Acot(B)

	case Isha:
		return calc.ishaAngle.Neg()

	default:
		return decimal.Zero
	}
}

func (calc Calculator) hoursToTime(hours decimal.Decimal) time.Time {
	y := calc.date.Year()
	m := calc.date.Month()
	d := calc.date.Day()
	location := calc.date.Location()
	seconds := int(hours.Mul(decimal.New(3600, 0)).Ceil().IntPart())
	newTime := time.Date(y, m, d, 0, 0, seconds, 0, location)

	if !calc.PreciseToSeconds {
		if newTime.Second() >= 30 {
			newTime = newTime.Add(time.Minute)
		}

		second := time.Duration(newTime.Second()) * time.Second
		newTime = newTime.Add(-second)
	}

	return newTime
}
