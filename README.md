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
	cfg := prayer.Config{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrJuristicMethod: prayer.Shafii,
		PreciseToSeconds:  false,
		Corrections: prayer.TimeCorrections{
			Fajr:    2 * time.Minute,
			Sunrise: -time.Minute,
			Zuhr:    2 * time.Minute,
			Asr:     time.Minute,
			Maghrib: time.Minute,
			Isha:    time.Minute,
		},
	}

	location := time.FixedZone("WIB", 7*3600)
	date := time.Date(2009, 6, 12, 0, 0, 0, 0, location)
	adhan, _ := prayer.GetTimes(date, cfg)

	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", adhan.Fajr.Format("15:04"))
	fmt.Println("Sunrise =", adhan.Sunrise.Format("15:04"))
	fmt.Println("Zuhr    =", adhan.Zuhr.Format("15:04"))
	fmt.Println("Asr     =", adhan.Asr.Format("15:04"))
	fmt.Println("Maghrib =", adhan.Maghrib.Format("15:04"))
	fmt.Println("Isha    =", adhan.Isha.Format("15:04"))
}
```

Which will give us following results :

```
2009-06-12
Fajr    = 04:38
Sunrise = 05:58
Zuhr    = 11:55
Asr     = 15:16
Maghrib = 17:48
Isha    = 19:02
```

## Accuracy

This package have been compared with various programs and the difference has been found to be within three minutes or so. You should bear that in mind and use the `Corrections` field when needed.

## Resources

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF](https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn))

## License

Go-Prayer is distributed using [MIT](http://choosealicense.com/licenses/mit/) license.
