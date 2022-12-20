package prayer

import "time"

var (
	fastingMin = 10*time.Hour + 17*time.Minute
	fastingMax = 17*time.Hour + 36*time.Minute
)

func calcHighLatShariNormalDay(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	// Get the nearest latitude
	latitude := cfg.Latitude
	if latitude > 45 {
		latitude = 45
	} else if latitude < -45 {
		latitude = -45
	}

	// Calculate schedule for the nearest latitude
	newCfg := Config{
		Latitude:           latitude,
		Longitude:          cfg.Longitude,
		Timezone:           cfg.Timezone,
		TwilightConvention: cfg.TwilightConvention,
		AsrConvention:      cfg.AsrConvention,
		HighLatConvention:  Disabled}
	newSchedules, _ := calcNormal(newCfg, year)

	// Apply schedules for the abnormal days
	for i, s := range schedules {
		fastingDuration := s.Maghrib.Sub(s.Fajr)
		if s.Fajr.IsZero() || s.Sunrise.IsZero() ||
			s.Maghrib.IsZero() || s.Isha.IsZero() ||
			fastingDuration < fastingMin ||
			fastingDuration > fastingMax {
			schedules[i] = newSchedules[i]
		}
	}

	return schedules
}
