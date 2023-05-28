package prayer

import (
	"math"
	"time"
)

// OneSeventhNight is adapter where the night period is divided into seven parts. Isha
// starts when the first seventh part ends, and Fajr starts when the last seventh part starts.
//
// This adapter depends on sunrise and sunset time, so it might not be suitable for area in
// extreme latitudes (>=65 degrees).
//
// Reference: http://praytimes.org/calculation
func OneSeventhNight() HighLatitudeAdapter {
	return highLatOneSeventhNight
}

func highLatOneSeventhNight(_ Config, _ int, schedules []PrayerSchedule) []PrayerSchedule {
	for i, s := range schedules {
		// Seventh night require Sunrise and Maghrib, and only done if Fajr or Isha missing
		if !s.Sunrise.IsZero() && !s.Maghrib.IsZero() && (s.Fajr.IsZero() || s.Isha.IsZero()) {
			// Calculate night duration
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := float64(24*60*60) - dayDuration

			// Calculate Fajr and Isha time
			seventhDuration := time.Duration(math.Round(nightDuration/7)) * time.Second
			schedules[i].Fajr = s.Sunrise.Add(-seventhDuration)
			schedules[i].Isha = s.Maghrib.Add(seventhDuration)
		}
	}

	return schedules
}
