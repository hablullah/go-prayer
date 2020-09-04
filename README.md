# Go-Prayer

[![Go Report Card](https://goreportcard.com/badge/github.com/RadhiFadlillah/go-prayer)](https://goreportcard.com/report/github.com/RadhiFadlillah/go-prayer)
[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-prayer?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-prayer)

Go-Prayer is a Go package for calculating prayer/salat times for a specific location at a specific date. As it is right now, this package should be able to give a quite accurate prayer times for **most** location in the world, except locations with latitude higher than 65 N or lower than 65 S.

## Usage Examples

For example, we want to get prayer times in Jakarta, Indonesia, at 4 September 2020 :

```go
package main

import (
	"fmt"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

func main() {
	// Prepare calculator
	calc := prayer.Calculator{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
		TimeCorrection: prayer.TimeCorrection{
			prayer.Fajr:    2 * time.Minute,
			prayer.Sunrise: -time.Minute,
			prayer.Zuhr:    2 * time.Minute,
			prayer.Asr:     time.Minute,
			prayer.Maghrib: time.Minute,
			prayer.Isha:    time.Minute,
		},
	}

	// Initiate the calculator. It must be run every time any
	// parameter in calculator above changed.
	calc.Init()

	// Specify the date that you want to calculate. This package
	// will use timezone from your date, so make sure to set it
	// to the timezone of the location that you want to calculate.
	zone := time.FixedZone("WIB", 7*3600)
	date := time.Date(2020, 9, 4, 0, 0, 0, 0, zone)
	calc.SetDate(date)

	// Now we just need to calculate
	result := calc.Calculate()
	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", result[prayer.Fajr].Format("15:04"))
	fmt.Println("Sunrise =", result[prayer.Sunrise].Format("15:04"))
	fmt.Println("Zuhr    =", result[prayer.Zuhr].Format("15:04"))
	fmt.Println("Asr     =", result[prayer.Asr].Format("15:04"))
	fmt.Println("Maghrib =", result[prayer.Maghrib].Format("15:04"))
	fmt.Println("Isha    =", result[prayer.Isha].Format("15:04"))
}
```

Which will give us following results :

```
2020-09-04
Fajr    = 04:36
Sunrise = 05:49
Zuhr    = 11:54
Asr     = 15:10
Maghrib = 17:54
Isha    = 19:02
```

By the way, method [`Init`](https://godoc.org/github.com/RadhiFadlillah/go-prayer#Calculator.Init) and [`SetDate`](https://godoc.org/github.com/RadhiFadlillah/go-prayer#Calculator.Init) is chainable, so you can shorten above code into this :

```go
package main

import (
	"fmt"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

func main() {
	// Specify the date and timezone.
	zone := time.FixedZone("WIB", 7*3600)
	date := time.Date(2020, 9, 4, 0, 0, 0, 0, zone)

	// Do it in one run
	result := (&prayer.Calculator{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
		TimeCorrection: prayer.TimeCorrection{
			prayer.Fajr:    2 * time.Minute,
			prayer.Sunrise: -time.Minute,
			prayer.Zuhr:    2 * time.Minute,
			prayer.Asr:     time.Minute,
			prayer.Maghrib: time.Minute,
			prayer.Isha:    time.Minute,
		},
	}).Init().SetDate(date).Calculate()

	// Print result
}
```

## Accuracy

This package have been compared with various programs and the difference has been found to be within three minutes or so. You should bear that in mind and use the `TimeCorrection` field when needed.

## Resources

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF](https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn))

## License

Go-Prayer is distributed using [MIT](http://choosealicense.com/licenses/mit/) license.
