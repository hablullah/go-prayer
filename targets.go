package prayer

// Target is the prayer times or special times to calculate.
type Target int

const (
	// Fajr is the time when the sky begins to lighten (dawn).
	Fajr Target = iota

	// Sunrise is the time at which the first part of the Sun appears above the horizon.
	Sunrise

	// Zuhr is the time when the Sun begins to decline after reaching its highest point in the sky.
	Zuhr

	// Asr is the time when the length of any object's shadow reaches a factor (usually 1 or 2)
	// of the length of the object itself plus the length of that object's shadow at noon.
	Asr

	// Maghrib is the time a little after sunset
	Maghrib

	// Isha is the time at which darkness falls and there is no scattered light in the sky.
	Isha
)
