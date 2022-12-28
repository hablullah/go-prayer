package prayer

// HighLatConvention is convention for calculating prayer times in area with latitude >45 degrees.
// Check out https://www.prayertimes.dk/story.html for more detail.
type HighLatConvention int

const (
	// Disabled will not calculate the adjustment for higher latitude and instead will return the
	// schedule as it is. For area in high or extreme latitude, it might return zero for Fajr,
	// Sunrise, Maghrib and Isha
	Disabled HighLatConvention = iota - 1

	// Mecca is convention based on Fatwa from Dar Al Iftah Al Misrriyah number 2806 dated at
	// 2010-08-08. In this convention, they propose that area with higher latitude to follows the
	// schedule in Mecca when abnormal days occured. In this convention, the day is considered
	// "abnormal" when there are no true night, the fasting time is more than 18 hours, or the
	// fasting time is less than 4 hours. See https://www.prayertimes.dk/fatawa.html
	Mecca

	// ForceMecca is similar with Mecca, except it will be applied every day and not only on
	// the "abnormal" days.
	ForceMecca

	// LocalRelativeEstimation is convention that created by cooperation between Fiqh Council
	// of Muslim World League and Islamic Crescents' Observation Project (ICOP). In short, this
	// convention uses average percentage to calculate Fajr and Isha time for abnormal times.
	// This convention used for area above 48.6 degrees. For more detail, check out
	// https://www.astronomycenter.net/latitude.html?l=en
	LocalRelativeEstimation

	// NearestDay is convention where the schedule for "abnormal" days will be taken from the
	// schedule of the last "normal" day. In this convention, the day considered "abnormal"
	// when there are no true night. See https://www.islamicity.com/prayertimes/Salat.pdf
	NearestDay

	// NearestLatitude is convention where the schedule for "abnormal" days will be taken from the
	// schedule of location at 48 degrees latitude. In this convention, the day considered
	// "abnormal" when there are no true night. See https://www.islamicity.com/prayertimes/Salat.pdf
	NearestLatitude

	// ForceNearestLatitude is similar with NearestLatitude, except it will be applied every day
	// and not only on the "abnormal" days.
	ForceNearestLatitude

	// ShariNormalDay is convention proposed by Mohamed Nabeel Tarabishy, Ph.D. In this convention,
	// they propose that a normal day is when the fasting period is between 10h17m and 17h36m. If
	// the day is "abnormal" then the fasting times is calculated using the schedule for area
	// with 45 degrees latitude. See https://www.astronomycenter.net/pdf/tarabishyshigh_2014.pdf
	ShariNormalDay

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
