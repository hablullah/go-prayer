package prayer

import (
	"math"
	"time"
)

// Type is the type of prayer, i.e Fajr, Zuhr, Asr, Maghrib and Isha.
type Type int

// CalculationMethod is the conventions for calculating prayer times, especially Fajr and Isha.
// For details, check http://praytimes.org/wiki/Calculation_Methods.
type CalculationMethod int

// AsrCalculationMethod is the conventions for calculating Asr time.
// For details, check http://www.prayerminder.com/faq.php#Fiqh.
type AsrCalculationMethod int

// TimeCorrections is the correction time for each prayer.
type TimeCorrections map[Type]time.Duration

// Times is the time for adhan or iqama.
type Times map[Type]time.Time

const (
	// Fajr is prayer that done at dawn to sunrise.
	Fajr Type = iota
	// Sunrise is sunrise. It's only included as prayer type to simplify code.
	Sunrise
	// Zuhr is prayer that done after true noon until Asr.
	Zuhr
	// Asr is prayer that done in afternoon, when the sun started to go down.
	Asr
	// Maghrib is prayer that done after sun until dusk.
	Maghrib
	// Isha is prayer that done at dusk and until dawn.
	Isha
)

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

// Calculator is object for calculating prayer times.
type Calculator struct {
	Latitude             float64
	Longitude            float64
	Elevation            float64
	CalculationMethod    CalculationMethod
	AsrCalculationMethod AsrCalculationMethod
	AdhanCorrections     TimeCorrections
	TimesToIqama         TimeCorrections
	FajrAngle            float64
	IshaAngle            float64

	sunDeclination float64
	equationOfTime float64
	transitTime    float64
}

// Calculate calculates the times of prayers on the submitted date.
func (calc *Calculator) Calculate(date time.Time) (Times, Times) {
	// Normalize date
	date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())

	// Get timezone
	_, offset := date.Zone()
	timezone := offset / 3600

	// Prepare variables
	var fajrAngle, ishaAngle, asrCoefficient float64

	julianDay := getJulianDay(date)
	calc.sunDeclination = getSunDeclination(julianDay)
	calc.equationOfTime = getEquationOfTime(julianDay)
	calc.transitTime = 12.0 + float64(timezone) - calc.Longitude/15 - calc.equationOfTime/60

	// Calculate fajr and isha angle
	switch calc.CalculationMethod {
	case MWL:
		fajrAngle = 18.0
		ishaAngle = 17.0
	case ISNA:
		fajrAngle = 15.0
		ishaAngle = 15.0
	case Egypt:
		fajrAngle = 19.5
		ishaAngle = 17.5
	case Karachi:
		fajrAngle = 18.0
		ishaAngle = 18.0
	default:
		fajrAngle = 20.0
		ishaAngle = 18.0
	}

	// If fajr and isha angle specified, use that
	if calc.FajrAngle != 0 {
		fajrAngle = calc.FajrAngle
	}

	if calc.IshaAngle != 0 {
		ishaAngle = calc.IshaAngle
	}

	// Calculate asr coefficient
	switch calc.AsrCalculationMethod {
	case Hanafi:
		asrCoefficient = 2.0
	default:
		asrCoefficient = 1.0
	}

	// Calculate sun altitude
	altitudes := map[Type]float64{
		Fajr:    -fajrAngle,
		Sunrise: -(5.0 / 6.0),
		Zuhr:    0,
		Asr:     arccot(asrCoefficient + tan(math.Abs(calc.sunDeclination-calc.Latitude))),
		Maghrib: -(5.0 / 6.0),
		Isha:    -ishaAngle,
	}

	// Calculate adhan and iqama
	adhanTimes := make(Times)
	iqamaTimes := make(Times)
	for prayer, altitude := range altitudes {
		adhan := calc.getAdhanTime(date, altitude, prayer)
		iqama := adhan.Add(calc.TimesToIqama[prayer])

		adhanTimes[prayer] = adhan
		if prayer != Sunrise {
			iqamaTimes[prayer] = iqama
		}
	}

	return adhanTimes, iqamaTimes
}

func (calc Calculator) getAdhanTime(date time.Time, altitude float64, prayer Type) time.Time {
	// Calculate hour angle
	hourAngle := 0.0
	if prayer != Zuhr {
		hourAngle = acos(
			(sin(altitude) - sin(calc.Latitude)*sin(calc.sunDeclination)) /
				(cos(calc.Latitude) * cos(calc.sunDeclination)))

		if prayer == Fajr || prayer == Sunrise {
			hourAngle *= -1
		}
	}

	// Calculate adhan time
	minutes := math.Floor((calc.transitTime+hourAngle/15.0)*60 + 0.5)
	year := date.Year()
	month := date.Month()
	day := date.Day()
	location := date.Location()

	return time.Date(year, month, day, 0, int(minutes), 0, 0, location).
		Add(calc.AdhanCorrections[prayer])
}
