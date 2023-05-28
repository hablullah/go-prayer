package prayer

import (
	"time"
)

// AlwaysMecca is similar with `Mecca`, except it will be applied every day and not only on
// the "abnormal" days.
//
// This adapter doesn't require the sunrise and sunset to be exist in a day, so it's usable
// for area in extreme latitudes (>=65 degrees).
func AlwaysMecca() HighLatitudeAdapter {
	return highLatAlwaysMecca
}

func highLatAlwaysMecca(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	// Calculate schedule for Mecca
	meccaTz, _ := time.LoadLocation("Asia/Riyadh")
	meccaCfg := Config{
		Latitude:           21.425506007708996,
		Longitude:          39.8254579358597,
		Timezone:           meccaTz,
		TwilightConvention: cfg.TwilightConvention,
		AsrConvention:      cfg.AsrConvention}
	meccaSchedules, _ := calcNormal(meccaCfg, year)

	// Apply schedules to current location, by matching it with duration in Mecca
	// using transit time (noon) as the base.
	for i, s := range schedules {
		// Calculate duration from Mecca schedule
		ms := meccaSchedules[i]
		msFajrTransit := ms.Zuhr.Sub(ms.Fajr)
		msRiseTransit := ms.Zuhr.Sub(ms.Sunrise)
		msTransitAsr := ms.Asr.Sub(ms.Zuhr)
		msTransitMaghrib := ms.Maghrib.Sub(ms.Zuhr)
		msTransitIsha := ms.Isha.Sub(ms.Zuhr)

		// Apply Mecca times
		s.Fajr = s.Zuhr.Add(-msFajrTransit)
		s.Sunrise = s.Zuhr.Add(-msRiseTransit)
		s.Asr = s.Zuhr.Add(msTransitAsr)
		s.Maghrib = s.Zuhr.Add(msTransitMaghrib)
		s.Isha = s.Zuhr.Add(msTransitIsha)
		schedules[i] = s
	}

	return schedules
}
