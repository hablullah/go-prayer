package prayer

import "time"

// TwilightConvention is the convention that specifies time for Fajr (dawn) and Isha (dusk). Most of
// the conventions use Solar angle elevation for both dawn and dusk time, however there are several
// convention where dusk times depends on sunset (Maghrib) times.
type TwilightConvention struct {
	FajrAngle       float64
	IshaAngle       float64
	MaghribDuration time.Duration
}

// AstronomicalTwilight is moment when Sun is 18 degrees below horizon. At this point most stars
// and other celestial objects still can be seen, however astronomers may be unable to observe
// some of the fainter stars and galaxies, hence the name of this twilight phase. This is the
// default twilight convention for this package.
func AstronomicalTwilight() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 18}
}

// MWL is calculation method from Muslim World League with Fajr at 18° and Isha at 17°.
// Usually used in Europe, Far East and parts of America.
func MWL() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 17}
}

// ISNA is calculation method from Islamic Society of North America with both Fajr and Isha at 15°.
// Used in North America i.e US and Canada.
func ISNA() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 15, IshaAngle: 15}
}

// UmmAlQura is calculation method from Umm al-Qura University in Makkah which used in Saudi Arabia.
// Fajr at 18.5° and Isha fixed at 90 minutes after Maghrib.
func UmmAlQura() *TwilightConvention {
	return &TwilightConvention{
		FajrAngle:       18.5,
		IshaAngle:       18.5,
		MaghribDuration: 90 * time.Minute}
}

// Gulf is calculation method that often used by countries in Gulf region like UAE and Kuwait.
// Fajr at 19.5° and Isha fixed at 90 minutes after Maghrib.
func Gulf() *TwilightConvention {
	return &TwilightConvention{
		FajrAngle:       19.5,
		IshaAngle:       19.5,
		MaghribDuration: 90 * time.Minute}
}

// Algerian is calculation method from Algerian Ministry of Religious Affairs and Wakfs.
// Fajr at 18° and Isha at 17°.
func Algerian() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 17}
}

// Karachi is calculation method from University of Islamic Sciences, Karachi, with both Fajr and
// Isha at 18°. Used in Pakistan, Afganistan, Bangladesh and India.
func Karachi() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 18}
}

// Diyanet is calculation method from Turkey's Diyanet İşleri Başkanlığı.
// It has the same value as MWL with Fajr at 18° and Isha at 17°.
func Diyanet() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 17}
}

// Egypt is calculation method from Egyptian General Authority of Survey with Fajr at 19.5° and
// Isha at 17.5°. Used in Africa, Syria and Lebanon.
func Egypt() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 19.5, IshaAngle: 17.5}
}

// EgyptBis is another version of calculation method from Egyptian General Authority of Survey.
// Fajr at 20° and Isha at 18°.
func EgyptBis() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 20, IshaAngle: 18}
}

// Kemenag is calculation method from Kementerian Agama Republik Indonesia.
// Fajr at 20° and Isha at 18°.
func Kemenag() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 20, IshaAngle: 18}
}

// MUIS is calculation method from Majlis Ugama Islam Singapura.
// Fajr at 20° and Isha at 18°.
func MUIS() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 20, IshaAngle: 18}
}

// JAKIM is calculation method from Jabatan Kemajuan Islam Malaysia.
// Fajr at 20° and Isha at 18°.
func JAKIM() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 20, IshaAngle: 18}
}

// UOIF is calculation method from Union Des Organisations Islamiques De France.
// Fajr and Isha both at 12°.
func UOIF() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 12, IshaAngle: 12}
}

// France15 is calculation method for France region with Fajr and Isha both at 15°.
func France15() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 15, IshaAngle: 15}
}

// France18 is calculation method for France region with Fajr and Isha both at 18°.
func France18() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 18}
}

// Tunisia is calculation method from Tunisian Ministry of Religious Affairs.
// Fajr and Isha both at 18°.
func Tunisia() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 18, IshaAngle: 18}
}

// Tehran is calculation method from Institute of Geophysics at University of Tehran.
// Fajr at 17.7° and Isha at 14°.
func Tehran() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 17.7, IshaAngle: 14}
}

// Jafari is calculation method from Shia Ithna Ashari that used in some Shia communities worldwide.
// Fajr at 16° and Isha at 14°.
func Jafari() *TwilightConvention {
	return &TwilightConvention{FajrAngle: 16, IshaAngle: 14}
}
