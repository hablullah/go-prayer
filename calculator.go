package prayer

import (
	"math"
	"time"

	"github.com/hablullah/go-sampa"
)

// PrayerSchedule is the result of prayer time calculation.
type PrayerSchedule struct {
	// Date is the ISO date, useful for logging.
	Date string

	// Fajr is the time when the sky begins to lighten (dawn) after previously completely dark.
	Fajr time.Time

	// Sunrise is the moment when the upper limb of the Sun appears on the horizon in the morning.
	Sunrise time.Time

	// Zuhr is the time when the Sun begins to decline after reaching the highest point in the sky,
	// so a bit after solar noon.
	Zuhr time.Time

	// Asr is the time when the length of any object's shadow reaches a factor of the length of the
	// object itself plus the length of that object's shadow at noon.
	Asr time.Time

	// Maghrib is sunset, i.e. the time when the upper limb of the Sun disappears below the horizon.
	Maghrib time.Time

	// Isha is the time when darkness falls and after this point the sky is no longer illuminated (dusk).
	Isha time.Time
}

// ScheduleCorrections is correction for each prayer time.
type ScheduleCorrections struct {
	Fajr    time.Duration
	Sunrise time.Duration
	Zuhr    time.Duration
	Asr     time.Duration
	Maghrib time.Duration
	Isha    time.Duration
}

// Config is configuration that used to calculate the prayer times.
type Config struct {
	// Latitude is the latitude of the location. Positive for north area and negative for south area.
	Latitude float64

	// Longitude is the longitude of the location. Positive for east area and negative for west area.
	Longitude float64

	// Elevation is the elevation of the location above sea level. It's used to improve calculation for
	// sunrise and sunset by factoring the value of atmospheric refraction. However, apparently most of
	// the prayer time calculator doesn't use it so it's fine to omit it.
	Elevation float64

	// Timezone is the time zone of the location.
	Timezone *time.Location

	// TwilightConvention is the convention that used to specify time for Fajr and Isha. By default
	// it will use `AstronomicalTwilight`.
	TwilightConvention *TwilightConvention

	// AsrConvention is the convention that used for calculating Asr time. There are two conventions,
	// Shafii and Hanafi. By default it will use Shafii.
	AsrConvention AsrConvention

	// HighLatConvention is the convention for calculation prayer times in area with high latitude
	// (>=48 degrees). By default it will use `LocalRelativeEstimation`.
	HighLatConvention HighLatConvention

	// Corrections is used to corrects calculated time for each specified prayer.
	Corrections ScheduleCorrections

	// PreciseToSeconds specify whether output time will omit the seconds or not.
	PreciseToSeconds bool
}

// Calculate calculates the prayer time for the entire year with specified configuration.
func Calculate(cfg Config, year int) ([]PrayerSchedule, error) {
	// Apply default config
	if cfg.TwilightConvention == nil {
		cfg.TwilightConvention = AstronomicalTwilight
	}

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
	var nAbnormalDays int
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

		// Apply time correction
		s.Fajr = applyCorrection(s.Fajr, cfg.Corrections.Fajr)
		s.Sunrise = applyCorrection(s.Sunrise, cfg.Corrections.Sunrise)
		s.Zuhr = applyCorrection(s.Zuhr, cfg.Corrections.Zuhr)
		s.Asr = applyCorrection(s.Asr, cfg.Corrections.Asr)
		s.Maghrib = applyCorrection(s.Maghrib, cfg.Corrections.Maghrib)
		s.Isha = applyCorrection(s.Isha, cfg.Corrections.Isha)

		// Check if the day is abnormal
		if s.Fajr.IsZero() || s.Sunrise.IsZero() || s.Maghrib.IsZero() || s.Isha.IsZero() {
			nAbnormalDays++
		}

		// Save the adjustment
		schedules[i] = s
	}

	// Calculate time for high latitude
	if nAbnormalDays > 0 && cfg.HighLatConvention > Disabled {
		// TODO
	}

	return schedules, nil
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

func applyCorrection(t time.Time, d time.Duration) time.Time {
	if !t.IsZero() {
		t = t.Add(d)
	}
	return t
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
