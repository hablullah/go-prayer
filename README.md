# Go-Prayer

[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-prayer?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-prayer)

Go-Prayer is a Go package for calculating prayer/salat times for a specific location at a specific date. As it is right now, this package should be able to give an accurate prayer times for **most** location in the world, except locations with high latitude, where the sun never sets.

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

## Resource

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF](https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn))

## License

Go-Prayer is distributed using [MIT](http://choosealicense.com/licenses/mit/) license.
