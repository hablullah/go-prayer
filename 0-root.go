package prayer

import (
	"math"
	"time"
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

	// Calculate the schedules
	schedules, nAbnormal := calcNormal(cfg, year)

	// Apply high latitude convention
	if math.Abs(cfg.Latitude) > 45 || nAbnormal > 0 {
		switch cfg.HighLatConvention {
		case Mecca:
			schedules = calcHighLatMecca(cfg, year, schedules)
		case ForceMecca:
			schedules = calcHighLatForceMecca(cfg, year, schedules)
		case LocalRelativeEstimation:
			schedules = calcLocalRelativeEstimation(cfg, year, schedules)
		case NearestDay:
			schedules = calcHighLatNearestDay(schedules)
		case NearestLatitude:
			schedules = calcHighLatNearestLatitude(cfg, year, schedules)
		case ForceNearestLatitude:
			schedules = calcHighLatForceNearestLatitude(cfg, year)
		case ShariNormalDay:
			schedules = calcHighLatShariNormalDay(cfg, year, schedules)
		case AngleBased:
			schedules = calcHighLatAngleBased(cfg, schedules)
		case OneSeventhNight:
			schedules = calcHighLatOneSeventhNight(schedules)
		case MiddleNight:
			schedules = calcHighLatMiddleNight(schedules)
		}
	}

	// Apply time correction
	for i, s := range schedules {
		s.Fajr = applyCorrection(s.Fajr, cfg.Corrections.Fajr)
		s.Sunrise = applyCorrection(s.Sunrise, cfg.Corrections.Sunrise)
		s.Zuhr = applyCorrection(s.Zuhr, cfg.Corrections.Zuhr)
		s.Asr = applyCorrection(s.Asr, cfg.Corrections.Asr)
		s.Maghrib = applyCorrection(s.Maghrib, cfg.Corrections.Maghrib)
		s.Isha = applyCorrection(s.Isha, cfg.Corrections.Isha)
		schedules[i] = s
	}

	return schedules, nil
}

func applyCorrection(t time.Time, d time.Duration) time.Time {
	if !t.IsZero() {
		t = t.Add(d)
	}
	return t
}
