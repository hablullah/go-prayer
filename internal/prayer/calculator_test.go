package prayer

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestTimeCalculator_getHourAngle(t *testing.T) {
	calc := TimeCalculator{
		Latitude:       decimal.NewFromFloat(-6.166667),
		sunDeclination: decimal.NewFromFloat(23.16099835),
	}

	scenarios := []struct {
		name     string
		sunAlt   float64
		expected float64
	}{{
		name:     "Jakarta, Fajr",
		sunAlt:   -20,
		expected: 109.1441394,
	}, {
		name:     "Jakarta, Sunrise",
		sunAlt:   -1.07869939,
		expected: 88.53151863,
	}, {
		name:     "Jakarta, Asr",
		sunAlt:   32.63075274,
		expected: 50.496359,
	}, {
		name:     "Jakarta, Isha",
		sunAlt:   -18,
		expected: 106.9681811,
	}}

	for _, s := range scenarios {
		sunAltitude := decimal.NewFromFloat(s.sunAlt)
		hourAngle := calc.getHourAngle(sunAltitude)
		diff := hourAngle.Sub(decimal.NewFromFloat(s.expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"name     : %s\n"+
				"expected : %f\n"+
				"get      : %s",
				s.name, s.expected, hourAngle.String())
		}
	}
}

func TestTimeCalculator_getEquationOfTime(t *testing.T) {
	calc := TimeCalculator{}

	scenarios := map[float64]float64{
		2455292:     -2.70815,
		2454994.708: 0.175,
	}

	for flJD, expected := range scenarios {
		jd := decimal.NewFromFloat(flJD)
		eoT := calc.getEquationOfTime(jd)
		diff := eoT.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"JD       : %s\n"+
				"expected : %f\n"+
				"get      : %s",
				jd.String(), expected, eoT.String())
		}
	}
}
