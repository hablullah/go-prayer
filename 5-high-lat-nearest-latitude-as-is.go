package prayer

// NearestLatitudeAsIs is similar with `NearestLatitude` except it will use the
// schedule from 45 degrees latitude as it is without any change. Like `NearestLatitude`,
// this method will change the schedule for entire year to prevent sudden changes in
// fasting time.
//
// This adapter doesn't require the sunrise and sunset to be exist in a day, so it's
// usable for area in extreme latitudes (>=65 degrees).
//
// Reference: https://fiqh.islamonline.net/en/praying-and-fasting-at-high-latitudes/
func NearestLatitudeAsIs() HighLatitudeAdapter {
	return highLatNearestLatitudeAsIs
}

func highLatNearestLatitudeAsIs(cfg Config, year int, _ []PrayerSchedule) []PrayerSchedule {
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
	newSchedules, _ := calcNormal(newCfg, year)
	return newSchedules
}
