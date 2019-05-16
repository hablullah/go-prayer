package prayer

import (
	"math"
	"time"

	"github.com/shopspring/decimal"
)

// CalculationMethod is the conventions for calculating prayer times, especially Fajr and Isha.
// For details, check http://praytimes.org/wiki/Calculation_Methods.
type CalculationMethod int

// AsrCalculationMethod is the conventions for calculating Asr time.
// For details, check http://www.prayerminder.com/faq.php#Fiqh.
type AsrCalculationMethod int

const (
	// Default is the default calculation method, with Fajr at 20° and Isha at 18°.
	Default CalculationMethod = iota
	// MWL is calculation method from Muslim World League, with Fajr at 18° and Isha at 17°.
	// Usually used in Europe, Far East and parts of US.
	MWL
	// ISNA is calculation method from Islamic Society of North America, with both Fajr and Isha at 15°.
	// Used in North America i.e US and Canada.
	ISNA
	// Egypt is calculation method from Egyptian General Authority of Survey, with Fajr at 19.5° and Isha at 17.5°.
	// Used in Africa, Syria, Lebanon and Malaysia.
	Egypt
	// Karachi is calculation method from University of Islamic Sciences, Karachi, with both Fajr and Isha at 18°.
	// Used in Pakistan, Afganistan, Bangladesh and India.
	Karachi
)

const (
	// Hanafi is the school which said that the Asr time is when the shadow of
	// an object is twice the length of the object plus the length of its shadow
	// when the sun is at its zenith.
	Hanafi AsrCalculationMethod = iota
	// Shafii is the school which said that the Asr time is when the shadow of
	// an object is equals the length of the object plus the length of its shadow
	// when the sun is at its zenith.
	Shafii
)

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

// IqamahDuration is an additional object that used to define
// duration between adhan and iqamah.
type IqamahDuration struct {
	Fajr    time.Duration
	Zuhr    time.Duration
	Asr     time.Duration
	Maghrib time.Duration
	Isha    time.Duration
}

// Times is time for adhan and iqamah of each prayer.
type Times struct {
	Fajr    time.Time
	Sunrise time.Time
	Zuhr    time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time

	IqamahFajr    time.Time
	IqamahZuhr    time.Time
	IqamahAsr     time.Time
	IqamahMaghrib time.Time
	IqamahIsha    time.Time
}

// Config is the configuration that used for calculating prayer times.
type Config struct {
	Latitude             float64
	Longitude            float64
	Elevation            float64
	FajrAngle            float64
	IshaAngle            float64
	CalculationMethod    CalculationMethod
	AsrCalculationMethod AsrCalculationMethod
	Corrections          TimeCorrections
	IqamahDuration       IqamahDuration
	PreciseToSeconds     bool
}

// Calculate calculates the times of prayers on the submitted date.
func Calculate(date time.Time, cfg Config) Times {
	// Make sure date is at 12 local time
	date = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		12, 0, 0, 0,
		date.Location())

	// Calculate variables
	jd := getJulianDay(date)
	timezone := getTimezone(date)

	eoT := getEquationOfTime(jd)
	sunDeclination := getSunDeclination(jd)
	transitTime := getTransitTime(timezone, cfg.Longitude, eoT)

	// Calculate adhan times
	fajr := getFajrTime(date, transitTime, sunDeclination, cfg)
	sunrise := getSunriseTime(date, transitTime, sunDeclination, cfg)
	zuhr := getZuhrTime(date, transitTime, cfg)
	asr := getAsrTime(date, transitTime, sunDeclination, cfg)
	maghrib := getMaghribTime(date, transitTime, sunDeclination, cfg)
	isha := getIshaTime(date, transitTime, sunDeclination, cfg)

	// Merge the correction time
	fajr = fajr.Add(cfg.Corrections.Fajr)
	sunrise = sunrise.Add(cfg.Corrections.Sunrise)
	zuhr = zuhr.Add(cfg.Corrections.Zuhr)
	asr = asr.Add(cfg.Corrections.Asr)
	maghrib = maghrib.Add(cfg.Corrections.Maghrib)
	isha = isha.Add(cfg.Corrections.Isha)

	// Calculate iqamah time
	iqamahFajr := fajr.Add(cfg.IqamahDuration.Fajr)
	iqamahZuhr := zuhr.Add(cfg.IqamahDuration.Zuhr)
	iqamahAsr := asr.Add(cfg.IqamahDuration.Asr)
	iqamahMaghrib := maghrib.Add(cfg.IqamahDuration.Maghrib)
	iqamahIsha := isha.Add(cfg.IqamahDuration.Isha)

	return Times{
		Fajr:    fajr,
		Sunrise: sunrise,
		Zuhr:    zuhr,
		Asr:     asr,
		Maghrib: maghrib,
		Isha:    isha,

		IqamahFajr:    iqamahFajr,
		IqamahZuhr:    iqamahZuhr,
		IqamahAsr:     iqamahAsr,
		IqamahMaghrib: iqamahMaghrib,
		IqamahIsha:    iqamahIsha,
	}
}

