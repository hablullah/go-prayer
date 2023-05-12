package prayer

import (
	"math"
	"time"
)

func calcHighLatMecca(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
	// Additional rule: day is normal if daylength is more than 4 hour
	for i, s := range schedules {
		if s.IsNormal {
			var dayLength time.Duration
			if !s.Maghrib.IsZero() && !s.Sunrise.IsZero() {
				dayLength = s.Maghrib.Sub(s.Sunrise)
			}
			schedules[i].IsNormal = dayLength >= 4*time.Hour
		}
	}

	// Extract abnormal schedules
	abnormalSummer, abnormalWinter := extractAbnormalSchedules(schedules)

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

	// Apply Mecca schedules in abnormal period by matching it with duration
	// in Mecca using transit time (noon) as the base.
	for _, as := range []AbnormalRange{abnormalSummer, abnormalWinter} {
		for _, i := range as.Indexes {
			// Calculate duration from Mecca schedule
			ms := meccaSchedules[i]
			msFajrTransit := ms.Zuhr.Sub(ms.Fajr)
			msRiseTransit := ms.Zuhr.Sub(ms.Sunrise)
			msTransitAsr := ms.Asr.Sub(ms.Zuhr)
			msTransitMaghrib := ms.Maghrib.Sub(ms.Zuhr)
			msTransitIsha := ms.Isha.Sub(ms.Zuhr)

			// Apply Mecca duration
			s := schedules[i]
			s.Fajr = s.Zuhr.Add(-msFajrTransit)
			s.Sunrise = s.Zuhr.Add(-msRiseTransit)
			s.Asr = s.Zuhr.Add(msTransitAsr)
			s.Maghrib = s.Zuhr.Add(msTransitMaghrib)
			s.Isha = s.Zuhr.Add(msTransitIsha)
			schedules[i] = s
		}
	}

	schedules = applyMeccaTransition(schedules, abnormalSummer, abnormalWinter)
	return schedules
}

