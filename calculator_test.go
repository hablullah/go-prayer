package prayer

import (
	"math"
	"testing"
	"time"
)

func TestCalculator_Calculate(t *testing.T) {
	calc := Calculator{
		Latitude:             -6.1751,
		Longitude:            106.8650,
		Elevation:            7.9,
		CalculationMethod:    Default,
		AsrCalculationMethod: Shafii,
		UseIhtiyat:           true,
	}

	location := time.FixedZone("Jakarta", 7*3600)
	date := time.Date(2018, 1, 7, 0, 0, 0, 0, location)
	adhan, _ := calc.Calculate(date)

	expectedResult := Times{
		Fajr:    time.Date(2018, 1, 7, 4, 22, 0, 0, location),
		Sunrise: time.Date(2018, 1, 7, 5, 42, 0, 0, location),
		Zuhr:    time.Date(2018, 1, 7, 12, 2, 0, 0, location),
		Asr:     time.Date(2018, 1, 7, 15, 27, 0, 0, location),
		Maghrib: time.Date(2018, 1, 7, 18, 15, 0, 0, location),
		Isha:    time.Date(2018, 1, 7, 19, 30, 0, 0, location),
	}

	for prayer, prayerTime := range adhan {
		expected := expectedResult[prayer]

		// 1 minute difference is still tolerated
		if math.Abs(prayerTime.Sub(expected).Minutes()) > 1 {
			t.Errorf("Prayer %d expected %s get %s", prayer,
				expected.Format("15:04:05"),
				prayerTime.Format("15:04:05"))
		}
	}
}
