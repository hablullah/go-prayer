package prayer_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/hablullah/go-prayer"
)

// In this test, we'll look for missing prayer times in area with normal latitude
// (below and up to 45 degrees). In this area, every prayer times must be found since
// the Sun rise and set properly.
func Test_checkMissingScheduleInNormalArea(t *testing.T) {
	for longitude := -180.; longitude <= 180; longitude += 5 {
		timezones := genTzList(longitude)

		for latitude := -45.; latitude <= 45; latitude += 5 {
			for _, tz := range timezones {
				t.Logf("check abnormal in normal: lat=%f, long=%f, tz=%v",
					latitude, longitude, tz)

				cfg := prayer.Config{
					Latitude:  latitude,
					Longitude: longitude,
					Timezone:  tz,
				}

				schedules, err := prayer.Calculate(cfg, 2022)
				if err != nil {
					t.Fatal(err)
				}

				for _, s := range schedules {
					if s.Fajr.IsZero() || s.Sunrise.IsZero() ||
						s.Zuhr.IsZero() || s.Asr.IsZero() ||
						s.Maghrib.IsZero() || s.Isha.IsZero() {
						t.Fatalf("abnormal days found in lat=%f, long=%f, tz: %v, date=%s: %s",
							latitude, longitude, tz, s.Date, strSchedule(s))
					}
				}
			}
		}
	}
}

// In this test, we'll look for missing prayer times in area with higher latitude
// (between 45 and 64 degrees). In this area, Sunrise, Asr and Maghrib should be exist.
func Test_checkMissingScheduleInHigherArea(t *testing.T) {
	higherLatitudes := append(
		genLatitudeList(-45, -64, 5),
		genLatitudeList(45, 64, 5)...,
	)

	for longitude := -180.; longitude <= 180; longitude += 5 {
		timezones := genTzList(longitude)

		for _, latitude := range higherLatitudes {
			for _, tz := range timezones {
				t.Logf("check abnormal in higher: lat=%f, long=%f, tz=%v",
					latitude, longitude, tz)

				cfg := prayer.Config{
					Latitude:  latitude,
					Longitude: longitude,
					Timezone:  tz,
				}

				schedules, err := prayer.Calculate(cfg, 2022)
				if err != nil {
					t.Fatal(err)
				}

				for _, s := range schedules {
					if s.Sunrise.IsZero() || s.Zuhr.IsZero() || s.Asr.IsZero() || s.Maghrib.IsZero() {
						t.Fatalf("abnormal days found in lat=%f, long=%f, tz: %v, date=%s: %s",
							latitude, longitude, tz, s.Date, strSchedule(s))
					}
				}
			}
		}
	}
}

func genTzList(longitude float64) []*time.Location {
	h := int(math.Round(longitude / 15))
	if h <= -12 || h >= 12 {
		h = 12
	}

	// Generate timezones
	tz := time.FixedZone(fmt.Sprintf("UTC%+02d", h), h*60*60)
	tzDST := time.FixedZone(fmt.Sprintf("UTC%+02d", h+1), (h+1)*60*60)

	return []*time.Location{tz, tzDST}
}

func genLatitudeList(start, end, step float64) []float64 {
	if end < start {
		end, start = start, end
	}

	newStart := math.Floor(start/step) * step
	newEnd := math.Ceil(end/step) * step

	var list []float64
	for i := newStart; i <= newEnd; i += step {
		if i < start {
			list = append(list, start)
		} else if i > end {
			list = append(list, end)
		} else {
			list = append(list, i)
		}
	}

	return list
}

func strSchedule(s prayer.Schedule) string {
	bt, _ := json.Marshal(&s)
	return string(bt)
}
