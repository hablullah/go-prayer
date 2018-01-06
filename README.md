# Go-Prayer

[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-prayer?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-prayer)

Go-Prayer is a Go package for calculating prayer/salat times for a specific location and a specific date.

## Usage Examples

For example, we want to get prayer times in Jakarta, Indonesia, at 6 January 2018 :

```go
package main

import (
	"fmt"
	"github.com/RadhiFadlillah/go-prayer"
	"time"
)

func main() {
	calculator := prayer.Calculator{
		Latitude:             -6.1751,
		Longitude:            106.865,
		Elevation:            7.9,
		CalculationMethod:    prayer.Default,
		AsrCalculationMethod: prayer.Shafii,
	}

	location := time.FixedZone("WIB", 7*3600)
	date := time.Date(2018, 1, 6, 0, 0, 0, 0, location)
	adhan, _ := calculator.Calculate(date)

	fmt.Println("Fajr:", adhan[prayer.Fajr].Format("15:04:05"))
	fmt.Println("Sunrise:", adhan[prayer.Sunrise].Format("15:04:05"))
	fmt.Println("Zuhr:", adhan[prayer.Zuhr].Format("15:04:05"))
	fmt.Println("Asr:", adhan[prayer.Asr].Format("15:04:05"))
	fmt.Println("Maghrib:", adhan[prayer.Maghrib].Format("15:04:05"))
	fmt.Println("Isha:", adhan[prayer.Isha].Format("15:04:05"))
}
```

Which will give us following results :

```
Fajr: 04:20:00
Sunrise: 05:43:00
Zuhr: 11:59:00
Asr: 15:25:00
Maghrib: 18:13:00
Isha: 19:28:00
```

## Resource

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF](https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn))

## License

Go-Prayer is distributed using [MIT](http://choosealicense.com/licenses/mit/) license.
