package prayer

import (
	"time"

	"github.com/RadhiFadlillah/go-prayer/internal/julianday"
	"github.com/shopspring/decimal"
)

// TimeCorrection is correction for each target time
type TimeCorrection map[Target]time.Duration

// AngleCorrection is value in degree, used to correct hour angle
type AngleCorrection map[Target]float64

// Calculator is calculator that used to calculate the prayer times.
type Calculator struct {
	// Latitude is the latitude of the location. Positive for north
	// area and negative for south area.
	Latitude float64

	// Longitude is the longitude of the location. Positive for east
	// area and negative for west area.
	Longitude float64

	// Elevation is the elevation of the location above sea level.
	// It's used to calculate sunrise and sunset. It's fine to omit
	// this field because the difference will be really small, but
	// it's better to specify it.
	Elevation float64

	// CalculationMethod is the method that used for calculating
	// Fajr and Isha time. It works by specifying Fajr angle, Isha
	// angle or Maghrib duration following one of the well-known
	// conventions. By default it will use MWL method.
	CalculationMethod CalculationMethod

	// FajrAngle is the angle of sun below horizon which mark the
	// start of Fajr time. If it's specified, the Fajr angle that
	// provided by CalculationMethod will be ignored.
	FajrAngle float64

	// IshaAngle is the angle of sun below horizon which mark the
	// start of Isha time. If it's specified, the Isha angle that
	// provided by CalculationMethod will be ignored.
	IshaAngle float64

	// MaghribDuration is the duration between Maghrib and Isha
	// If it's specified, the Maghrib duration that provided by
	// CalculationMethod will be ignored. Isha angle will be
	// ignored as well since the Isha time will be calculated
	// from Maghrib time.
	MaghribDuration time.Duration

	// AsrConvention is the convention that used for calculating
	// Asr time. There are two conventions, Shafii and Hanafi.
	// By default it will use Shafii.
	AsrConvention AsrConvention

	// PreciseToSeconds specify whether output time will omit
	// the seconds or not. If it set to false, the minutes
	// will be rounded up if seconds >= 30, and rounded down
	// if seconds less than 30.
	PreciseToSeconds bool

	// TimeCorrection is map which used to corrects calculated
	// time for each specified target.
	TimeCorrection TimeCorrection

	// AngleCorrection is map which used to corrects hour angle
	// for each specified target. It might be easier to use
	// `TimeCorrection` instead of this field, but some people
	// might prefer this.
	AngleCorrection AngleCorrection

	// HighLatitudeMethods is methods that used for calculating Fajr
	// and Isha time in higher latitude area (more than 48.5 degree
	// from equator) where the sun might never set or rise for an
	// entire season. By default it will use angle-based method.
	HighLatitudeMethod HighLatitudeMethods

	latitude       decimal.Decimal
	longitude      decimal.Decimal
	fajrAngle      decimal.Decimal
	ishaAngle      decimal.Decimal
	asrCoefficient decimal.Decimal

	date           time.Time
	timezone       decimal.Decimal
	transitTime    decimal.Decimal
	sunDeclination decimal.Decimal
}

// Init initiates the calculator. Must be run after changing
// any fields in Calculator.
func (calc *Calculator) Init() *Calculator {
	// Save location
	calc.latitude = decimal.NewFromFloat(calc.Latitude)
	calc.longitude = decimal.NewFromFloat(calc.Longitude)

	// Apply calculation method
	var maghribDuration time.Duration
	var fajrAngle, ishaAngle float64

	switch calc.CalculationMethod {
	case MWL, Algerian, Diyanet:
		fajrAngle, ishaAngle = 18, 17
	case ISNA:
		fajrAngle, ishaAngle = 15, 15
	case UmmAlQura:
		fajrAngle, maghribDuration = 18.5, 90*time.Minute
	case Gulf:
		fajrAngle, maghribDuration = 19.5, 90*time.Minute
	case Karachi, France18, Tunisia:
		fajrAngle, ishaAngle = 18, 18
	case Egypt:
		fajrAngle, ishaAngle = 19.5, 17.5
	case EgyptBis, Kemenag, MUIS, JAKIM:
		fajrAngle, ishaAngle = 20, 18
	case UOIF:
		fajrAngle, ishaAngle = 12, 12
	case France15:
		fajrAngle, ishaAngle = 15, 15
	case Tehran:
		fajrAngle, ishaAngle = 17.7, 14
	case Jafari:
		fajrAngle, ishaAngle = 16, 14
	}

	if calc.FajrAngle != 0 {
		fajrAngle = calc.FajrAngle
	}

	if calc.IshaAngle != 0 {
		ishaAngle = calc.IshaAngle
	}

	if calc.MaghribDuration != 0 {
		maghribDuration = calc.MaghribDuration
	}

	calc.fajrAngle = decimal.NewFromFloat(fajrAngle)
	calc.ishaAngle = decimal.NewFromFloat(ishaAngle)
	calc.MaghribDuration = maghribDuration

	// Set asr coefficient
	switch calc.AsrConvention {
	case Hanafi:
		calc.asrCoefficient = decimal.New(2, 0)
	default:
		calc.asrCoefficient = decimal.New(1, 0)
	}

	return calc
}

// SetDate specifies active date to calculate.
// It will also calculates the timezone from the date location.
func (calc *Calculator) SetDate(date time.Time) *Calculator {
	// Make sure date is at 12 local time
	y := date.Year()
	m := date.Month()
	d := date.Day()
	location := date.Location()
	calc.date = time.Date(y, m, d, 12, 0, 0, 0, location)

	// Save timezone
	_, utcOffset := calc.date.Zone()
	calc.timezone = decimal.New(int64(utcOffset), 0).
		Div(decimal.New(3600, 0))

	// Calculate transit time and sun declination
	jd := julianday.Convert(calc.date)
	calc.transitTime = calc.getTransitTime(jd)
	calc.sunDeclination = calc.getSunDeclination(jd)
	return calc
}

// Calculate calculates times for all possible targets. If the target
// is not available, it will be omitted from result.
func (calc Calculator) Calculate() map[Target]time.Time {
	times := map[Target]time.Time{}

	// Get all target's time
	for target := Fajr; target <= Isha; target++ {
		if targetTime, isNA := calc.calculate(target); !isNA {
			times[target] = targetTime
		}
	}

	// If Fajr or Isha is not calculable
	// but Sunrise and Sunset is, use high latitude rules
	_, hasFajr := times[Fajr]
	_, hasIsha := times[Isha]
	sunrise, hasSunrise := times[Sunrise]
	sunset, hasSunset := times[Maghrib]
	if (!hasFajr || !hasIsha) && (hasSunrise && hasSunset) {
		times[Fajr], times[Isha] = calc.adjustHighLatitudeTime(sunrise, sunset)
	}

	return times
}
