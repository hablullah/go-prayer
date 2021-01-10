package prayer

import (
	"math"
	"time"

	"github.com/hablullah/go-juliandays"
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

	// FajrAngle is the angle of sun below horizon which mark the start of Fajr time. If it's specified,
	// the Fajr angle that provided by CalculationMethod will be ignored.
	FajrAngle float64

	// IshaAngle is the angle of sun below horizon which mark the start of Isha time. If it's specified,
	// the Isha angle that provided by CalculationMethod will be ignored.
	IshaAngle float64

	// MaghribDuration is the duration between Maghrib and Isha If it's specified, the Maghrib duration
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
	// (more than 45 degree from equator) where the sun might never set or rise for an entire season. By
	// default it will use angle-based method.
	HighLatitudeMethod HighLatitudeMethod
}

func Calculate(cfg Config, date time.Time) (Times, error) {
	cfg = adjustHighLatitudeConfig(cfg)

	times, err := calculate(cfg, date)
	if err != nil {
		return Times{}, err
	}

	times = applyTimeCorrections(cfg, times)
	return times, nil
}

func calculate(cfg Config, date time.Time) (Times, error) {
	// For initial calculation, get prayer times using noon as base time
	date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())
	times, err := calculateByBase(cfg, date)
	if err != nil {
		return Times{}, err
	}

	// To increase accuracy, redo calculation for each prayer using initial time as base.
	// The increased accuracy is not really much though, only around 1-15 seconds.
	if !times.Fajr.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Fajr)
		times.Fajr = accTimes.Fajr
	}

	if !times.Sunrise.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Sunrise)
		times.Sunrise = accTimes.Sunrise
	}

	if !times.Zuhr.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Zuhr)
		times.Zuhr = accTimes.Zuhr
	}

	if !times.Asr.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Asr)
		times.Asr = accTimes.Asr
	}

	if !times.Maghrib.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Maghrib)
		times.Maghrib = accTimes.Maghrib
	}

	if !times.Isha.IsZero() {
		accTimes, _ := calculateByBase(cfg, times.Isha)
		times.Isha = accTimes.Isha
	}

	// Adjust prayer time in higher latitude
	times = adjustHighLatitudeTimes(cfg, times)
	return times, nil
}

