package prayer

// AsrConvention is the convention for calculating Asr time.
type AsrConvention int

const (
	// Shafii is the school which said that the Asr time is when the shadow of an object is equals the
	// length of the object plus the length of its shadow when the Sun is at its zenith.
	Shafii AsrConvention = iota

	// Hanafi is the school which said that the Asr time is when the shadow of an object is twice the
	// length of the object plus the length of its shadow when the Sun is at its zenith.
	Hanafi
)
