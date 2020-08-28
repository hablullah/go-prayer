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

	calc.Init().SetDate(date)
	fajr := calc.Calculate(prayer.Fajr)
	sunrise := calc.Calculate(prayer.Sunrise)
	zuhr := calc.Calculate(prayer.Zuhr)
	asr := calc.Calculate(prayer.Asr)
	maghrib := calc.Calculate(prayer.Maghrib)
	isha := calc.Calculate(prayer.Isha)

	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", fajr.Format("15:04"))
	fmt.Println("Sunrise =", sunrise.Format("15:04"))
	fmt.Println("Zuhr    =", zuhr.Format("15:04"))
	fmt.Println("Asr     =", asr.Format("15:04"))
	fmt.Println("Maghrib =", maghrib.Format("15:04"))
	fmt.Println("Isha    =", isha.Format("15:04"))
}
