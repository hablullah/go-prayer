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

	// Make sure the new Sunrise and Sunset is not more than 5 minutes apart
	sunriseTimes = interpolateTimes(year, sunriseTimes, emptySunriseIndexGroups, 5*time.Minute)
	maghribTimes = interpolateTimes(year, maghribTimes, emptyMaghribIndexGroups, -5*time.Minute)

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

	fajrTimes = interpolateTimes(year, fajrTimes, emptyFajrIndexGroups, 5*time.Minute)
	ishaTimes = interpolateTimes(year, ishaTimes, emptyIshaIndexGroups, -5*time.Minute)

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

func groupEmptyTimes(times []time.Time) [][]int {
	// Group the empty times
	lastIdx := -2
	var currentGroup []int
	var indexesGroups [][]int
	for i, t := range times {
		if t.IsZero() {
			if i > lastIdx+1 && len(currentGroup) > 0 {
				indexesGroups = append(indexesGroups, currentGroup)
				currentGroup = []int{i}
			} else {
				currentGroup = append(currentGroup, i)
			}
			lastIdx = i
		}
	}

	// If:
	// 1) the final group not empty,
	// 2) the final group end with the last index of `times`
	// 3) the first group start with 0
	// then merge the final group with the first one
	if len(currentGroup) > 0 {
		firstOfFirstGroup := -1
		if len(indexesGroups) > 0 && len(indexesGroups[0]) > 0 {
			firstOfFirstGroup = indexesGroups[0][0]
		}

		lastOfFinalGroup := currentGroup[len(currentGroup)-1]
		if firstOfFirstGroup == 0 && lastOfFinalGroup == len(times)-1 {
			indexesGroups[0] = append(currentGroup, indexesGroups[0]...)
			currentGroup = []int{} // empty the current group
		}
	}

	// If current group is not empty, save
	if len(currentGroup) > 0 {
		indexesGroups = append(indexesGroups, currentGroup)
	}

	return indexesGroups
}

func interpolateTimes(year int, times []time.Time, emptyIndexGroups [][]int, step time.Duration) []time.Time {
	for _, emptyIndexes := range emptyIndexGroups {
		// Get min, max, mid time
		half := len(emptyIndexes) / 2

		// Split indexes into two
		firstHalf := emptyIndexes[:half+1]
		lastHalf := emptyIndexes[half+1:]

		// Fix the first half
		for _, idx := range firstHalf {
			// Get time from previous day, then set it date to today
			today := times[idx]
			prevDay := times[getRealIndex(idx-1, times)]
			prevDay = prevDay.AddDate(0, 0, 1)

			// Make sure the year for previous day is correct
			prevDay = time.Date(year, prevDay.Month(), prevDay.Day(),
				prevDay.Hour(), prevDay.Minute(), prevDay.Second(),
				prevDay.Nanosecond(), prevDay.Location())

			// Adjust the time
			adjusted := adjustTime(today, prevDay, step)
			times[idx] = adjusted
		}

		// Fix the last half
		for i := len(lastHalf) - 1; i >= 0; i-- {
			idx := lastHalf[i]
			today := times[idx]

			// Get time from the next day, then set it date to today
			nextDay := times[getRealIndex(idx+1, times)]
			nextDay = nextDay.AddDate(0, 0, -1)

			// Make sure the year for next day is correct
			nextDay = time.Date(year, nextDay.Month(), nextDay.Day(),
				nextDay.Hour(), nextDay.Minute(), nextDay.Second(),
				nextDay.Nanosecond(), nextDay.Location())

			// Adjust the time
			adjusted := adjustTime(today, nextDay, step)
			times[idx] = adjusted
		}
	}

	return times
}

func adjustTime(curent, reference time.Time, step time.Duration) time.Time {
	diff := curent.Sub(reference).Abs()
	if diff > step.Abs() {
		return reference.Add(step)
	} else {
		return curent
	}
}

func getRealIndex(idx int, list []time.Time) int {
	if idx < 0 {
		return len(list) - 1
	} else if idx > len(list)-1 {
		return 0
	} else {
		return idx
	}
}
