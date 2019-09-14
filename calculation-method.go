package prayer

import (
	"time"

	"github.com/shopspring/decimal"
)

// CalculationMethod is the conventions for calculating prayer times, especially Fajr and Isha.
// For references, check these website :
// - http://praytimes.org/wiki/Calculation_Methods
// - http://www.islamicfinder.us/index.php/api/index
// - https://www.muslimpro.com/en/prayer-times
type CalculationMethod int

// AsrJuristicMethod is the conventions for calculating Asr time.
// For details, check http://www.prayerminder.com/faq.php#Fiqh.
type AsrJuristicMethod int

const (
	// Default is the default calculation method with the same value as MWL.
	Default CalculationMethod = iota

	// MWL is calculation method from Muslim World League with Fajr at 18° and Isha at 17°.
	// Usually used in Europe, Far East and parts of US.
	MWL

	// ISNA is calculation method from Islamic Society of North America with both Fajr and Isha at 15°.
	// Used in North America i.e US and Canada.
	ISNA

	// UmmAlQura is calculation method from Umm al-Qura University in Makkah which used in Saudi Arabia.
	// Fajr at 18.5° and Isha fixed at 90 minutes after Maghrib.
	UmmAlQura

	// Gulf is calculation method that often used by countries in Gulf region like UAE and Kuwait.
	// Fajr at 19.5° and Isha fixed at 90 minutes after Maghrib.
	Gulf

	// Algerian is calculation method from Algerian Ministry of Religious Affairs and Wakfs.
	// Fajr at 18° and Isha at 17°.
	Algerian

	// Karachi is calculation method from University of Islamic Sciences, Karachi, with both Fajr and Isha at 18°.
	// Used in Pakistan, Afganistan, Bangladesh and India.
	Karachi

	// Diyanet is calculation method from Turkey's Diyanet İşleri Başkanlığı.
	// It has the same value as MWL with Fajr at 18° and Isha at 17°.
	Diyanet

	// Egypt is calculation method from Egyptian General Authority of Survey with Fajr at 19.5° and Isha at 17.5°.
	// Used in Africa, Syria and Lebanon.
	Egypt

	// EgyptBis is another version of calculation method from Egyptian General Authority of Survey.
	// Fajr at 20° and Isha at 18°.
	EgyptBis

	// Kemenag is calculation method from Kementerian Agama Republik Indonesia.
	// Fajr at 20° and Isha at 18°.
	Kemenag

	// MUIS is calculation method from Majlis Ugama Islam Singapura.
	// Fajr at 20° and Isha at 18°.
	MUIS

	// JAKIM is calculation method from Jabatan Kemajuan Islam Malaysia.
	// Fajr at 20° and Isha at 18°.
	JAKIM

	// UOIF is calculation method from Union Des Organisations Islamiques De France.
	// Fajr and Isha both at 12°.
	UOIF

	// France15 is calculation method for France region with Fajr and Isha both at 15°.
	France15

	// France18 is calculation method for France region with Fajr and Isha both at 18°.
	France18

	// Tunisia is calculation method from Tunisian Ministry of Religious Affairs.
	// Fajr and Isha both at 18°.
	Tunisia

	// Tehran is calculation method from Institute of Geophysics at University of Tehran.
	// Fajr at 17.7° and Isha at 14°.
	Tehran

	// Jafari is calcuation method from Shia Ithna Ashari that used in some Shia communities worldwide.
	// Fajr at 16° and Isha at 14°.
	Jafari
)

const (
	// Hanafi is the school which said that the Asr time is when the shadow of
	// an object is twice the length of the object plus the length of its shadow
	// when the sun is at its zenith.
	Hanafi AsrJuristicMethod = iota

	// Shafii is the school which said that the Asr time is when the shadow of
	// an object is equals the length of the object plus the length of its shadow
	// when the sun is at its zenith.
	Shafii
)

func getCalculationAngle(cfg Config) (decimal.Decimal, decimal.Decimal, time.Duration) {
	var maghribDuration time.Duration
	var fajrAngle, ishaAngle float64

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

	return decimal.NewFromFloat(fajrAngle),
		decimal.NewFromFloat(ishaAngle),
		maghribDuration
}

func getAsrCoefficient(cfg Config) decimal.Decimal {
	switch cfg.AsrJuristicMethod {
	case Hanafi:
		return decimal.New(2, 0)
	default:
		return decimal.New(1, 0)
	}
}
