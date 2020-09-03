package main

import (
	"fmt"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

func main() {
	location := time.FixedZone("WIB", 7*3600)
	date := time.Date(2009, 6, 12, 0, 0, 0, 0, location)

	calc := prayer.Calculator{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
	}

	result := calc.Init().SetDate(date).Calculate()

	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", result[prayer.Fajr].Format("15:04"))
	fmt.Println("Sunrise =", result[prayer.Sunrise].Format("15:04"))
	fmt.Println("Zuhr    =", result[prayer.Zuhr].Format("15:04"))
	fmt.Println("Asr     =", result[prayer.Asr].Format("15:04"))
	fmt.Println("Maghrib =", result[prayer.Maghrib].Format("15:04"))
	fmt.Println("Isha    =", result[prayer.Isha].Format("15:04"))
}
