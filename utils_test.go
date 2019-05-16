package prayer

import (
	"math"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func Test_getJulianDay(t *testing.T) {
	jkt := time.FixedZone("WIB", 7*60*60)
	scenarios := map[time.Time]float64{
		time.Date(-4712, 1, 1, 12, 0, 0, 0, time.UTC):  0,
		time.Date(-4712, 1, 2, 0, 0, 0, 0, time.UTC):   0.5,
		time.Date(-4712, 1, 2, 12, 0, 0, 0, time.UTC):  1,
		time.Date(1582, 10, 4, 0, 0, 0, 0, time.UTC):   2299159.5,
		time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC):  2299160.5,
		time.Date(1945, 8, 17, 0, 0, 0, 0, time.UTC):   2431684.5,
		time.Date(1974, 9, 27, 0, 0, 0, 0, time.UTC):   2442317.5,
		time.Date(624, 2, 26, 0, 0, 0, 0, time.UTC):    1949029.5,
		time.Date(-2961, 1, 1, 19, 47, 4, 0, time.UTC): 639553.32435,
		time.Date(2009, 6, 12, 12, 0, 0, 0, jkt):       2454994.7083,
	}

	for date, expected := range scenarios {
		jd := getJulianDay(date)
		diff := decimal.NewFromFloat(jd).
			Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"date     : %s\n"+
				"expected : %f\n"+
				"get      : %f",
				date.Format("2006-01-02 15:04:05 -07"),
				expected, jd)
		}
	}
}

func Test_getEquationOfTime(t *testing.T) {
	scenarios := map[float64]float64{
		2455292:     -2.70815,
		2454994.708: 0.175,
	}

	for jd, expected := range scenarios {
		eoT := getEquationOfTime(jd)

		eoT = round(eoT, 3)
		expected = round(expected, 3)

		if eoT-expected != 0 {
			t.Errorf("%f expected %f get %f", jd, expected, eoT)
		}
	}
}

func Test_sin(t *testing.T) {
	scenarios := map[float64]float64{
		0:  0,
		30: 0.5,
		45: 1.0 / math.Sqrt(2),
		60: math.Sqrt(3) / 2,
		90: 1,
	}

	for degree, expected := range scenarios {
		result := sin(degree)

		result = round(result, 3)
		expected = round(expected, 3)

		if result-expected != 0 {
			t.Errorf("%f degree expected %f get %f", degree, expected, result)
		}
	}
}

func Test_cos(t *testing.T) {
	scenarios := map[float64]float64{
		0:  1,
		30: math.Sqrt(3) / 2,
		45: 1.0 / math.Sqrt(2),
		60: 0.5,
		90: 0,
	}

	for degree, expected := range scenarios {
		result := cos(degree)

		result = round(result, 3)
		expected = round(expected, 3)

		if result-expected != 0 {
			t.Errorf("%f degree expected %f get %f", degree, expected, result)
		}
	}
}

func Test_tan(t *testing.T) {
	scenarios := map[float64]float64{
		0:  0,
		30: 1.0 / math.Sqrt(3),
		45: 1,
		60: math.Sqrt(3),
	}

	for degree, expected := range scenarios {
		result := tan(degree)

		result = round(result, 3)
		expected = round(expected, 3)

		if result-expected != 0 {
			t.Errorf("%f degree expected %f get %f", degree, expected, result)
		}
	}
}

func Test_acos(t *testing.T) {
	scenarios := map[float64]float64{
		1:                  0,
		math.Sqrt(3) / 2:   30,
		1.0 / math.Sqrt(2): 45,
		0.5:                60,
		0:                  90,
	}

	for src, expected := range scenarios {
		degree := acos(src)

		degree = round(degree, 3)
		expected = round(expected, 3)

		if degree-expected != 0 {
			t.Errorf("%f expected %f get %f", src, expected, degree)
		}
	}
}

func Test_round(t *testing.T) {
	scenarios := map[float64]float64{
		1.0 / 3.0: 0.333,
		1.0 / 4.0: 0.25,
		1.0 / 6.0: 0.167,
		1.0 / 7.0: 0.143,
		1.0 / 8.0: 0.125,
		1.0 / 9.0: 0.111,
	}

	for src, expected := range scenarios {
		result := round(src, 3)

		if result-expected != 0 {
			t.Errorf("Expected %f get %f", expected, result)
		}
	}
}
