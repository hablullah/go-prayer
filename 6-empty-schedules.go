package prayer

import (
	"time"
)

type EmptyRange struct {
	Indexes      []int
	LeftIndexes  []int
	RightIndexes []int
}

func fixEmptyTimes(schedules []PrayerSchedule) []PrayerSchedule {
	// Split schedules for each time
	nSchedules := len(schedules)
	transits := make([]time.Time, nSchedules)
	fajrTimes := make([]time.Time, nSchedules)
	sunriseTimes := make([]time.Time, nSchedules)
	asrTimes := make([]time.Time, nSchedules)
	maghribTimes := make([]time.Time, nSchedules)
	ishaTimes := make([]time.Time, nSchedules)

	var fajrHasEmpty bool
	var sunriseHasEmpty bool
	var asrHasEmpty bool
	var maghribHasEmpty bool
	var ishaHasEmpty bool

	for idx, s := range schedules {
		transits[idx] = s.Zuhr
		fajrTimes[idx] = s.Fajr
		sunriseTimes[idx] = s.Sunrise
		asrTimes[idx] = s.Asr
		maghribTimes[idx] = s.Maghrib
		ishaTimes[idx] = s.Isha

		fajrHasEmpty = fajrHasEmpty || s.Fajr.IsZero()
		sunriseHasEmpty = sunriseHasEmpty || s.Sunrise.IsZero()
		asrHasEmpty = asrHasEmpty || s.Asr.IsZero()
		maghribHasEmpty = maghribHasEmpty || s.Maghrib.IsZero()
		ishaHasEmpty = ishaHasEmpty || s.Isha.IsZero()
	}

	// Fill the empty times
	if fajrHasEmpty {
		emptyRanges := extractEmptyTimes(schedules, fajrTimes)
		fajrTimes = interpolateEmptyTimes(transits, fajrTimes, emptyRanges)
	}

	if sunriseHasEmpty {
		emptyRanges := extractEmptyTimes(schedules, sunriseTimes)
		sunriseTimes = interpolateEmptyTimes(transits, sunriseTimes, emptyRanges)
	}

	if asrHasEmpty {
		emptyRanges := extractEmptyTimes(schedules, asrTimes)
		asrTimes = interpolateEmptyTimes(transits, asrTimes, emptyRanges)
	}

	if maghribHasEmpty {
		emptyRanges := extractEmptyTimes(schedules, maghribTimes)
		maghribTimes = interpolateEmptyTimes(transits, maghribTimes, emptyRanges)
	}

	if ishaHasEmpty {
		emptyRanges := extractEmptyTimes(schedules, ishaTimes)
		ishaTimes = interpolateEmptyTimes(transits, ishaTimes, emptyRanges)
	}

	// Put back times to schedule
	for idx, s := range schedules {
		s.Fajr = fajrTimes[idx]
		s.Sunrise = sunriseTimes[idx]
		s.Asr = asrTimes[idx]
		s.Maghrib = maghribTimes[idx]
		s.Isha = ishaTimes[idx]
		schedules[idx] = s
	}

	return schedules
}

func extractEmptyTimes(schedules []PrayerSchedule, times []time.Time) []EmptyRange {
	// If there are no schedules, return empty
	if len(schedules) == 0 || len(times) == 0 || len(schedules) != len(times) {
		return nil
	}

	// Loop each time
	lastEmptyIdx := -2
	var ranges []EmptyRange
	var currentRange EmptyRange
	for i, s := range schedules {
		t := times[i]

		// If schedule is abnormal or time not empty, skip
		if !s.IsNormal || !t.IsZero() {
			continue
		}

		// If current schedule is continuation of the previous, just add the index
		if i == lastEmptyIdx+1 {
			currentRange.Indexes = append(currentRange.Indexes, i)
		} else {
			// If current range is not empty, save it to range list
			if len(currentRange.Indexes) > 0 {
				ranges = append(ranges, currentRange)
			}

			// Re-initiate current range
			currentRange = EmptyRange{Indexes: []int{i}}
		}

		lastEmptyIdx = i
	}

	// Handle leftover range
	if len(currentRange.Indexes) > 0 {
		// Check if we can merge the leftover range to the first range:
		// - there are more than one range
		// - the first index of first range = zero
		// - the last index of current range = last schedule index
		// If we can't merge it, just append current range to the end
		if len(ranges) > 0 &&
			ranges[0].Indexes[0] == 0 &&
			currentRange.Indexes[len(currentRange.Indexes)-1] == len(schedules)-1 {
			ranges[0].Indexes = append(currentRange.Indexes, ranges[0].Indexes...)
		} else {
			ranges = append(ranges, currentRange)
		}
	}

	// At this point we might have several empty groups, so now we need to fetch its boundaries
	for i, r := range ranges {
		// Get indexes of empty ranges
		firstIdx, _ := firstSliceItem(r.Indexes)
		lastIdx, _ := lastSliceItem(r.Indexes)

		// Catch the left boundary
		var leftTime time.Time
		leftIdx := sliceRealIdx(schedules, firstIdx-1)
		if schedules[leftIdx].IsNormal && !times[leftIdx].IsZero() {
			leftTime = times[leftIdx]
		}

		// Catch the right boundary
		var rightTime time.Time
		rightIdx := sliceRealIdx(schedules, lastIdx+1)
		if schedules[rightIdx].IsNormal && !times[rightIdx].IsZero() {
			rightTime = times[rightIdx]
		}

		// If both left and right exist, use both
		if !leftTime.IsZero() && !rightTime.IsZero() {
			r.LeftIndexes = []int{leftIdx}
			r.RightIndexes = []int{rightIdx}
			ranges[i] = r
			continue
		}

		// If we only have left boundaries, try to catch one more
		if !leftTime.IsZero() && rightTime.IsZero() {
			lefterIdx := sliceRealIdx(schedules, leftIdx-1)
			if schedules[lefterIdx].IsNormal && !times[lefterIdx].IsZero() {
				r.LeftIndexes = []int{lefterIdx, leftIdx}
			} else {
				r.LeftIndexes = []int{leftIdx}
			}
			ranges[i] = r
			continue
		}

		// If we only have right boundaries, also try to catch one more
		if leftTime.IsZero() && !rightTime.IsZero() {
			righterIdx := sliceRealIdx(schedules, rightIdx+1)
			if schedules[righterIdx].IsNormal && !times[righterIdx].IsZero() {
				r.RightIndexes = []int{rightIdx, righterIdx}
			} else {
				r.RightIndexes = []int{rightIdx}
			}
			ranges[i] = r
			continue
		}
	}

	return ranges
}

