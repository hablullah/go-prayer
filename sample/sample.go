package main

import (
	"fmt"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

func main() {
	cfg := prayer.Config{
		Latitude:             -6.21,
		Longitude:            106.85,
		Elevation:            50,
		CalculationMethod:    prayer.Default,
		AsrCalculationMethod: prayer.Shafii,
		PreciseToSeconds:     false,
		Corrections: prayer.TimeCorrections{
			Fajr:    2 * time.Minute,
			Sunrise: -time.Minute,
			Asr:     time.Minute,
			Maghrib: time.Minute,
			Isha:    time.Minute,
		},
	}

	location := time.FixedZone("WIB", 7*3600)
	date := time.Date(2009, 6, 12, 0, 0, 0, 0, location)
	times := prayer.Calculate(date, cfg)

	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", times.Fajr.Format("15:04"))
	fmt.Println("Sunrise =", times.Sunrise.Format("15:04"))
	fmt.Println("Zuhr    =", times.Zuhr.Format("15:04"))
	fmt.Println("Asr     =", times.Asr.Format("15:04"))
	fmt.Println("Maghrib =", times.Maghrib.Format("15:04"))
	fmt.Println("Isha    =", times.Isha.Format("15:04"))
}
