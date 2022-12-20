package prayer

import (
	"math"
	"time"
)

func calcHighLatAngleBased(cfg Config, schedules []PrayerSchedule) []PrayerSchedule {
	// Fetch the twilight angle
	var fajrAngle, ishaAngle float64
	if cfg.TwilightConvention != nil {
		fajrAngle = cfg.TwilightConvention.FajrAngle
		ishaAngle = cfg.TwilightConvention.IshaAngle
	}

	// If twilight angle missing, use the astronomical twilight
	if fajrAngle == 0 {
		fajrAngle = AstronomicalTwilight.FajrAngle
	}

	if ishaAngle == 0 {
		ishaAngle = AstronomicalTwilight.IshaAngle
	}

	// Apply schedules
	for i, s := range schedules {
		// Angle based require Sunrise and Maghrib
		if s.Sunrise.IsZero() || s.Maghrib.IsZero() {
			continue
		}

		// Calculate night duration
		dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
		nightDuration := float64(24*60*60) - dayDuration

		// Calculate Fajr time
		fajrDuration := time.Duration(math.Round(nightDuration/60*fajrAngle)) * time.Second
		schedules[i].Fajr = s.Sunrise.Add(-fajrDuration)

		// Calculate Isha time
		ishaDuration := time.Duration(math.Round(nightDuration/60*ishaAngle)) * time.Second
		schedules[i].Isha = s.Maghrib.Add(ishaDuration)
	}

	return schedules
}