func interpolateEmptyTimes(transits, times []time.Time, emptyRanges []EmptyRange) []time.Time {
	for _, r := range emptyRanges {
		// Scenario 1: empty range has both left and right boundary
		if len(r.LeftIndexes) == 1 && len(r.RightIndexes) == 1 {
			nDays := len(r.Indexes)
			leftIdx, rightIdx := r.LeftIndexes[0], r.RightIndexes[0]
			leftTime, rightTime := times[leftIdx], times[rightIdx]
			leftTransit, rightTransit := transits[leftIdx], transits[rightIdx]

			leftDuration := leftTransit.Sub(leftTime)
			rightDuration := rightTransit.Sub(rightTime)
			durationDiff := rightDuration - leftDuration
			diffStep := durationDiff / time.Duration(nDays+1)

			currentDuration := leftDuration
			for _, idx := range r.Indexes {
				currentDuration += diffStep
				times[idx] = transits[idx].Add(-currentDuration)
			}

			continue
		}

		// Scenario 2: empty range only has complete left boundary
		if len(r.LeftIndexes) == 2 && len(r.RightIndexes) == 0 {
			lefterIdx, leftIdx := r.LeftIndexes[0], r.LeftIndexes[1]
			lefterTime, leftTime := times[lefterIdx], times[leftIdx]
			lefterTransit, leftTransit := transits[lefterIdx], transits[leftIdx]

			lefterDuration := lefterTransit.Sub(lefterTime)
			leftDuration := leftTransit.Sub(leftTime)
			diffStep := leftDuration - lefterDuration

			currentDuration := leftDuration
			for _, idx := range r.Indexes {
				currentDuration += diffStep
				times[idx] = transits[idx].Add(-currentDuration)
			}

			continue
		}

		// Scenario 3: empty range only has complete right boundary
		if len(r.LeftIndexes) == 0 && len(r.RightIndexes) == 2 {
			rightIdx, righterIdx := r.RightIndexes[0], r.RightIndexes[1]
			rightTime, righterTime := times[rightIdx], times[righterIdx]
			rightTransit, righterTransit := transits[rightIdx], transits[righterIdx]

			rightDuration := rightTransit.Sub(rightTime)
			righterDuration := righterTransit.Sub(righterTime)
			diffStep := righterDuration - rightDuration

			currentDuration := rightDuration
			for j := len(r.Indexes) - 1; j >= 0; j-- {
				idx := r.Indexes[j]
				currentDuration -= diffStep // minus because backward
				times[idx] = transits[idx].Add(-currentDuration)
			}

			continue
		}

		// Scenario 4: empty range only has partial boundary
		if len(r.LeftIndexes) == 1 || len(r.RightIndexes) == 1 {
			var boundaryIdx int
			if len(r.LeftIndexes) == 1 {
				boundaryIdx = r.LeftIndexes[0]
			} else {
				boundaryIdx = r.RightIndexes[0]
			}

			boundaryTime := times[boundaryIdx]
			boundaryTransit := transits[boundaryIdx]
			currentDuration := boundaryTransit.Sub(boundaryTime)
			for _, idx := range r.Indexes {
				times[idx] = transits[idx].Add(-currentDuration)
			}

			continue
		}
	}

	return times
}
