package prayer

import (
	"time"

	"github.com/RadhiFadlillah/go-prayer/internal/prayer"
	"github.com/shopspring/decimal"
)

// Config is the configuration that used for calculating prayer times.
type Config struct {
	Latitude          float64
	Longitude         float64
	Elevation         float64
	FajrAngle         float64
	IshaAngle         float64
	CalculationMethod CalculationMethod
	AsrJuristicMethod AsrJuristicMethod
	Corrections       TimeCorrections
	Iqamah            IqamahDelay
	PreciseToSeconds  bool
}

// TimeCorrections is an additional object that used to define
// manual correction for each prayer times.
type TimeCorrections struct {
	Fajr    time.Duration
	Sunrise time.Duration
	Zuhr    time.Duration
	Asr     time.Duration
	Maghrib time.Duration
	Isha    time.Duration
}

// IqamahDelay is delay after adhan before iqamah is commenced.
type IqamahDelay struct {
	Fajr    time.Duration
	Zuhr    time.Duration
	Asr     time.Duration
	Maghrib time.Duration
	Isha    time.Duration
}

// Times is time for adhan or iqamah of each prayer.
type Times struct {
	Fajr    time.Time
	Sunrise time.Time
	Zuhr    time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
}

// GetTimes calculates the times of adhan and iqamah for the submitted date.
func GetTimes(date time.Time, cfg Config) (Times, Times) {
	// Parse value from config
	asrCoefficient := getAsrCoefficient(cfg)
	fajrAngle, ishaAngle, maghribDuration := getCalculationAngle(cfg)

	// Prepare calculator
	calc := prayer.TimeCalculator{
		Elevation:        cfg.Elevation,
		Latitude:         decimal.NewFromFloat(cfg.Latitude),
		Longitude:        decimal.NewFromFloat(cfg.Longitude),
		AsrCoefficient:   asrCoefficient,
		FajrAngle:        fajrAngle,
		IshaAngle:        ishaAngle,
		MaghribDuration:  maghribDuration,
		PreciseToSeconds: cfg.PreciseToSeconds}
	calc.SetDate(date)

	// Calculate adhan times and its correction value
	adhan := Times{
		Fajr:    calc.GetFajrTime().Add(cfg.Corrections.Fajr),
		Sunrise: calc.GetSunriseTime().Add(cfg.Corrections.Sunrise),
		Zuhr:    calc.GetZuhrTime().Add(cfg.Corrections.Zuhr),
		Asr:     calc.GetAsrTime().Add(cfg.Corrections.Asr),
		Maghrib: calc.GetMaghribTime().Add(cfg.Corrections.Maghrib),
		Isha:    calc.GetIshaTime().Add(cfg.Corrections.Isha),
	}

	// Calculate iqamah time
	iqamah := Times{
		Fajr:    adhan.Fajr.Add(cfg.Iqamah.Fajr),
		Zuhr:    adhan.Zuhr.Add(cfg.Iqamah.Zuhr),
		Asr:     adhan.Asr.Add(cfg.Iqamah.Asr),
		Maghrib: adhan.Maghrib.Add(cfg.Iqamah.Maghrib),
		Isha:    adhan.Isha.Add(cfg.Iqamah.Isha),
	}

	return adhan, iqamah
}
