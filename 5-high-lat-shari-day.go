package prayer

import "time"

// ShariNormalDay is adapter following method that proposed by Mohamed Nabeel Tarabishy, Ph.D.
//
// He proposes that a normal day is defined as day when the fasting period is between 10h17m
// and 17h36m. If the day is "abnormal" then the fasting times is calculated using the schedule
// for area with 45 degrees latitude.
//
// This adapter doesn't require the sunrise and sunset to be exist in a day, so it's usable
// for area in extreme latitudes (>=65 degrees).
//
// Do note in this method there will be sudden changes in the length of the day of fasting.
// To avoid this issue, the author has given suggestion to just use the schedule from 45°
// on permanent basis. So, following that suggestion, you might be better using other adapter
// like `NearestLatitude` or `NearestLatitudeAsIs`.
//
// Reference: https://www.astronomycenter.net/pdf/tarabishyshigh_2014.pdf
func ShariNormalDay() HighLatitudeAdapter {
	return highLatShariNormalDay
}

func highLatShariNormalDay(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
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
		AsrConvention:      cfg.AsrConvention}
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
