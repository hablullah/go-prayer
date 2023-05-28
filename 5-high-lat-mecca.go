package prayer

import (
	"math"
	"time"
)

// Mecca is adapter based on Fatwa from Dar Al Iftah Al Misrriyah number 2806 dated at 2010-08-08.
// They propose that area with higher latitude to follows the schedule in Mecca when abnormal days
// occured, using transit time as the common point. Here the day is considered "abnormal" when there
// are no true night, or the day length is less than 4 hours.
//
// To prevent sudden schedule changes, this method uses transition period for maximum one month
// before and after the abnormal periods.
//
// This adapter doesn't require the sunrise and sunset to be exist in a day, so it's usable
// for area in extreme latitudes (>=65 degrees).
//
// Reference: https://www.prayertimes.dk/fatawa.html
func Mecca() HighLatitudeAdapter {
	return highLatMecca
}

func highLatMecca(cfg Config, year int, schedules []PrayerSchedule) []PrayerSchedule {
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
		AsrConvention:      cfg.AsrConvention}
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
		nTransitionDays := leftoverDays / 2

		// Apply transition times
		fajrTimes = createMeccaPreTransition(fajrTimes, abnormalPeriod, nTransitionDays)
		fajrTimes = createMeccaPostTransition(fajrTimes, abnormalPeriod, nTransitionDays)
		sunriseTimes = createMeccaPreTransition(sunriseTimes, abnormalPeriod, nTransitionDays)
		sunriseTimes = createMeccaPostTransition(sunriseTimes, abnormalPeriod, nTransitionDays)
		asrTimes = createMeccaPreTransition(asrTimes, abnormalPeriod, nTransitionDays)
		asrTimes = createMeccaPostTransition(asrTimes, abnormalPeriod, nTransitionDays)
		maghribTimes = createMeccaPreTransition(maghribTimes, abnormalPeriod, nTransitionDays)
		maghribTimes = createMeccaPostTransition(maghribTimes, abnormalPeriod, nTransitionDays)
		ishaTimes = createMeccaPreTransition(ishaTimes, abnormalPeriod, nTransitionDays)
		ishaTimes = createMeccaPostTransition(ishaTimes, abnormalPeriod, nTransitionDays)
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
		winterSummerTransitionDays := winterSummerGap / 2 // for post-winter and pre-summer
		summerWinterTransitionDays := summerWinterGap / 2 // for pre-winter and post-summer

		preWinterTransitionDays := summerWinterTransitionDays
		postWinterTransitionDays := winterSummerTransitionDays
		preSummerTransitionDays := winterSummerTransitionDays
		postSummerTransitionDays := summerWinterTransitionDays

		// Create winter transition
		fajrTimes = createMeccaPreTransition(fajrTimes, abnormalWinter, preWinterTransitionDays)
		fajrTimes = createMeccaPostTransition(fajrTimes, abnormalWinter, postWinterTransitionDays)
		sunriseTimes = createMeccaPreTransition(sunriseTimes, abnormalWinter, preWinterTransitionDays)
		sunriseTimes = createMeccaPostTransition(sunriseTimes, abnormalWinter, postWinterTransitionDays)
		asrTimes = createMeccaPreTransition(asrTimes, abnormalWinter, preWinterTransitionDays)
		asrTimes = createMeccaPostTransition(asrTimes, abnormalWinter, postWinterTransitionDays)
		maghribTimes = createMeccaPreTransition(maghribTimes, abnormalWinter, preWinterTransitionDays)
		maghribTimes = createMeccaPostTransition(maghribTimes, abnormalWinter, postWinterTransitionDays)
		ishaTimes = createMeccaPreTransition(ishaTimes, abnormalWinter, preWinterTransitionDays)
		ishaTimes = createMeccaPostTransition(ishaTimes, abnormalWinter, postWinterTransitionDays)

		// Create summer transition
		fajrTimes = createMeccaPreTransition(fajrTimes, abnormalSummer, preSummerTransitionDays)
		fajrTimes = createMeccaPostTransition(fajrTimes, abnormalSummer, postSummerTransitionDays)
		sunriseTimes = createMeccaPreTransition(sunriseTimes, abnormalSummer, preSummerTransitionDays)
		sunriseTimes = createMeccaPostTransition(sunriseTimes, abnormalSummer, postSummerTransitionDays)
		asrTimes = createMeccaPreTransition(asrTimes, abnormalSummer, preSummerTransitionDays)
		asrTimes = createMeccaPostTransition(asrTimes, abnormalSummer, postSummerTransitionDays)
		maghribTimes = createMeccaPreTransition(maghribTimes, abnormalSummer, preSummerTransitionDays)
		maghribTimes = createMeccaPostTransition(maghribTimes, abnormalSummer, postSummerTransitionDays)
		ishaTimes = createMeccaPreTransition(ishaTimes, abnormalSummer, preSummerTransitionDays)
		ishaTimes = createMeccaPostTransition(ishaTimes, abnormalSummer, postSummerTransitionDays)
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

func createMeccaPreTransition(times []time.Time, abnormalPeriod AbnormalRange, nTransitionDays int) []time.Time {
	// Fix transition days
	if nTransitionDays > 30 {
		nTransitionDays = 30
	}

	// Get data where transition end, i.e. when abnormality start
	endIdx, _ := firstSliceItem(abnormalPeriod.Indexes)
	endTime := times[endIdx]

	// Get data where transition begin
	startIdx := sliceRealIdx(times, endIdx-nTransitionDays)
	startTime := times[startIdx]
	if startTime.After(endTime) {
		startTime = startTime.AddDate(-1, 0, 0)
	}

	// Calculate duration step
	durationDiff := endTime.Sub(startTime)
	diffStep := durationDiff / time.Duration(nTransitionDays)

	// Apply transition time
	ci, ct := endIdx, endTime
	for i := nTransitionDays - 1; i > 0; i-- { // minus 1 to exclude `startIdx`
		ci = sliceRealIdx(times, ci-1)
		ct = ct.Add(-diffStep)
		times[ci] = ct
	}

	return times
}

func createMeccaPostTransition(times []time.Time, abnormalPeriod AbnormalRange, nTransitionDays int) []time.Time {
	// Fix transition days
	if nTransitionDays > 30 {
		nTransitionDays = 30
	}

	// Get data where transition start, i.e. when abnormality end
	startIdx, _ := lastSliceItem(abnormalPeriod.Indexes)
	startTime := times[startIdx]

	// Get data where transition end
	endIdx := sliceRealIdx(times, startIdx+nTransitionDays)
	endTime := times[endIdx]
	if endTime.Before(startTime) {
		endTime = endTime.AddDate(1, 0, 0)
	}

	// Calculate duration step
	durationDiff := endTime.Sub(startTime)
	diffStep := durationDiff / time.Duration(nTransitionDays)

	// Apply transition time
	ci, ct := startIdx, startTime
	for i := 1; i < nTransitionDays; i++ {
		ci = sliceRealIdx(times, ci+1)
		ct = ct.Add(diffStep)
		times[ci] = ct
	}

	return times
}
