package prayer

func calcHighLatNearestDay(schedules []PrayerSchedule) []PrayerSchedule {
	// Fetch the first last normal day
	var lastNormalIdx int
	var lastNormalSchedule PrayerSchedule

	if isScheduleNormal(schedules[0]) {
		for i := 1; i < len(schedules); i++ {
			if !isScheduleNormal(schedules[i]) {
				lastNormalIdx = i - 1
				lastNormalSchedule = schedules[i-1]
				break
			}
		}
	} else {
		for i := len(schedules) - 1; i >= 0; i-- {
			if isScheduleNormal(schedules[i]) {
				lastNormalIdx = i
				lastNormalSchedule = schedules[i]
				break
			}
		}
	}

	// Fix the schedule starting from the last normal idx
	i := lastNormalIdx + 1
	for {
		if i >= len(schedules) {
			i = 0
		}

		if i == lastNormalIdx {
			break
		}

		if s := schedules[i]; isScheduleNormal(s) {
			lastNormalSchedule = s
		} else {
			schedules[i] = lastNormalSchedule
		}

		i++
	}

	return schedules
}
