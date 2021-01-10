package prayer

// CalculationMethod is the conventions for calculating prayer times, especially Fajr and Isha.
type CalculationMethod int

// AsrConvention is the convention for calculating Asr time.
type AsrConvention int

// HighLatitudeMethod is the method for calculating prayer time in area with higher latitude
// (more than 45N or 45S).
type HighLatitudeMethod int

const (
	// MWL is calculation method from Muslim World League with Fajr at 18° and Isha at 17°.
	// Usually used in Europe, Far East and parts of US.
	MWL CalculationMethod = iota

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
	// Shafii is the school which said that the Asr time is when the shadow of an object is equals the
	// length of the object plus the length of its shadow when the sun is at its zenith.
	Shafii AsrConvention = iota

	// Hanafi is the school which said that the Asr time is when the shadow of an object is twice the
	// length of the object plus the length of its shadow when the sun is at its zenith.
	Hanafi
)

const (
	// AngleBased is method that created after conference between Muslim scholars in Brussels, 25-26 May
	// 2009. It's one of the most common method and used by some recent prayer time calculators. Let a
	// be the twilight angle for Isha, and let t = a/60. The period between sunset and sunrise is divided
	// into t parts. Isha begins after the first part. For example, if the twilight angle for Isha is 15,
	// then Isha begins at the end of the first quarter (15/60) of the night. Time for Fajr is calculated
	// similarly. This method is only used in area between 48.6 and 66.6 latitude.
	AngleBased HighLatitudeMethod = iota

	// OneSevenNight is method where the period between sunset and sunrise is divided into seven parts.
	// Isha starts when the first seventh part ends, and Fajr starts when the last seventh part starts.
	// Like angle-based, this method is only used in area between 48.6 and 66.6 latitude.
	OneSeventhNight

	// MiddleNight is method where the period from sunset to sunrise is divided into two halves. The first
	// half is considered to be the "night" and the other half as "day break". Fajr and Isha in this method
	// are assumed to be at mid-night during the abnormal periods. Like angle-based, this method is only
	// used in area between 48.6 and 66.6 latitude.
	MiddleNight

	// NormalRegion is the latest method (from 2014), created by Mohamed Nabeel Tarabishy, Ph.D. In this
	// method, if in a day the fasting time (duration between Fajr and Sunset) is too short (less than 10h
	// 17m) or too long (more than 17h 36m), then that day is considered abnormal. In those abnormal days,
	// the prayer times is calculated by setting the latitude into 45. With that said, this method only
	// used in area above 45 latitude.
	NormalRegion

	// ForcedNormalRegion is similar with NormalRegion method. The difference is in this method, all area
	// above 45 latitude will be calculated as if it's located at latitude 45. This is because when using
	// NormalRegion method, there will be sudden changes in the length of the day of fasting so for some
	// cases it might be better to set the latitude permanently into 45 degrees.
	ForcedNormalRegion
)
