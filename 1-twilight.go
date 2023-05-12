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

var (
	// AstronomicalTwilight is moment when Sun is 18 degrees below horizon. At this point most stars
	// and other celestial objects still can be seen, however astronomers may be unable to observe
	// some of the fainter stars and galaxies, hence the name of this twilight phase. This is the
	// default twilight convention for this package.
	AstronomicalTwilight = tc(18, 18, 0)

	// MWL is calculation method from Muslim World League with Fajr at 18° and Isha at 17°.
	// Usually used in Europe, Far East and parts of America.
	MWL = tc(18, 17, 0)

	// ISNA is calculation method from Islamic Society of North America with both Fajr and Isha at 15°.
	// Used in North America i.e US and Canada.
	ISNA = tc(15, 15, 0)

	// UmmAlQura is calculation method from Umm al-Qura University in Makkah which used in Saudi Arabia.
	// Fajr at 18.5° and Isha fixed at 90 minutes after Maghrib.
	UmmAlQura = tc(18.5, 18, time.Minute*90)

	// Gulf is calculation method that often used by countries in Gulf region like UAE and Kuwait.
	// Fajr at 19.5° and Isha fixed at 90 minutes after Maghrib.
	Gulf = tc(19.5, 18, time.Minute*90)

	// Algerian is calculation method from Algerian Ministry of Religious Affairs and Wakfs.
	// Fajr at 18° and Isha at 17°.
	Algerian = tc(18, 17, 0)

	// Karachi is calculation method from University of Islamic Sciences, Karachi, with both Fajr and Isha at 18°.
	// Used in Pakistan, Afganistan, Bangladesh and India.
	Karachi = tc(18, 18, 0)

	// Diyanet is calculation method from Turkey's Diyanet İşleri Başkanlığı.
	// It has the same value as MWL with Fajr at 18° and Isha at 17°.
	Diyanet = tc(18, 17, 0)

	// Egypt is calculation method from Egyptian General Authority of Survey with Fajr at 19.5° and Isha at 17.5°.
	// Used in Africa, Syria and Lebanon.
	Egypt = tc(19.5, 17.5, 0)

	// EgyptBis is another version of calculation method from Egyptian General Authority of Survey.
	// Fajr at 20° and Isha at 18°.
	EgyptBis = tc(20, 18, 0)

	// Kemenag is calculation method from Kementerian Agama Republik Indonesia.
	// Fajr at 20° and Isha at 18°.
	Kemenag = tc(20, 18, 0)

	// MUIS is calculation method from Majlis Ugama Islam Singapura.
	// Fajr at 20° and Isha at 18°.
	MUIS = tc(20, 18, 0)

	// JAKIM is calculation method from Jabatan Kemajuan Islam Malaysia.
	// Fajr at 20° and Isha at 18°.
	JAKIM = tc(20, 18, 0)

	// UOIF is calculation method from Union Des Organisations Islamiques De France.
	// Fajr and Isha both at 12°.
	UOIF = tc(12, 12, 0)

	// France15 is calculation method for France region with Fajr and Isha both at 15°.
	France15 = tc(15, 15, 0)

	// France18 is calculation method for France region with Fajr and Isha both at 18°.
	France18 = tc(18, 18, 0)

	// Tunisia is calculation method from Tunisian Ministry of Religious Affairs.
	// Fajr and Isha both at 18°.
	Tunisia = tc(18, 18, 0)

	// Tehran is calculation method from Institute of Geophysics at University of Tehran.
	// Fajr at 17.7° and Isha at 14°.
	Tehran = tc(17.7, 14, 0)

	// Jafari is calculation method from Shia Ithna Ashari that used in some Shia communities worldwide.
	// Fajr at 16° and Isha at 14°.
	Jafari = tc(16, 14, 0)
)

func tc(fajrAngle float64, ishaAngle float64, maghribDuration time.Duration) *TwilightConvention {
	return &TwilightConvention{
		FajrAngle:       fajrAngle,
		IshaAngle:       ishaAngle,
		MaghribDuration: maghribDuration,
	}
}
