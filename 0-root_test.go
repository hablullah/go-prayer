package prayer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hablullah/go-prayer"
	"github.com/hablullah/go-prayer/internal/datatest"
)

func TestCalculate(t *testing.T) {
	testCalculate(t, datatest.Tromso)     // North Frigid
	testCalculate(t, datatest.London)     // North Temperate
	testCalculate(t, datatest.Jakarta)    // Torrid
	testCalculate(t, datatest.Wellington) // South Temperate
}

func testCalculate(t *testing.T, td datatest.TestData) {
	// Calculate schedules
	schedules, err := prayer.Calculate(prayer.Config{
		Latitude:            td.Latitude,
		Longitude:           td.Longitude,
		Timezone:            td.Timezone,
		TwilightConvention:  prayer.AstronomicalTwilight(),
		AsrConvention:       prayer.Shafii,
		HighLatitudeAdapter: prayer.NearestLatitude(),
		PreciseToSeconds:    true,
	}, 2023)

	msg := fmt.Sprintf("schedule in %s has error: %v", td.Name, err)
	assertNil(t, err, msg)

	nExpected, nResult := len(td.Schedules), len(schedules)
	msg = fmt.Sprintf("%s schedule size: want %d got %d", td.Name, nExpected, nResult)
	assertEqual(t, nExpected, nResult, msg)

	for i := range schedules {
		result := schedules[i]
		expected := td.Schedules[i]
		assertSchedule(t, td, expected, result)
	}
}

func assertSchedule(t *testing.T, td datatest.TestData, e, r prayer.Schedule) {
	// Calculate diff
	diffFajr := e.Fajr.Sub(r.Fajr).Abs()
	diffSunrise := e.Sunrise.Sub(r.Sunrise).Abs()
	diffZuhr := e.Zuhr.Sub(r.Zuhr).Abs()
	diffAsr := e.Asr.Sub(r.Asr).Abs()
	diffMaghrib := e.Maghrib.Sub(r.Maghrib).Abs()
	diffIsha := e.Isha.Sub(r.Isha).Abs()

	// Prepare log message
	msgFormat := "%s, %s => want %q got %q (%v)"
	fajrMsg := fmt.Sprintf(msgFormat, td.Name, "Fajr", e.Fajr, r.Fajr, diffFajr)
	sunriseMsg := fmt.Sprintf(msgFormat, td.Name, "Sunrise", e.Sunrise, r.Sunrise, diffSunrise)
	zuhrMsg := fmt.Sprintf(msgFormat, td.Name, "Zuhr", e.Zuhr, r.Zuhr, diffZuhr)
	asrMsg := fmt.Sprintf(msgFormat, td.Name, "Asr", e.Asr, r.Asr, diffAsr)
	maghribMsg := fmt.Sprintf(msgFormat, td.Name, "Maghrib", e.Maghrib, r.Maghrib, diffMaghrib)
	ishaMsg := fmt.Sprintf(msgFormat, td.Name, "Isha", e.Isha, r.Isha, diffIsha)

	// Diff only allowed up to 5 seconds
	maxDiff := 5 * time.Second
	assertLTE(t, diffFajr, maxDiff, fajrMsg)
	assertLTE(t, diffSunrise, maxDiff, sunriseMsg)
	assertLTE(t, diffZuhr, maxDiff, zuhrMsg)
	assertLTE(t, diffAsr, maxDiff, asrMsg)
	assertLTE(t, diffMaghrib, maxDiff, maghribMsg)
	assertLTE(t, diffIsha, maxDiff, ishaMsg)
}

func assertNil(t *testing.T, v any, msg string) {
	if v != nil {
		t.Error(msg)
	}
}

func assertEqual[T comparable](t *testing.T, a T, b T, msg string) {
	if a != b {
		t.Error(msg)
	}
}

func assertLTE[T int | float64 | time.Duration](t *testing.T, a T, b T, msg string) {
	if a > b {
		t.Error(msg)
	}
}
