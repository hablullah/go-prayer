package prayer

import "time"

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
	nearestSchedules, _ := calcNormal(newCfg, year)

	// Apply schedules for the abnormal days using schedules from nearest latitude
	// with transit as common point.
	minFastingDuration := 10*time.Hour + 17*time.Minute
	maxFastingDuration := 17*time.Hour + 36*time.Minute
	for i, s := range schedules {
		// If day is normal, just continue
		fastingDuration := s.Maghrib.Sub(s.Fajr)
		if s.IsNormal && fastingDuration >= minFastingDuration && fastingDuration <= maxFastingDuration {
			continue
		}

		// Calculate duration from schedule for nearest latitude
		ns := nearestSchedules[i]
		nsFajrTransit := ns.Zuhr.Sub(ns.Fajr)
		nsRiseTransit := ns.Zuhr.Sub(ns.Sunrise)
		nsTransitAsr := ns.Asr.Sub(ns.Zuhr)
		nsTransitMaghrib := ns.Maghrib.Sub(ns.Zuhr)
		nsTransitIsha := ns.Isha.Sub(ns.Zuhr)

		// Apply durations
		s.Fajr = s.Zuhr.Add(-nsFajrTransit)
		s.Sunrise = s.Zuhr.Add(-nsRiseTransit)
		s.Asr = s.Zuhr.Add(nsTransitAsr)
		s.Maghrib = s.Zuhr.Add(nsTransitMaghrib)
		s.Isha = s.Zuhr.Add(nsTransitIsha)
		schedules[i] = s
	}

	return schedules
}
