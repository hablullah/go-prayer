package prayer

func calcHighLatNearestDay(schedules []PrayerSchedule) []PrayerSchedule {
	abnormalSummer, abnormalWinter := extractAbnormalSchedules(schedules)

	for _, as := range []AbnormalRange{abnormalSummer, abnormalWinter} {
		// If this abnormal period is empty, skip
		if as.IsEmpty() {
			continue
		}

		// Get the last normal schedule
		abnormalIdxStart := as.Indexes[0]
		lastNormalSchedule := sliceAt(schedules, abnormalIdxStart-1)

		// Use the last normal schedule for the entire abnormal period
		for _, idx := range as.Indexes {
			schedules[idx] = lastNormalSchedule
		}
	}

	return schedules
}
