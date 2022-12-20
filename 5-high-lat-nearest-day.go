package prayer

func calcHighLatNearestDay(schedules []PrayerSchedule) []PrayerSchedule {
	// Helper function
	isNormal := func(s PrayerSchedule) bool {
		return !s.Fajr.IsZero() && !s.Sunrise.IsZero() && !s.Asr.IsZero() &&
			!s.Maghrib.IsZero() && !s.Isha.IsZero()
	}

	// Fetch the first last normal day
	var lastNormalIdx int
	var lastNormalSchedule PrayerSchedule

	if isNormal(schedules[0]) {
		lastNormalIdx = 0
		lastNormalSchedule = schedules[0]
	} else {
		for i := len(schedules) - 1; i >= 0; i-- {
			if isNormal(schedules[i]) {
				lastNormalIdx = i
				lastNormalSchedule = schedules[i]
				break
			}
		}
	}

	// Fix the schedule starting from the last normal idx
	for i := lastNormalIdx + 1; i != lastNormalIdx; i++ {
		if i >= len(schedules) {
			i = 0
		}

		if s := schedules[i]; isNormal(s) {
			lastNormalSchedule = s
		} else {
			schedules[i] = lastNormalSchedule
		}
	}

	return schedules
}
