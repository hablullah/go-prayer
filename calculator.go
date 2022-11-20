package prayer

import (
	"time"
)

// Times is the result of calculation.
type Times struct {
	Fajr    time.Time
	Sunrise time.Time
	Zuhr    time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
}

// TimeCorrections is correction for each prayer time.
type TimeCorrections struct {
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

	// CalculationMethod is the method that used for calculating Fajr and Isha time. It works by specifying
	// Fajr angle, Isha angle or Maghrib duration following one of the well-known conventions. By default
	// it will use MWL method.
	CalculationMethod CalculationMethod

	// FajrAngle is the altitude of Sun below horizon which mark the start of Fajr time. If it's specified,
	// the Fajr angle that provided by CalculationMethod will be ignored.
	FajrAngle float64

	// IshaAngle is the altitude of Sun below horizon which mark the start of Isha time. If it's specified,
	// the Isha angle that provided by CalculationMethod will be ignored.
	IshaAngle float64

	// MaghribDuration is the duration between Maghrib and Isha. If it's specified, the Maghrib duration
	// that provided by CalculationMethod will be ignored. Isha angle will be ignored as well since the
	// Isha time will be calculated from Maghrib time.
	MaghribDuration time.Duration

	// AsrConvention is the convention that used for calculating Asr time. There are two conventions,
	// Shafii and Hanafi. By default it will use Shafii.
	AsrConvention AsrConvention

	// PreciseToSeconds specify whether output time will omit the seconds or not.
	PreciseToSeconds bool

	// TimeCorrections is used to corrects calculated time for each specified prayer.
	TimeCorrections TimeCorrections

	// HighLatitudeMethods is methods that used for calculating Fajr and Isha time in higher latitude area
	// (more than 45 degree from equator) where the Sun might never set or rise for an entire season. By
	// default it will use angle-based method.
	HighLatitudeMethod HighLatitudeMethod
}

// Calculate calculates the prayer time for specified date with specified configuration.
func Calculate(cfg Config, date time.Time) (Times, error) {
	return Times{}, nil
}
