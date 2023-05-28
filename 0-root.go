package prayer

import (
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

	// IsNormal specify whether the day have a normal day night period or not. It will be false in
	// area with higher latitude, when Sun never rise or set in extreme periods.
	IsNormal bool
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

// TwilightConvention is the convention that specifies time for Fajr (dawn) and Isha (dusk). Most of
// the conventions use Solar angle elevation for both dawn and dusk time, however there are several
// convention where dusk times depends on sunset (Maghrib) times.
type TwilightConvention struct {
	FajrAngle       float64
	IshaAngle       float64
	MaghribDuration time.Duration
}

// HighLatitudeAdapter is function for calculating prayer times in area with latitude >45 degrees.
// Check out https://www.prayertimes.dk/story.html for why this is needed.
type HighLatitudeAdapter func(cfg Config, year int, currentSchedules []PrayerSchedule) []PrayerSchedule

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

	// HighLatitudeAdapter is the function for adjusting prayer times in area with high latitude
	// (>=45 degrees). If not specified, it will not calculate the adjustment for higher latitude
	// and instead will return the schedule as it is. For area in high or extreme latitude, it might
	// return zero for Fajr, Sunrise, Maghrib and Isha.
	HighLatitudeAdapter HighLatitudeAdapter

	// Corrections is used to corrects calculated time for each specified prayer.
	Corrections ScheduleCorrections

	// PreciseToSeconds specify whether output time will omit the seconds or not.
	PreciseToSeconds bool
}

// Calculate calculates the prayer time for the entire year with specified configuration.
func Calculate(cfg Config, year int) ([]PrayerSchedule, error) {
	// Apply default config
	if cfg.TwilightConvention == nil {
		cfg.TwilightConvention = AstronomicalTwilight()
	}

	// Calculate the schedules
	schedules, nAbnormal := calcNormal(cfg, year)

	// Apply high latitude adapter
	if nAbnormal > 0 && cfg.HighLatitudeAdapter != nil {
		schedules = cfg.HighLatitudeAdapter(cfg, year, schedules)
	}

	// Apply Isha times for convention where Isha time is fixed after Maghrib
	if d := cfg.TwilightConvention.MaghribDuration; d > 0 {
		for i, s := range schedules {
			schedules[i].Isha = s.Maghrib.Add(d)
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
