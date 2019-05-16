package prayer

import (
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
		diff := jd.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"date     : %s\n"+
				"expected : %f\n"+
				"get      : %s",
				date.Format("2006-01-02 15:04:05 -07"),
				expected, jd.String())
		}
	}
}

func Test_getEquationOfTime(t *testing.T) {
	scenarios := map[float64]float64{
		2455292:     -2.70815,
		2454994.708: 0.175,
	}

	for jd, expected := range scenarios {
		decJD := decimal.NewFromFloat(jd)
		eoT := getEquationOfTime(decJD)
		diff := eoT.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"JD       : %s\n"+
				"expected : %f\n"+
				"get      : %s",
				decJD.String(), expected, eoT.String())
		}
	}
}

func Test_getTimezone(t *testing.T) {
	scenarios := map[time.Time]int64{
		time.Date(2009, 01, 02, 0, 0, 0, 0, time.UTC):                               0,
		time.Date(2009, 01, 02, 0, 0, 0, 0, time.FixedZone("WIB-A", 7*60*60)):       7,
		time.Date(2009, 01, 02, 0, 0, 0, 0, time.FixedZone("WIB-B", 7*60*60+1800)):  7,
		time.Date(2009, 01, 02, 0, 0, 0, 0, time.FixedZone("EST-A", -5*60*60)):      -5,
		time.Date(2009, 01, 02, 0, 0, 0, 0, time.FixedZone("EST-B", -5*60*60-1800)): -5,
	}

	for date, expected := range scenarios {
		if timezone := getTimezone(date); timezone != expected {
			t.Errorf("\n"+
				"location : %s\n"+
				"expected : %d\n"+
				"get      : %d",
				date.Location().String(),
				expected, timezone)
		}
	}
}

func Test_getTransitTime(t *testing.T) {
	scenarios := []struct {
		name      string
		date      time.Time
		longitude float64
		expected  float64
	}{{
		name:      "Jakarta, 2009-06-12",
		date:      time.Date(2009, 6, 12, 0, 0, 0, 0, time.FixedZone("WIB", 7*60*60)),
		longitude: 106.85,
		expected:  11.87375,
	}}

	for _, s := range scenarios {
		date := time.Date(
			s.date.Year(),
			s.date.Month(),
			s.date.Day(),
			12, 0, 0, 0,
			s.date.Location())

		jd := getJulianDay(date)
		eoT := getEquationOfTime(jd)
		timezone := getTimezone(date)
		transitTime := getTransitTime(timezone, s.longitude, eoT)
		diff := transitTime.Sub(decimal.NewFromFloat(s.expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"name     : %s\n"+
				"expected : %f\n"+
				"get      : %s",
				s.name, s.expected, transitTime.String())
		}
	}
}

func Test_getHourAngle(t *testing.T) {
	latitude := float64(-6.166667)
	sunDeclination := decimal.NewFromFloat(23.16099835)

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
		sunAlt := decimal.NewFromFloat(s.sunAlt)
		hourAngle := getHourAngle(latitude, sunAlt, sunDeclination)
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