func calculateByBase(cfg Config, baseTime time.Time) (Times, error) {
	// Calculate Julian Days
	jd, err := juliandays.FromTime(baseTime)
	if err != nil {
		return Times{}, err
	}

	// Convert Julian Days to Julian Century
	jc := (jd - 2451545) / 36525

	// Get timezone offset from date
	_, utcOffset := baseTime.Zone()
	timezone := float64(utcOffset) / 3600

	// Calculate position of the sun
	earthOrbitEccent := 0.016708634 - jc*(0.000042037+0.0000001267*jc)
	sunMeanLongitude := math.Mod(280.46646+jc*(36000.76983+jc*0.0003032), 360)
	sunMeanAnomaly := 357.52911 + jc*(35999.05029-0.0001537*jc)
	sunEqOfCenter := sin(sunMeanAnomaly)*(1.914602-jc*(0.004817+0.000014*jc)) +
		sin(2*sunMeanAnomaly)*(0.019993-0.000101*jc) +
		sin(3*sunMeanAnomaly)*0.000289
	sunTrueLongitude := sunMeanLongitude + sunEqOfCenter
	sunAppLongitude := sunTrueLongitude - 0.00569 - 0.00478*sin(125.04-1934.136*jc)
	meanObliqEcliptic := 23 + (26+(21.448-jc*(46.815+jc*(0.00059-jc*0.001813)))/60)/60
	obliqCorrection := meanObliqEcliptic + 0.00256*cos(125.04-1934.136*jc)
	sunDeclination := asin(sin(obliqCorrection) * sin(sunAppLongitude))

	// Calculate equation of time
	tmp := tan(obliqCorrection/2) * tan(obliqCorrection/2)
	eqOfTime := 4 * degree(tmp*sin(2*sunMeanLongitude)-
		2*earthOrbitEccent*sin(sunMeanAnomaly)+
		4*earthOrbitEccent*tmp*sin(sunMeanAnomaly)*cos(2*sunMeanLongitude)-
		0.5*math.Pow(tmp, 2)*sin(4*sunMeanLongitude)-
		1.25*math.Pow(earthOrbitEccent, 2)*sin(2*sunMeanAnomaly))

	// Calculate solar noon
	solarNoon := 720 - 4*cfg.Longitude - eqOfTime + float64(timezone)*60

	// Calculate sunrise and sunset (Maghrib)
	sunriseSunAltitude := -0.833333 - 0.0347*math.Sqrt(cfg.Elevation)
	haSunrise := getHourAngle(cfg.Latitude, sunriseSunAltitude, sunDeclination)
	sunriseTime := solarNoon - haSunrise*4
	maghribTime := solarNoon + haSunrise*4

	// Calculate Fajr and Isha time
	fajrAngle, ishaAngle, maghribDuration := getNightPrayerConfig(cfg)

	fajrSunAltitude := -fajrAngle
	haFajr := getHourAngle(cfg.Latitude, fajrSunAltitude, sunDeclination)
	fajrTime := solarNoon - haFajr*4

	var ishaTime float64
	if maghribDuration != 0 {
		ishaTime = maghribTime + maghribDuration.Minutes()
	} else {
		ishaSunAltitude := -ishaAngle
		haIsha := getHourAngle(cfg.Latitude, ishaSunAltitude, sunDeclination)
		ishaTime = solarNoon + haIsha*4
	}

	// Calculate Asr time
	asrSunAltitude := acot(getAsrCoefficient(cfg) + tan(math.Abs(sunDeclination-cfg.Latitude)))
	haAsr := getHourAngle(cfg.Latitude, asrSunAltitude, sunDeclination)
	asrTime := solarNoon + haAsr*4

	// Return all times
	return Times{
		Fajr:    minutesToTime(cfg, baseTime, fajrTime),
		Sunrise: minutesToTime(cfg, baseTime, sunriseTime),
		Zuhr:    minutesToTime(cfg, baseTime, solarNoon),
		Asr:     minutesToTime(cfg, baseTime, asrTime),
		Maghrib: minutesToTime(cfg, baseTime, maghribTime),
		Isha:    minutesToTime(cfg, baseTime, ishaTime),
	}, nil
}