func applyMeccaTransition(schedules []PrayerSchedule, abnormalSummer, abnormalWinter AbnormalRange) []PrayerSchedule {
	// If there are no abnormality, return as it is
	if abnormalSummer.IsEmpty() && abnormalWinter.IsEmpty() {
		return schedules
	}

	// Split schedules for each time
	nSchedules := len(schedules)
	fajrTimes := make([]time.Time, nSchedules)
	sunriseTimes := make([]time.Time, nSchedules)
	asrTimes := make([]time.Time, nSchedules)
	maghribTimes := make([]time.Time, nSchedules)
	ishaTimes := make([]time.Time, nSchedules)

	for idx, s := range schedules {
		fajrTimes[idx] = s.Fajr
		sunriseTimes[idx] = s.Sunrise
		asrTimes[idx] = s.Asr
		maghribTimes[idx] = s.Maghrib
		ishaTimes[idx] = s.Isha
	}

	// Check if there is only one abnormal period
	onlySummer := abnormalWinter.IsEmpty() && !abnormalSummer.IsEmpty()
	onlyWinter := abnormalSummer.IsEmpty() && !abnormalWinter.IsEmpty()
	if onlySummer || onlyWinter {
		// Merge into one abnormal period
		abnormalPeriod := abnormalSummer
		if abnormalPeriod.IsEmpty() {
			abnormalPeriod = abnormalWinter
		}

		// Calculate transition duration from leftover days
		leftoverDays := len(schedules) - len(abnormalPeriod.Indexes)
		nTransitionDays := int(math.Floor(float64(leftoverDays) / 2))

		// Apply transition times
		fajrTimes = createMeccaPreTransition(schedules, fajrTimes, abnormalPeriod, nTransitionDays)
		fajrTimes = createMeccaPostTransition(schedules, fajrTimes, abnormalPeriod, nTransitionDays)
		sunriseTimes = createMeccaPreTransition(schedules, sunriseTimes, abnormalPeriod, nTransitionDays)
		sunriseTimes = createMeccaPostTransition(schedules, sunriseTimes, abnormalPeriod, nTransitionDays)
		asrTimes = createMeccaPreTransition(schedules, asrTimes, abnormalPeriod, nTransitionDays)
		asrTimes = createMeccaPostTransition(schedules, asrTimes, abnormalPeriod, nTransitionDays)
		maghribTimes = createMeccaPreTransition(schedules, maghribTimes, abnormalPeriod, nTransitionDays)
		maghribTimes = createMeccaPostTransition(schedules, maghribTimes, abnormalPeriod, nTransitionDays)
		ishaTimes = createMeccaPreTransition(schedules, ishaTimes, abnormalPeriod, nTransitionDays)
		ishaTimes = createMeccaPostTransition(schedules, ishaTimes, abnormalPeriod, nTransitionDays)
	} else if !abnormalSummer.IsEmpty() && !abnormalWinter.IsEmpty() {
		// Fetch indexes
		summerIdxStart, _ := firstSliceItem(abnormalSummer.Indexes)
		summerIdxEnd, _ := lastSliceItem(abnormalSummer.Indexes)
		winterIdxStart, _ := firstSliceItem(abnormalWinter.Indexes)
		winterIdxEnd, _ := lastSliceItem(abnormalWinter.Indexes)

		// Calculate gap between period
		winterSummerGap := int(math.Abs(float64(summerIdxStart-winterIdxEnd))) - 1 // after winter end, before summer start
		summerWinterGap := int(math.Abs(float64(winterIdxStart-summerIdxEnd))) - 1 // after summer end, before winter start

		// Calculate transition duration
		winterSummerTransitionDays := int(math.Floor(float64(winterSummerGap) / 2)) // for post-winter and pre-summer
		summerWinterTransitionDays := int(math.Floor(float64(summerWinterGap) / 2)) // for pre-winter and post-summer

		preWinterTransitionDays := summerWinterTransitionDays
		postWinterTransitionDays := winterSummerTransitionDays
		preSummerTransitionDays := winterSummerTransitionDays
		postSummerTransitionDays := summerWinterTransitionDays

		// Create winter transition
		fajrTimes = createMeccaPreTransition(schedules, fajrTimes, abnormalWinter, preWinterTransitionDays)
		fajrTimes = createMeccaPostTransition(schedules, fajrTimes, abnormalWinter, postWinterTransitionDays)
		sunriseTimes = createMeccaPreTransition(schedules, sunriseTimes, abnormalWinter, preWinterTransitionDays)
		sunriseTimes = createMeccaPostTransition(schedules, sunriseTimes, abnormalWinter, postWinterTransitionDays)
		asrTimes = createMeccaPreTransition(schedules, asrTimes, abnormalWinter, preWinterTransitionDays)
		asrTimes = createMeccaPostTransition(schedules, asrTimes, abnormalWinter, postWinterTransitionDays)
		maghribTimes = createMeccaPreTransition(schedules, maghribTimes, abnormalWinter, preWinterTransitionDays)
		maghribTimes = createMeccaPostTransition(schedules, maghribTimes, abnormalWinter, postWinterTransitionDays)
		ishaTimes = createMeccaPreTransition(schedules, ishaTimes, abnormalWinter, preWinterTransitionDays)
		ishaTimes = createMeccaPostTransition(schedules, ishaTimes, abnormalWinter, postWinterTransitionDays)

		// Create summer transition
		fajrTimes = createMeccaPreTransition(schedules, fajrTimes, abnormalSummer, preSummerTransitionDays)
		fajrTimes = createMeccaPostTransition(schedules, fajrTimes, abnormalSummer, postSummerTransitionDays)
		sunriseTimes = createMeccaPreTransition(schedules, sunriseTimes, abnormalSummer, preSummerTransitionDays)
		sunriseTimes = createMeccaPostTransition(schedules, sunriseTimes, abnormalSummer, postSummerTransitionDays)
		asrTimes = createMeccaPreTransition(schedules, asrTimes, abnormalSummer, preSummerTransitionDays)
		asrTimes = createMeccaPostTransition(schedules, asrTimes, abnormalSummer, postSummerTransitionDays)
		maghribTimes = createMeccaPreTransition(schedules, maghribTimes, abnormalSummer, preSummerTransitionDays)
		maghribTimes = createMeccaPostTransition(schedules, maghribTimes, abnormalSummer, postSummerTransitionDays)
		ishaTimes = createMeccaPreTransition(schedules, ishaTimes, abnormalSummer, preSummerTransitionDays)
		ishaTimes = createMeccaPostTransition(schedules, ishaTimes, abnormalSummer, postSummerTransitionDays)
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

func createMeccaPreTransition(schedules []PrayerSchedule, times []time.Time, abnormalPeriod AbnormalRange, nTransitionDays int) []time.Time {
	// Fix transition days
	if nTransitionDays > 30 {
		nTransitionDays = 30
	}

	// Get data where transition end, i.e. when abnormality start
	endIdx, _ := firstSliceItem(abnormalPeriod.Indexes)
	endTime := times[endIdx]
	endTransit := schedules[endIdx].Zuhr
	endDuration := endTransit.Sub(endTime)

	// Get data where transition begin
	startIdx := sliceRealIdx(schedules, endIdx-nTransitionDays)
	startTime := times[startIdx]
	startTransit := schedules[startIdx].Zuhr
	startDuration := startTransit.Sub(startTime)

	// Calculate duration step
	durationDiff := endDuration - startDuration
	diffStep := durationDiff / time.Duration(nTransitionDays)

	// Apply transition time
	ci, duration := endIdx, endDuration
	for i := nTransitionDays - 1; i > 0; i-- { // minus 1 to exclude `startIdx`
		ci = sliceRealIdx(times, ci-1)
		duration -= diffStep // minus because backward
		times[ci] = schedules[ci].Zuhr.Add(-duration)
	}

	return times
}

func createMeccaPostTransition(schedules []PrayerSchedule, times []time.Time, abnormalPeriod AbnormalRange, nTransitionDays int) []time.Time {
	// Fix transition days
	if nTransitionDays > 30 {
		nTransitionDays = 30
	}

	// Get data where transition start, i.e. when abnormality end
	startIdx, _ := lastSliceItem(abnormalPeriod.Indexes)
	startTime := times[startIdx]
	startTransit := schedules[startIdx].Zuhr
	startDuration := startTransit.Sub(startTime)

	// Get data where transition end
	endIdx := sliceRealIdx(schedules, startIdx+nTransitionDays)
	endTime := times[endIdx]
	endTransit := schedules[endIdx].Zuhr
	endDuration := endTransit.Sub(endTime)

	// Calculate duration step
	durationDiff := endDuration - startDuration
	diffStep := durationDiff / time.Duration(nTransitionDays)

	// Apply transition time
	ci, duration := startIdx, startDuration
	for i := 1; i < nTransitionDays; i++ {
		ci = sliceRealIdx(times, ci+1)
		duration += diffStep
		times[ci] = schedules[ci].Zuhr.Add(-duration)
	}

	return times
}
