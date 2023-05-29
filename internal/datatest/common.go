package datatest

import (
	"regexp"
	"strconv"
	"time"

	"github.com/hablullah/go-prayer"
)

var rxDT = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})`)

type TestData struct {
	Name      string
	Latitude  float64
	Longitude float64
	Timezone  *time.Location
	Schedules []prayer.Schedule
}

func schedule(date string, fajr, sunrise, zuhr, asr, maghrib, isha string, tz *time.Location) prayer.Schedule {
	return prayer.Schedule{
		Date:    date,
		Fajr:    st(fajr, tz),
		Sunrise: st(sunrise, tz),
		Zuhr:    st(zuhr, tz),
		Asr:     st(asr, tz),
		Maghrib: st(maghrib, tz),
		Isha:    st(isha, tz),
	}
}

func st(s string, tz *time.Location) time.Time {
	parts := rxDT.FindStringSubmatch(s)
	if len(parts) != 7 {
		return time.Time{}
	}

	year, _ := strconv.Atoi(parts[1])
	month, _ := strconv.Atoi(parts[2])
	day, _ := strconv.Atoi(parts[3])
	hour, _ := strconv.Atoi(parts[4])
	minute, _ := strconv.Atoi(parts[5])
	second, _ := strconv.Atoi(parts[6])
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, tz)
}
