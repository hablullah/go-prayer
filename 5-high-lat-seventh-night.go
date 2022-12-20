package prayer

import (
	"math"
	"time"
)

func calcHighLatOneSeventhNight(schedules []PrayerSchedule) []PrayerSchedule {
	for i, s := range schedules {
		// Sevent Night require Sunrise and Maghrib
		if s.Sunrise.IsZero() || s.Maghrib.IsZero() {
			continue
		}

		// Calculate night duration
		dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
		nightDuration := float64(24*60*60) - dayDuration

		// Calculate Fajr and Isha time
		seventhDuration := time.Duration(math.Round(nightDuration/7)) * time.Second
		schedules[i].Fajr = s.Sunrise.Add(-seventhDuration)
		schedules[i].Isha = s.Maghrib.Add(seventhDuration)
	}

	return schedules
}