func getNightPrayerConfig(cfg Config) (fajrAngle, ishaAngle float64, maghribDuration time.Duration) {
	switch cfg.CalculationMethod {
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

	if cfg.FajrAngle != 0 {
		fajrAngle = cfg.FajrAngle
	}

	if cfg.IshaAngle != 0 {
		ishaAngle = cfg.IshaAngle
	}

	if cfg.MaghribDuration != 0 {
		maghribDuration = cfg.MaghribDuration
	}

	return
}

func getAsrCoefficient(cfg Config) float64 {
	if cfg.AsrConvention == Hanafi {
		return 2
	}
	return 1
}

func getHourAngle(latitude, sunAltitude, sunDeclination float64) float64 {
	return acos(
		(sin(float64(sunAltitude)) - sin(latitude)*sin(sunDeclination)) /
			(cos(latitude) * cos(sunDeclination)),
	)
}

func minutesToTime(cfg Config, date time.Time, minutes float64) time.Time {
	if math.IsNaN(minutes) {
		return time.Time{}
	}

	y := date.Year()
	m := date.Month()
	d := date.Day()
	location := date.Location()

	if cfg.PreciseToSeconds {
		seconds := math.Round(minutes * 60)
		return time.Date(y, m, d, 0, 0, int(seconds), 0, location)
	} else {
		minutes = math.Round(minutes)
		return time.Date(y, m, d, 0, int(minutes), 0, 0, location)
	}
}

func adjustHighLatitudeConfig(cfg Config) Config {
	// If high latitude convention is forced normal region, any latitude above (or below) 45 N(S)
	// will be changed to 45
	if cfg.HighLatitudeMethod == ForcedNormalRegion {
		if cfg.Latitude > 45 {
			cfg.Latitude = 45
		} else if cfg.Latitude < -45 {
			cfg.Latitude = -45
		}
	}

	return cfg
}

func adjustHighLatitudeTimes(cfg Config, times Times) Times {
	switch cfg.HighLatitudeMethod {
	case NormalRegion:
		// This adjustment only used in area above latitude 45 (north and south)
		if math.Abs(cfg.Latitude) <= 45 {
			return times
		}

		// Make sure the fasting time is outside normal
		// The normal fasting duration is between 10h 17m and 17h 36m
		fastingDuration := times.Maghrib.Sub(times.Fajr).Minutes()
		if fastingDuration >= 617 && fastingDuration <= 1056 {
			return times
		}

		// Convert latitude to 45
		if cfg.Latitude > 0 {
			cfg.Latitude = 45
		} else {
			cfg.Latitude = -45
		}

		adjustedTimes, _ := calculate(cfg, times.Zuhr)
		return adjustedTimes

	case AngleBased, OneSeventhNight, MiddleNight:
		// These conventions requires sunrise and sunset to be exist
		if times.Sunrise.IsZero() || times.Maghrib.IsZero() {
			return times
		}

		// This adjustment only used in latitude between 48.6 and 66.6 north and south
		absLatitude := math.Abs(cfg.Latitude)
		if absLatitude < 48.6 || absLatitude > 66.6 {
			return times
		}

		// Get Fajr and Isha angle
		fajrAngle, ishaAngle, maghribDuration := getNightPrayerConfig(cfg)

		// Calculate night duration
		dayDuration := times.Maghrib.Sub(times.Sunrise).Minutes()
		nightDuration := (24 * 60) - dayDuration

		// Calculate night portion
		var fajrPortion, ishaPortion float64
		switch cfg.HighLatitudeMethod {
		case MiddleNight:
			fajrPortion = 0.5
			ishaPortion = 0.5
		case OneSeventhNight:
			fajrPortion = 1 / 7
			ishaPortion = 1 / 7
		default:
			fajrPortion = fajrAngle / 60
			ishaPortion = ishaAngle / 60
		}

		// Calculate new Fajr time
		fajrDuration := math.Round(fajrPortion * nightDuration)
		newFajr := times.Sunrise.Add(time.Duration(-fajrDuration) * time.Minute)
		if times.Fajr.IsZero() || newFajr.After(times.Fajr) {
			times.Fajr = newFajr
		}

		// Calculate new Isha time
		if maghribDuration == 0 {
			ishaDuration := math.Round(ishaPortion * nightDuration)
			newIsha := times.Maghrib.Add(time.Duration(ishaDuration) * time.Minute)
			if times.Isha.IsZero() || newIsha.Before(times.Isha) {
				times.Isha = newIsha
			}
		}

		return times

	default:
		return times
	}
}

func applyTimeCorrections(cfg Config, times Times) Times {
	if !times.Fajr.IsZero() {
		times.Fajr = times.Fajr.Add(cfg.TimeCorrections.Fajr)
	}

	if !times.Sunrise.IsZero() {
		times.Sunrise = times.Sunrise.Add(cfg.TimeCorrections.Sunrise)
	}

	if !times.Zuhr.IsZero() {
		times.Zuhr = times.Zuhr.Add(cfg.TimeCorrections.Zuhr)
	}

	if !times.Asr.IsZero() {
		times.Asr = times.Asr.Add(cfg.TimeCorrections.Asr)
	}

	if !times.Maghrib.IsZero() {
		times.Maghrib = times.Maghrib.Add(cfg.TimeCorrections.Maghrib)
	}

	if !times.Isha.IsZero() {
		times.Isha = times.Isha.Add(cfg.TimeCorrections.Isha)
	}

	return times
}
