// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/hablullah/go-prayer"
)

func main() {
	// Prepare configuration
	cfg := prayer.Config{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
		TimeCorrections: prayer.TimeCorrections{
			Fajr:    2 * time.Minute,
			Sunrise: -time.Minute,
			Zuhr:    2 * time.Minute,
			Asr:     time.Minute,
			Maghrib: time.Minute,
			Isha:    time.Minute,
		},
	}

	// Specify the date that you want to calculate. This package will use timezone from your date,
	// so make sure to set it to the timezone of the location that you want to calculate.
	zone := time.FixedZone("WIB", 7*3600)
	date := time.Date(2020, 9, 4, 0, 0, 0, 0, zone)

	// Now we just need to calculate it
	result, _ := prayer.Calculate(cfg, date)
	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", result.Fajr.Format("15:04"))
	fmt.Println("Sunrise =", result.Sunrise.Format("15:04"))
	fmt.Println("Zuhr    =", result.Zuhr.Format("15:04"))
	fmt.Println("Asr     =", result.Asr.Format("15:04"))
	fmt.Println("Maghrib =", result.Maghrib.Format("15:04"))
	fmt.Println("Isha    =", result.Isha.Format("15:04"))
}
