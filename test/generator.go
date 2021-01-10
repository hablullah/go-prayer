// +build ignore

package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

var timezones = map[string]int{
	"longyearbyen": 1,
	"oslo":         1,
	"ottawa":       -5,
	"cairo":        2,
	"sana":         3,
	"singapore":    8,
	"brasilia":     -3,
	"maputo":       2,
	"canberra":     11,
	"wellington":   13,
	"king-edward":  -2,
}

var configs = map[string]prayer.Config{
	"longyearbyen": {
		Latitude:           78.216667,
		Longitude:          15.633333,
		Elevation:          20,
		CalculationMethod:  prayer.MWL,
		PreciseToSeconds:   true,
		HighLatitudeMethod: prayer.NormalRegion,
	},
	"oslo": {
		Latitude:          59.913889,
		Longitude:         10.752222,
		Elevation:         11,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
	"ottawa": {
		Latitude:          45.424722,
		Longitude:         -75.695,
		Elevation:         76,
		CalculationMethod: prayer.ISNA,
		PreciseToSeconds:  true,
	},
	"cairo": {
		Latitude:          30.033333,
		Longitude:         31.233333,
		Elevation:         22,
		CalculationMethod: prayer.Egypt,
		PreciseToSeconds:  true,
	},
	"sana": {
		Latitude:          15.348333,
		Longitude:         44.206389,
		Elevation:         2266,
		CalculationMethod: prayer.UmmAlQura,
		PreciseToSeconds:  true,
	},
	"singapore": {
		Latitude:          1.283333,
		Longitude:         103.833333,
		Elevation:         93,
		CalculationMethod: prayer.MUIS,
		PreciseToSeconds:  true,
	},
	"brasilia": {
		Latitude:          -15.793889,
		Longitude:         -47.882778,
		Elevation:         1091,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
	"maputo": {
		Latitude:          -25.966667,
		Longitude:         32.583333,
		Elevation:         20,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
	"canberra": {
		Latitude:          -35.293056,
		Longitude:         149.126944,
		Elevation:         577,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
	"wellington": {
		Latitude:          -41.288889,
		Longitude:         174.777222,
		Elevation:         13,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
	"king-edward": {
		Latitude:          -54.283333,
		Longitude:         -36.5,
		Elevation:         3,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	},
}

func main() {
	for city, cfg := range configs {
		log.Println("generating test for", city)

		f, err := os.Create("test/" + city + ".csv")
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		zone := time.FixedZone("ZONE", timezones[city]*3600)
		for date := time.Date(2021, 1, 1, 0, 0, 0, 0, zone); date.Year() == 2021; date = date.AddDate(0, 0, 1) {
			result, _ := prayer.Calculate(cfg, date)
			err = writer.Write([]string{
				date.Format("2006-01-02"),
				result.Fajr.Format("15:04:05"),
				result.Sunrise.Format("15:04:05"),
				result.Zuhr.Format("15:04:05"),
				result.Asr.Format("15:04:05"),
				result.Maghrib.Format("15:04:05"),
				result.Isha.Format("15:04:05"),
			})

			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
