package prayer

import (
	"math"
	"testing"
	"time"
)

func TestCalculator_Calculate(t *testing.T) {
	// Prepare expected value.
	// The value is taken from "Mekanika Benda Langit", page 93.
	wib := time.FixedZone("WIB", 7*3600)
	expected := Times{
		Fajr:    time.Date(2009, 6, 12, 4, 35, 51, 0, wib),
		Sunrise: time.Date(2009, 6, 12, 5, 58, 18, 0, wib),
		Zuhr:    time.Date(2009, 6, 12, 11, 54, 26, 0, wib),
		Asr:     time.Date(2009, 6, 12, 15, 14, 25, 0, wib),
		Maghrib: time.Date(2009, 6, 12, 17, 46, 33, 0, wib),
		Isha:    time.Date(2009, 6, 12, 19, 00, 18, 0, wib),
	}

	// Calculate the times
	cfg := Config{
		Latitude:             -6.166667,
		Longitude:            106.85,
		Elevation:            50,
		CalculationMethod:    Default,
		AsrCalculationMethod: Shafii,
		PreciseToSeconds:     true,
	}

	date := time.Date(2009, 6, 12, 0, 0, 0, 0, wib)
	results := Calculate(date, cfg)

	// Compare between result and expectation, tolerated 1 second differences
	compare := func(name string, want, result time.Time) {
		if diff := want.Sub(result).Seconds(); math.Abs(diff) > 1 {
			t.Errorf("\n"+
				"name     : %s\n"+
				"expected : %s\n"+
				"get      : %s",
				name, want.Format("15:04:05"), result.Format("15:04:05"))
		}
	}

	compare("Fajr", expected.Fajr, results.Fajr)
	compare("Sunrise", expected.Sunrise, results.Sunrise)
	compare("Zuhr", expected.Zuhr, results.Zuhr)
	compare("Asr", expected.Asr, results.Asr)
	compare("Maghrib", expected.Maghrib, results.Maghrib)
	compare("Isha", expected.Isha, results.Isha)
}
