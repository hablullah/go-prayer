package prayer

// HighLatConvention is convention for calculating prayer times in area with latitude >45 degrees.
type HighLatConvention int

const (
	// LocalRelativeEstimation is convention that created by cooperation between Fiqh Council
	// of Muslim World League and Islamic Crescents' Observation Project (ICOP). In short, this
	// convention uses average percentage to calculate Fajr and Isha time for abnormal times.
	// This convention used for area above 48.6 degrees. For more detail, check out
	// https://www.astronomycenter.net/latitude.html?l=en
	LocalRelativeEstimation HighLatConvention = iota

	// Mecca is convention based on Fatwa from Dar Al Iftah Al Misrriyah number 2806 dated at
	// 2010-08-08. In this convention, they propose that area with higher latitude when fasting
	// time is too long (more than 18 hours), to follows the fasting time in Mecca. This convention
	// is used for area above 48.5 degrees. See https://www.prayertimes.dk/fatawa.html
	Mecca

	// ShariNormalDay is convention proposed by Mohamed Nabeel Tarabishy, Ph.D. In this convention,
	// they propose that a maximum daylight duration (for fasting) is 17 hours and 36 minutes.
	// If the day is "abnormal" then the fasting times is calculated using the schedule for area
	// with 45 degrees latitude. See https://www.astronomycenter.net/pdf/tarabishyshigh_2014.pdf
	ShariNormalDay

	// NearestDay is convention where the schedule for "abnormal" days will be taken from the
	// schedule of the last "normal" day. In this convention, the day considered "abnormal"
	// when there are no true night. See https://www.islamicity.com/prayertimes/Salat.pdf
	NearestDay

	// NearestLatitude is convention where the schedule for "abnormal" days will be taken from the
	// schedule of location at 48 degrees latitude. In this convention, the day considered
	// "abnormal" when there are no true night. See https://www.islamicity.com/prayertimes/Salat.pdf
	NearestLatitude

	// ForcedNearestLatitude is similar with NearestLatitude, except it will be applied every day
	// and not only on the "abnormal" days.
	ForcedNearestLatitude

	// AngleBased is convention that used by some recent prayer time calculators. Let a be the
	// twilight angle for Isha, and let t = a/60. The period between sunset and sunrise is divided
	// into t parts. Isha begins after the first part. For example, if the twilight angle for Isha
	// is 15, then Isha begins at the end of the first quarter (15/60) of the night. Time for Fajr
	// is calculated similarly. See http://praytimes.org/calculation
	AngleBased

	// OneSeventhNight is convention where the period between sunset and sunrise is divided into
	// seven parts. Isha starts when the first seventh part ends, and Fajr starts when the last
	// seventh part starts. See http://praytimes.org/calculation
	OneSeventhNight

	// MiddleNight is convention where the period from sunset to sunrise is divided into two halves.
	// The first half is considered to be the "night" and the other half as "day break". Fajr and
	// Isha in this method are assumed to be at mid-night during the abnormal periods. See
	// http://praytimes.org/calculation
	MiddleNight
)
