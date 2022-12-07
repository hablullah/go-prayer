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

	// TwilightConvention is the convention that used to specify time for Fajr and Isha. By default
	// it will use `AstronomicalTwilight`.
	TwilightConvention *TwilightConvention

	// AsrConvention is the convention that used for calculating Asr time. There are two conventions,
	// Shafii and Hanafi. By default it will use Shafii.
	AsrConvention AsrConvention

	// HighLatConvention is the convention for calculation prayer times in area with high latitude
	// (>=48 degrees). By default it will use `LocalRelativeEstimation`.
	HighLatConvention HighLatConvention

	// TimeCorrections is used to corrects calculated time for each specified prayer.
	TimeCorrections TimeCorrections

	// PreciseToSeconds specify whether output time will omit the seconds or not.
	PreciseToSeconds bool
}

// Calculate calculates the prayer time for specified date with specified configuration.
func Calculate(cfg Config, date time.Time) (Times, error) {
	return Times{}, nil
}
