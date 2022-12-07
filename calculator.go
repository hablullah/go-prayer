package prayer

import (
	"math"
	"time"

	"github.com/hablullah/go-sampa"
)

// Times is the result of calculation.
type Times struct {
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

	// TimeCorrections is used to corrects calculated time for each specified prayer.
	TimeCorrections TimeCorrections

	// PreciseToSeconds specify whether output time will omit the seconds or not.
	PreciseToSeconds bool
}

// Calculate calculates the prayer time for the entire year with specified configuration.
func Calculate(cfg Config, year int) ([]Times, error) {
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

	// Loop for every day of the year
	var prayerTimes []Times
	for dt := time.Date(year, 1, 1, 12, 0, 0, 0, cfg.Timezone); dt.Year() == year; dt = dt.AddDate(0, 0, 1) {
		// Calculate Sun events
		events, err := sampa.GetSunEvents(dt, location, nil, customEvents...)
		if err != nil {
			return nil, err
		}

		// Extract the time
		fajr := events.Others["fajr"].DateTime
		sunrise := events.Sunrise.DateTime
		zuhr := events.Transit.DateTime
		asr := events.Others["asr"].DateTime
		maghrib := events.Sunset.DateTime
		isha := events.Others["isha"].DateTime

		// Save the times
		prayerTimes = append(prayerTimes, Times{
			Date:    dt.Format("2006-01-02"),
			Fajr:    applyTc(fajr, cfg.TimeCorrections.Fajr),
			Sunrise: applyTc(sunrise, cfg.TimeCorrections.Sunrise),
			Zuhr:    applyTc(zuhr, cfg.TimeCorrections.Zuhr),
			Asr:     applyTc(asr, cfg.TimeCorrections.Asr),
			Maghrib: applyTc(maghrib, cfg.TimeCorrections.Maghrib),
			Isha:    applyTc(isha, cfg.TimeCorrections.Isha),
		})

		// TODO: Calculate time for high latitude
	}

	return prayerTimes, nil
}

func applyTc(t time.Time, d time.Duration) time.Time {
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
