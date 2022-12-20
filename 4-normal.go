package prayer

import (
	"math"
	"time"

	"github.com/hablullah/go-sampa"
)

func calcNormal(cfg Config, year int) ([]PrayerSchedule, int) {
	// Prepare location
	location := sampa.Location{
		Latitude:  cfg.Latitude,
		Longitude: cfg.Longitude,
		Elevation: cfg.Elevation,
	}

	// Prepare custom Sun events
	customEvents := []sampa.CustomSunEvent{{
		Name:          "fajr",
		BeforeTransit: true,
		Elevation: func(_ sampa.SunPosition) float64 {
			return -cfg.TwilightConvention.FajrAngle
		},
	}, {
		Name:          "isha",
		BeforeTransit: false,
		Elevation: func(_ sampa.SunPosition) float64 {
			return -cfg.TwilightConvention.IshaAngle
		},
	}, {
		Name:          "asr",
		BeforeTransit: false,
		Elevation: func(todayData sampa.SunPosition) float64 {
			a := getAsrCoefficient(cfg.AsrConvention)
			b := math.Abs(todayData.TopocentricDeclination - cfg.Latitude)
			elevation := acot(a + math.Tan(degToRad(b)))
			return radToDeg(elevation)
		},
	}}

	// Calculate schedules for each day in a year. Here we also calculate the first day
	// of next year and the last day of the previous year. This is useful to check if
	// some schedules chained to tomorrow or yesterday events.
	base := time.Date(year, 1, 1, 0, 0, 0, 0, cfg.Timezone)
	start := base.AddDate(0, 0, -1)
	limit := base.AddDate(1, 0, 0)
	nDays := int(limit.Sub(start).Hours()/24) + 1

	// Create slice to contain result
	schedules := make([]PrayerSchedule, nDays)

	// Calculate each day
	var idx int
	for dt := start; !dt.After(limit); dt = dt.AddDate(0, 0, 1) {
		// Calculate the schedules
		e, _ := sampa.GetSunEvents(dt, location, nil, customEvents...)
		transit := e.Transit.DateTime

		fajr := e.Others["fajr"].DateTime
		sunrise := e.Sunrise.DateTime
		asr := e.Others["asr"].DateTime
		maghrib := e.Sunset.DateTime
		isha := e.Others["isha"].DateTime

		// Adjust the index
		fajrIdx := adjustScheduleIdx(schedules, idx, fajr, transit, true)
		sunriseIdx := adjustScheduleIdx(schedules, idx, sunrise, transit, true)
		maghribIdx := adjustScheduleIdx(schedules, idx, maghrib, transit, false)
		ishaIdx := adjustScheduleIdx(schedules, idx, isha, transit, false)

		// Save the schedules
		schedules[idx].Date = dt.Format("2006-01-02")
		schedules[fajrIdx].Fajr = fajr
		schedules[sunriseIdx].Sunrise = sunrise
		schedules[idx].Zuhr = transit
		schedules[idx].Asr = asr
		schedules[maghribIdx].Maghrib = maghrib
		schedules[ishaIdx].Isha = isha

		idx++
	}

	// Adjust slice so we only see schedules for this year
	schedules = schedules[1 : len(schedules)-1]

	// Final check
	var nAbnormal int
	for i, s := range schedules {
		// Clean up schedule
		if !s.Fajr.Before(s.Zuhr) || (!s.Sunrise.IsZero() && !s.Fajr.Before(s.Sunrise)) {
			s.Fajr = time.Time{}
		}

		if !s.Sunrise.Before(s.Zuhr) {
			s.Sunrise = time.Time{}
		}

		if !s.Maghrib.After(s.Zuhr) {
			s.Maghrib = time.Time{}
		}

		if !s.Isha.After(s.Zuhr) || (!s.Maghrib.IsZero() && !s.Isha.After(s.Maghrib)) {
			s.Isha = time.Time{}
		}

		if !s.Asr.After(s.Zuhr) || (!s.Maghrib.IsZero() && !s.Asr.Before(s.Maghrib)) {
			s.Asr = time.Time{}
		}

		// Check if twilight convention use fixed duration for Maghrib
		if cfg.TwilightConvention.MaghribDuration != 0 && !s.Maghrib.IsZero() {
			s.Isha = s.Maghrib.Add(cfg.TwilightConvention.MaghribDuration)
		}

		// Check if the day is abnormal
		if s.Fajr.IsZero() || s.Sunrise.IsZero() || s.Asr.IsZero() || s.Maghrib.IsZero() || s.Isha.IsZero() {
			nAbnormal++
		}

		// Save the adjustment
		schedules[i] = s
	}

	return schedules, nAbnormal
}

func adjustScheduleIdx(schedules []PrayerSchedule, idx int, t, transit time.Time, beforeTransit bool) int {
	// return idx
	if !t.IsZero() && !transit.IsZero() {
		// If event is supposed to occur before transit but in calculation it
		// happened after, then it's event that chained with tomorrow schedules
		//
		// If event is supposed to occur after transit but in calculation it
		// happened before, then it's event that chained with yesterday schedules
		if beforeTransit && t.After(transit) {
			idx++
		} else if !beforeTransit && t.Before(transit) {
			idx--
		}
	}

	// Fix the index
	if idx >= len(schedules) {
		idx = 0
	} else if idx < 0 {
		idx = len(schedules) - 1
	}

	return idx
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func acot(cotValue float64) float64 {
	return math.Atan(1 / cotValue)
}
