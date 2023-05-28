package prayer

import (
	"math"
	"time"
)

// AngleBased is adapter where the night period is divided into several parts,
// depending on the value of twilight angle for Fajr and Isha.
//
// For example, let a be the twilight angle for Isha, and let t = a/60. The period
// between sunset and sunrise is divided into t parts. Isha begins after the first
// part. So, if the twilight angle for Isha is 15, then Isha begins at the end of the
// first quarter (15/60) of the night. Time for Fajr is calculated similarly.
//
// This adapter depends on sunrise and sunset time, so it might not be suitable for
// area in extreme latitudes (>=65 degrees).
//
// Reference: http://praytimes.org/calculation
func AngleBased() HighLatitudeAdapter {
	return highLatAngleBased
}

func highLatAngleBased(cfg Config, year int, schedules []Schedule) []Schedule {
	// Fetch the twilight angle
	var fajrAngle, ishaAngle float64
	if cfg.TwilightConvention != nil {
		fajrAngle = cfg.TwilightConvention.FajrAngle
		ishaAngle = cfg.TwilightConvention.IshaAngle
	}

	// If twilight angle missing, use the astronomical twilight
	astronomical := AstronomicalTwilight()
	if fajrAngle == 0 {
		fajrAngle = astronomical.FajrAngle
	}

	if ishaAngle == 0 {
		ishaAngle = astronomical.IshaAngle
	}

	// Apply schedules
	for i, s := range schedules {
		// Angle based require Sunrise and Maghrib, and only done if Fajr or Isha missing
		if !s.Sunrise.IsZero() && !s.Maghrib.IsZero() && (s.Fajr.IsZero() || s.Isha.IsZero()) {
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
	}

	return schedules
}
