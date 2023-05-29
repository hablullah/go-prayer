package main

import (
	"fmt"
	"time"

	"github.com/hablullah/go-prayer"
)

func main() {
	// Calculate prayer schedule in Jakarta for 2023
	asiaJakarta, _ := time.LoadLocation("Asia/Jakarta")
	schedules, _ := prayer.Calculate(prayer.Config{
		Latitude:            -6.14,
		Longitude:           106.81,
		Timezone:            asiaJakarta,
		TwilightConvention:  prayer.Kemenag(),
		AsrConvention:       prayer.Shafii,
		HighLatitudeAdapter: prayer.NearestLatitude(),
		PreciseToSeconds:    false,
	}, 2023)

	// Print the schedules
	for _, s := range schedules {
		fmt.Println(
			"'"+s.Date+"'",
			s.Fajr.Format("'2006-01-02 15:04:05'"),
			s.Sunrise.Format("'2006-01-02 15:04:05'"),
			s.Zuhr.Format("'2006-01-02 15:04:05'"),
			s.Asr.Format("20'06-01-02 15:04:05'"),
			s.Maghrib.Format("'2006-01-02 15:04:05'"),
			s.Isha.Format("'2006-01-02 15:04:05'"),
		)
	}
}
