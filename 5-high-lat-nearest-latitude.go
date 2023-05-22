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
	nearestSchedules, _ := calcNormal(newCfg, year)

	// Apply schedules from nearest latitude in abnormal period, using transit as common point.
	for i, s := range schedules {
		if s.IsNormal {
			continue
		}

		// Calculate duration from schedule for nearest latitude
		ns := nearestSchedules[i]
		nsFajrRise := ns.Sunrise.Sub(ns.Fajr)
		nsRiseTransit := ns.Zuhr.Sub(ns.Sunrise)
		nsTransitMaghrib := ns.Maghrib.Sub(ns.Zuhr)
		nsMaghribIsha := ns.Isha.Sub(ns.Maghrib)

		// Apply the duration
		if s.Sunrise.IsZero() {
			s.Sunrise = s.Zuhr.Add(-nsRiseTransit)
		}

		if s.Maghrib.IsZero() {
			s.Maghrib = s.Zuhr.Add(nsTransitMaghrib)
		}

		s.Fajr = s.Sunrise.Add(-nsFajrRise)
		s.Isha = s.Maghrib.Add(nsMaghribIsha)
		schedules[i] = s
	}

	return schedules
}
