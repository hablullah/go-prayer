package prayer

import (
	"time"
)

// MiddleNight is adapter where the night period is divided into two halves. The
// first half is considered to be the "night" and the other half as "day break".
// Fajr and Isha in this method are assumed to be at mid-night during the abnormal
// periods.
//
// This adapter depends on sunrise and sunset time, so it might not be suitable for
// area in extreme latitudes (>=65 degrees).
//
// Reference: http://praytimes.org/calculation
func MiddleNight() HighLatitudeAdapter {
	return highLatMiddleNight
}

func highLatMiddleNight(_ Config, _ int, schedules []Schedule) []Schedule {
	for i, s := range schedules {
		// Middle night require Sunrise and Maghrib, and only done if Fajr or Isha missing
		if !s.Sunrise.IsZero() && !s.Maghrib.IsZero() && (s.Fajr.IsZero() || s.Isha.IsZero()) {
			// Calculate night duration
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := float64(24*60*60) - dayDuration

			// Calculate Fajr and Isha time
			halfDuration := time.Duration(nightDuration * 0.5 * float64(time.Second))
			schedules[i].Fajr = s.Sunrise.Add(-halfDuration)
			schedules[i].Isha = s.Maghrib.Add(halfDuration)
		}
	}

	return schedules
}
