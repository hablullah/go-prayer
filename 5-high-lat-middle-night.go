package prayer

import (
	"math"
	"time"
)

func calcHighLatMiddleNight(schedules []PrayerSchedule) []PrayerSchedule {
	for i, s := range schedules {
		// Middle night require Sunrise and Maghrib, and only done if Fajr or Isha missing
		if !s.Sunrise.IsZero() && !s.Maghrib.IsZero() && (s.Fajr.IsZero() || s.Isha.IsZero()) {
			// Calculate night duration
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := float64(24*60*60) - dayDuration

			// Calculate Fajr and Isha time
			halfDuration := time.Duration(math.Floor(nightDuration/2)) * time.Second
			schedules[i].Fajr = s.Sunrise.Add(-halfDuration)
			schedules[i].Isha = s.Maghrib.Add(halfDuration)
		}
	}

	return schedules
}
