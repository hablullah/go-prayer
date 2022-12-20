package prayer

func calcHighLatNearestLatitude(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
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
		if s.Fajr.IsZero() || s.Sunrise.IsZero() ||
			s.Maghrib.IsZero() || s.Isha.IsZero() {
			schedules[i] = newSchedules[i]
		}
	}

	return schedules
}