func getFajrAngle(cfg Config) float64 {
	if cfg.FajrAngle != 0 {
		return cfg.FajrAngle
	}

	switch cfg.CalculationMethod {
	case MWL:
		return 18.0
	case ISNA:
		return 15.0
	case Egypt:
		return 19.5
	case Karachi:
		return 18.0
	default:
		return 20.0
	}
}

func getIshaAngle(cfg Config) float64 {
	if cfg.IshaAngle != 0 {
		return cfg.IshaAngle
	}

	switch cfg.CalculationMethod {
	case MWL:
		return 17.0
	case ISNA:
		return 15.0
	case Egypt:
		return 17.5
	case Karachi:
		return 18.0
	default:
		return 18.0
	}
}

func getAsrCoefficient(cfg Config) float64 {
	switch cfg.AsrCalculationMethod {
	case Hanafi:
		return 2.0
	default:
		return 1.0
	}
}

func getFajrTime(date time.Time, transitTime, sunDecli decimal.Decimal, cfg Config) time.Time {
	fajrAngle := getFajrAngle(cfg)
	sunAltitude := decimal.NewFromFloat(-fajrAngle)

	hourAngle := getHourAngle(cfg.Latitude, sunAltitude, sunDecli)
	hours := transitTime.Sub(hourAngle.Div(decimal.New(15, 0)))

	return attachHours(hours, date, cfg)
}

func getSunriseTime(date time.Time, transitTime, sunDecli decimal.Decimal, cfg Config) time.Time {
	A := decimal.New(-5, 0).Div(decimal.New(6, 0))
	B := decimal.NewFromFloat(0.0347).Mul(decimal.NewFromFloat(math.Sqrt(cfg.Elevation)))
	sunAltitude := A.Sub(B)

	hourAngle := getHourAngle(cfg.Latitude, sunAltitude, sunDecli)
	hours := transitTime.Sub(hourAngle.Div(decimal.New(15, 0)))

	return attachHours(hours, date, cfg)
}

func getZuhrTime(date time.Time, transitTime decimal.Decimal, cfg Config) time.Time {
	twoMinutes := decimal.New(2, 0).Div(decimal.New(60, 0))
	hours := transitTime.Add(twoMinutes)

	return attachHours(hours, date, cfg)
}

func getAsrTime(date time.Time, transitTime, sunDecli decimal.Decimal, cfg Config) time.Time {
	latitude := decimal.NewFromFloat(cfg.Latitude)
	asrCoefficient := decimal.NewFromFloat(getAsrCoefficient(cfg))
	sunAltitude := acot(asrCoefficient.Add(tan(sunDecli.Sub(latitude).Abs())))

	hourAngle := getHourAngle(cfg.Latitude, sunAltitude, sunDecli)
	hours := transitTime.Add(hourAngle.Div(decimal.New(15, 0)))

	return attachHours(hours, date, cfg)
}

func getMaghribTime(date time.Time, transitTime, sunDecli decimal.Decimal, cfg Config) time.Time {
	A := decimal.New(-5, 0).Div(decimal.New(6, 0))
	B := decimal.NewFromFloat(0.0347).Mul(decimal.NewFromFloat(math.Sqrt(cfg.Elevation)))
	sunAltitude := A.Sub(B)

	hourAngle := getHourAngle(cfg.Latitude, sunAltitude, sunDecli)
	hours := transitTime.Add(hourAngle.Div(decimal.New(15, 0)))

	return attachHours(hours, date, cfg)
}

func getIshaTime(date time.Time, transitTime, sunDecli decimal.Decimal, cfg Config) time.Time {
	ishaAngle := getIshaAngle(cfg)
	sunAltitude := decimal.NewFromFloat(-ishaAngle)

	hourAngle := getHourAngle(cfg.Latitude, sunAltitude, sunDecli)
	hours := transitTime.Add(hourAngle.Div(decimal.New(15, 0)))

	return attachHours(hours, date, cfg)
}

func attachHours(hours decimal.Decimal, date time.Time, cfg Config) time.Time {
	var minutes, seconds int

	if cfg.PreciseToSeconds {
		seconds = int(hours.Mul(decimal.New(3600, 0)).Ceil().IntPart())
	} else {
		minutes = int(hours.Mul(decimal.New(60, 0)).Ceil().IntPart())
	}

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, minutes, seconds, 0,
		date.Location())
}
