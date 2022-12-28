package prayer

import "time"

func calcLocalRelativeEstimation(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	// Split prayer times and calculate percentage
	var nDaySample, nFajrSample, nIshaSample int
	var sumDayPercents, sumFajrPercents, sumIshaPercents float64
	var fajrTimes, sunriseTimes, maghribTimes, ishaTimes []time.Time

	for _, s := range schedules {
		// Split prayer times
		fajrTimes = append(fajrTimes, s.Fajr)
		sunriseTimes = append(sunriseTimes, s.Sunrise)
		maghribTimes = append(maghribTimes, s.Maghrib)
		ishaTimes = append(ishaTimes, s.Isha)

		// Calculate percentage in normal days
		if !s.Sunrise.IsZero() && !s.Maghrib.IsZero() {
			// Calculate day percentage
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := float64(24*60*60) - dayDuration
			sumDayPercents += dayDuration / (24 * 60 * 60)
			nDaySample++

			// Calculate Fajr percentage
			if !s.Fajr.IsZero() {
				fajrDuration := s.Sunrise.Sub(s.Fajr).Seconds()
				sumFajrPercents += fajrDuration / nightDuration
				nFajrSample++
			}

			// Calculate Isha percentage
			if !s.Isha.IsZero() {
				ishaDuration := s.Isha.Sub(s.Maghrib).Seconds()
				sumIshaPercents += ishaDuration / nightDuration
				nIshaSample++
			}
		}
	}

	// Calculate average percentage
	avgDayPercents := sumDayPercents / float64(nDaySample)
	avgFajrPercents := sumFajrPercents / float64(nFajrSample)
	avgIshaPercents := sumIshaPercents / float64(nIshaSample)
	avgDayDuration := time.Second * time.Duration(avgDayPercents*24*60*60)

	// Group empty time indexes
	emptyFajrIndexGroups := groupEmptyTimes(fajrTimes)
	emptySunriseIndexGroups := groupEmptyTimes(sunriseTimes)
	emptyMaghribIndexGroups := groupEmptyTimes(maghribTimes)
	emptyIshaIndexGroups := groupEmptyTimes(ishaTimes)

	// Fix Sunrise and Maghrib time
	for i := range schedules {
		sunrise, maghrib := sunriseTimes[i], maghribTimes[i]

		switch {
		case !sunrise.IsZero() && maghrib.IsZero():
			maghribTimes[i] = sunrise.Add(avgDayDuration)
		case sunrise.IsZero() && !maghrib.IsZero():
			sunriseTimes[i] = maghrib.Add(-avgDayDuration)
		case sunrise.IsZero() && maghrib.IsZero():
			noon := schedules[i].Zuhr
			sunriseTimes[i] = noon.Add(-(avgDayDuration / 2))
			maghribTimes[i] = sunrise.Add(avgDayDuration)
		}
	}

	maxDiff := 5 * time.Minute
	sunriseTimes = interpolateEmptyTimes(year, sunriseTimes, emptySunriseIndexGroups, maxDiff)
	maghribTimes = interpolateEmptyTimes(year, maghribTimes, emptyMaghribIndexGroups, maxDiff)

	// Fix Fajr and Isha times
	for i := range schedules {
		sunrise, maghrib := sunriseTimes[i], maghribTimes[i]
		dayDuration := maghrib.Sub(sunrise).Seconds()
		nightDuration := float64(24*60*60) - dayDuration

		if fajr := fajrTimes[i]; fajr.IsZero() {
			fajrDuration := time.Duration(nightDuration*avgFajrPercents) * time.Second
			fajrTimes[i] = sunrise.Add(-fajrDuration)
		}

		if isha := ishaTimes[i]; isha.IsZero() {
			ishaDuration := time.Duration(nightDuration*avgIshaPercents) * time.Second
			ishaTimes[i] = maghrib.Add(ishaDuration)
		}
	}

	fajrTimes = interpolateEmptyTimes(year, fajrTimes, emptyFajrIndexGroups, maxDiff)
	ishaTimes = interpolateEmptyTimes(year, ishaTimes, emptyIshaIndexGroups, maxDiff)

	// Apply the corrected times
	for i, s := range schedules {
		s.Fajr = fajrTimes[i]
		s.Sunrise = sunriseTimes[i]
		s.Maghrib = maghribTimes[i]
		s.Isha = ishaTimes[i]
		schedules[i] = s
	}

	return schedules
}

func interpolateEmptyTimes(year int, times []time.Time, emptyIndexGroups [][]int, step time.Duration) []time.Time {
	for _, emptyIndexes := range emptyIndexGroups {
		// Split indexes into two
		half := len(emptyIndexes) / 2
		firstHalf := emptyIndexes[:half+1]
		lastHalf := emptyIndexes[half+1:]

		// Fix the first half
		for _, idx := range firstHalf {
			// Get time from previous day, then set it date to transit
			prevDay := times[getRealIndex(idx-1, times)]
			prevDay = prevDay.AddDate(0, 0, 1)
			prevDay = time.Date(
				year, prevDay.Month(), prevDay.Day(),
				prevDay.Hour(), prevDay.Minute(), prevDay.Second(),
				prevDay.Nanosecond(), prevDay.Location())

			// Adjust the time
			today := times[idx]
			times[idx] = limitTimeDiff(today, prevDay, step, true)
		}

		// Fix the last half
		for i := len(lastHalf) - 1; i >= 0; i-- {
			idx := lastHalf[i]

			// Get time from the next day, then set it date to transit
			nextDay := times[getRealIndex(idx+1, times)]
			nextDay = nextDay.AddDate(0, 0, -1)
			nextDay = time.Date(
				year, nextDay.Month(), nextDay.Day(),
				nextDay.Hour(), nextDay.Minute(), nextDay.Second(),
				nextDay.Nanosecond(), nextDay.Location())

			// Adjust the time
			today := times[idx]
			times[idx] = limitTimeDiff(today, nextDay, step, false)
		}
	}

	return times
}
