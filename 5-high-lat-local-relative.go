package prayer

import (
	"time"
)

func calcLocalRelativeEstimation(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	var (
		nFajrSample     int
		nIshaSample     int
		sumFajrPercents float64
		sumIshaPercents float64
	)

	for _, s := range schedules {
		// This conventions only works if daytime exists (in other words, sunrise
		// and Maghrib must exist). So if there are days where those time don't
		// exist, stop and just return the schedule as it is.
		// TODO: maybe put some warning log later.
		if s.Sunrise.IsZero() || s.Maghrib.IsZero() {
			return schedules
		}

		// Calculate percentage in normal days
		if s.IsNormal {
			// Calculate day and night
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := 24*60*60 - dayDuration

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
	avgFajrPercents := sumFajrPercents / float64(nFajrSample)
	avgIshaPercents := sumIshaPercents / float64(nIshaSample)

	// Extract abnormal schedules
	abnormalSummer, abnormalWinter := extractAbnormalSchedules(schedules)

	// Fix Fajr and Isha times in abnormal days
	for _, as := range []AbnormalRange{abnormalSummer, abnormalWinter} {
		for _, i := range as.Indexes {
			s := schedules[i]
			dayDuration := s.Maghrib.Sub(s.Sunrise).Seconds()
			nightDuration := 24*60*60 - dayDuration

			if !s.IsNormal {
				fajrDuration := time.Duration(nightDuration*avgFajrPercents) * time.Second
				schedules[i].Fajr = s.Sunrise.Add(-fajrDuration)

				ishaDuration := time.Duration(nightDuration*avgIshaPercents) * time.Second
				schedules[i].Isha = s.Maghrib.Add(ishaDuration)
			}
		}
	}

	schedules = applyLocalRelativeTransition(schedules, abnormalSummer)
	schedules = applyLocalRelativeTransition(schedules, abnormalWinter)
	return schedules
}

func applyLocalRelativeTransition(schedules []PrayerSchedule, abnormalPeriod AbnormalRange) []PrayerSchedule {
	// If there are no abnormality, return as it is
	if abnormalPeriod.IsEmpty() {
		return schedules
	}

	// Split the abnormal period into two
	nAbnormalDays := len(abnormalPeriod.Indexes)
	maxTransitionDays := nAbnormalDays / 2
	firstHalf := abnormalPeriod.Indexes[:maxTransitionDays]
	secondHalf := abnormalPeriod.Indexes[nAbnormalDays-maxTransitionDays:]

	// Fix the time in first half
	for _, idx := range firstHalf {
		today := sliceAt(schedules, idx)
		yesterday := sliceAt(schedules, idx-1)

		// If idx is zero, it means today is the first day of the year.
		// Therefore, yesterday is occured last year.
		if idx == 0 {
			yesterday.Fajr = yesterday.Fajr.AddDate(-1, 0, 0)
			yesterday.Isha = yesterday.Isha.AddDate(-1, 0, 0)
		}

		var fajrChanged, ishaChanged bool
		schedules[idx].Fajr, fajrChanged = applyLocalRelativeTransitionTime(yesterday.Fajr, today.Fajr)
		schedules[idx].Isha, ishaChanged = applyLocalRelativeTransitionTime(yesterday.Isha, today.Isha)
		if !fajrChanged && !ishaChanged {
			break
		}
	}

	// Fix the time in second half, do it backward
	for i := len(secondHalf) - 1; i >= 0; i-- {
		idx := secondHalf[i]
		today := sliceAt(schedules, idx)
		tomorrow := sliceAt(schedules, idx+1)

		// If idx is last, it means today is the last day of the year.
		// Therefore, tomorrow will occur next year.
		if idx == len(schedules)-1 {
			tomorrow.Fajr = tomorrow.Fajr.AddDate(1, 0, 0)
			tomorrow.Isha = tomorrow.Isha.AddDate(1, 0, 0)
		}

		var fajrChanged, ishaChanged bool
		schedules[idx].Fajr, fajrChanged = applyLocalRelativeTransitionTime(tomorrow.Fajr, today.Fajr)
		schedules[idx].Isha, ishaChanged = applyLocalRelativeTransitionTime(tomorrow.Isha, today.Isha)
		if !fajrChanged && !ishaChanged {
			break
		}
	}

	return schedules
}

func applyLocalRelativeTransitionTime(reference, today time.Time) (time.Time, bool) {
	// Calculate diff between today and reference
	var diff time.Duration
	var referenceIsForward bool
	if today.After(reference) {
		diff = today.Sub(reference)
		referenceIsForward = false
	} else {
		diff = reference.Sub(today)
		referenceIsForward = true
	}

	// Limit the difference
	maxDiff := 24*time.Hour + 5*time.Minute
	minDiff := 24*time.Hour - 5*time.Minute

	if diff > maxDiff {
		diff = maxDiff
	} else if diff < minDiff {
		diff = minDiff
	} else {
		return today, false // the diff is within limit, nothing to change
	}

	// Adjust the time
	var newTime time.Time
	if referenceIsForward {
		newTime = reference.Add(-diff)
	} else {
		newTime = reference.Add(diff)
	}

	// Fix the year
	if todayYear := today.Year(); newTime.Year() != todayYear {
		newTime = time.Date(todayYear, newTime.Month(), newTime.Day(),
			newTime.Hour(), newTime.Minute(), newTime.Second(), newTime.Nanosecond(),
			newTime.Location())
	}

	return newTime, true
}
