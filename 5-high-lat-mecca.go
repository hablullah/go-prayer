package prayer

import "time"

func calcHighLatMecca(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	// Clear and split prayer times
	var fajrTimes, sunriseTimes, asrTimes, maghribTimes, ishaTimes []time.Time
	for _, s := range schedules {
		if !isScheduleNormal(s) {
			s.Fajr = time.Time{}
			s.Sunrise = time.Time{}
			s.Asr = time.Time{}
			s.Maghrib = time.Time{}
			s.Isha = time.Time{}
		}

		fajrTimes = append(fajrTimes, s.Fajr)
		sunriseTimes = append(sunriseTimes, s.Sunrise)
		asrTimes = append(asrTimes, s.Asr)
		maghribTimes = append(maghribTimes, s.Maghrib)
		ishaTimes = append(ishaTimes, s.Isha)
	}

	// Group empty time indexes
	emptyFajrIndexGroups := groupEmptyTimes(fajrTimes)
	emptySunriseIndexGroups := groupEmptyTimes(sunriseTimes)
	emptyAsrIndexGroups := groupEmptyTimes(asrTimes)
	emptyMaghribIndexGroups := groupEmptyTimes(maghribTimes)
	emptyIshaIndexGroups := groupEmptyTimes(ishaTimes)

	// Calculate schedule for Mecca
	meccaTz, _ := time.LoadLocation("Asia/Riyadh")
	meccaCfg := Config{
		Latitude:           21.425506007708996,
		Longitude:          39.8254579358597,
		Timezone:           meccaTz,
		TwilightConvention: cfg.TwilightConvention,
		AsrConvention:      cfg.AsrConvention,
		HighLatConvention:  Disabled}
	meccaSchedules, _ := calcNormal(meccaCfg, year)

	// Apply Mecca schedules to current location, by matching it with duration
	// in Mecca using transit time (noon) as the base.
	for i, s := range schedules {
		// Calculate duration from Mecca schedule
		ms := meccaSchedules[i]
		msFajrTransit := ms.Zuhr.Sub(ms.Fajr)
		msRiseTransit := ms.Zuhr.Sub(ms.Sunrise)
		msTransitAsr := ms.Asr.Sub(ms.Zuhr)
		msTransitMaghrib := ms.Maghrib.Sub(ms.Zuhr)
		msTransitIsha := ms.Isha.Sub(ms.Zuhr)

		// Apply Mecca times
		if fajrTimes[i].IsZero() {
			fajrTimes[i] = s.Zuhr.Add(-msFajrTransit)
		}

		if sunriseTimes[i].IsZero() {
			sunriseTimes[i] = s.Zuhr.Add(-msRiseTransit)
		}

		if asrTimes[i].IsZero() {
			asrTimes[i] = s.Zuhr.Add(msTransitAsr)
		}

		if maghribTimes[i].IsZero() {
			maghribTimes[i] = s.Zuhr.Add(msTransitMaghrib)
		}

		if ishaTimes[i].IsZero() {
			ishaTimes[i] = s.Zuhr.Add(msTransitIsha)
		}
	}

	// Create transition times
	maxDiff := 3 * time.Minute
	sunriseTimes = genTransitionTimes(year, sunriseTimes, emptySunriseIndexGroups, maxDiff)
	maghribTimes = genTransitionTimes(year, maghribTimes, emptyMaghribIndexGroups, maxDiff)
	asrTimes = genTransitionTimes(year, asrTimes, emptyAsrIndexGroups, maxDiff)
	fajrTimes = genTransitionTimes(year, fajrTimes, emptyFajrIndexGroups, maxDiff)
	ishaTimes = genTransitionTimes(year, ishaTimes, emptyIshaIndexGroups, maxDiff)

	// Apply the corrected times
	for i, s := range schedules {
		s.Fajr = fajrTimes[i]
		s.Sunrise = sunriseTimes[i]
		s.Asr = asrTimes[i]
		s.Maghrib = maghribTimes[i]
		s.Isha = ishaTimes[i]
		schedules[i] = s
	}

	return schedules
}

func genTransitionTimes(year int, times []time.Time, emptyIndexGroups [][]int, step time.Duration) []time.Time {
	for _, emptyIndexes := range emptyIndexGroups {
		// Get starter and end of empty indexes
		firstIdx := emptyIndexes[0]
		lastIdx := emptyIndexes[len(emptyIndexes)-1]

		// Create transition for before the empty group
		for i := firstIdx; i > firstIdx-30; i-- {
			// Get the reference, then move it to prev day
			refIdx := getRealIndex(i, times)
			refDay := times[refIdx]
			refDay = refDay.AddDate(0, 0, -1)
			refDay = time.Date(
				year, refDay.Month(), refDay.Day(),
				refDay.Hour(), refDay.Minute(), refDay.Second(),
				refDay.Nanosecond(), refDay.Location())

			// Get the prev day and adjust its limit
			prevIdx := getRealIndex(refIdx-1, times)
			prevDay := times[prevIdx]
			if prevDay.Sub(refDay).Abs() > step {
				times[prevIdx] = limitTimeDiff(prevDay, refDay, step, false)
			} else {
				break
			}
		}

		// Create transition for after the empty group
		for i := lastIdx; i < firstIdx+30; i++ {
			// Get the reference, then move it to next day
			refIdx := getRealIndex(i, times)
			refDay := times[refIdx]
			refDay = refDay.AddDate(0, 0, 1)
			refDay = time.Date(
				year, refDay.Month(), refDay.Day(),
				refDay.Hour(), refDay.Minute(), refDay.Second(),
				refDay.Nanosecond(), refDay.Location())

			// Get the next day and adjust its limit
			nextIdx := getRealIndex(refIdx+1, times)
			nextDay := times[nextIdx]
			if nextDay.Sub(refDay).Abs() > step {
				times[nextIdx] = limitTimeDiff(nextDay, refDay, step, true)
			} else {
				break
			}
		}
	}

	return times
}
