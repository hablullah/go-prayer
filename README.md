# Go-Prayer

[![Go Report Card](https://goreportcard.com/badge/github.com/RadhiFadlillah/go-prayer)](https://goreportcard.com/report/github.com/RadhiFadlillah/go-prayer)
[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-prayer?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-prayer)

Go-Prayer is a Go package for calculating prayer/salat times for a specific location at a specific date. As it is right now, this package should be able to give a quite accurate prayer times for **most** location in the world, except locations with high latitude where the sun never sets.

## Usage Examples

For example, we want to get prayer times in Jakarta, Indonesia, at 6 January 2018 :

```go
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
		TimeCorrection: prayer.TimeCorrection{
			prayer.Fajr:    2 * time.Minute,
			prayer.Sunrise: -time.Minute,
			prayer.Zuhr:    2 * time.Minute,
			prayer.Asr:     time.Minute,
			prayer.Maghrib: time.Minute,
			prayer.Isha:    time.Minute,
		},
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
```

Which will give us following results :

```
2009-06-12
Fajr    = 04:38
Sunrise = 05:57
Zuhr    = 11:54
Asr     = 15:15
Maghrib = 17:48
Isha    = 19:01
```

## Accuracy

This package have been compared with various programs and the difference has been found to be within three minutes or so. You should bear that in mind and use the `TimeCorrection` field when needed.

## Resources

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF](https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn))

## License

Go-Prayer is distributed using [MIT](http://choosealicense.com/licenses/mit/) license.
