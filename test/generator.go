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
	"ottawa":     -5,
	"cairo":      2,
	"sana":       3,
	"singapore":  8,
	"brasilia":   -3,
	"maputo":     2,
	"canberra":   11,
	"wellington": 13,
}

var calculators = map[string]*prayer.Calculator{
	"ottawa": (&prayer.Calculator{
		Latitude:          45.424722,
		Longitude:         -75.695,
		Elevation:         76,
		CalculationMethod: prayer.ISNA,
		PreciseToSeconds:  true,
	}).Init(),
	"cairo": (&prayer.Calculator{
		Latitude:          30.033333,
		Longitude:         31.233333,
		Elevation:         22,
		CalculationMethod: prayer.Egypt,
		PreciseToSeconds:  true,
	}).Init(),
	"sana": (&prayer.Calculator{
		Latitude:          15.348333,
		Longitude:         44.206389,
		Elevation:         2266,
		CalculationMethod: prayer.UmmAlQura,
		PreciseToSeconds:  true,
	}).Init(),
	"singapore": (&prayer.Calculator{
		Latitude:          1.283333,
		Longitude:         103.833333,
		Elevation:         93,
		CalculationMethod: prayer.MUIS,
		PreciseToSeconds:  true,
	}).Init(),
	"brasilia": (&prayer.Calculator{
		Latitude:          -15.793889,
		Longitude:         -47.882778,
		Elevation:         1091,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	}).Init(),
	"maputo": (&prayer.Calculator{
		Latitude:          -25.966667,
		Longitude:         32.583333,
		Elevation:         20,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	}).Init(),
	"canberra": (&prayer.Calculator{
		Latitude:          -35.293056,
		Longitude:         149.126944,
		Elevation:         577,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	}).Init(),
	"wellington": (&prayer.Calculator{
		Latitude:          -41.288889,
		Longitude:         174.777222,
		Elevation:         13,
		CalculationMethod: prayer.MWL,
		PreciseToSeconds:  true,
	}).Init(),
}

func main() {
	for city, calc := range calculators {
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
			result := calc.SetDate(date).Calculate()
			err = writer.Write([]string{
				date.Format("2006-01-02"),
				result[prayer.Fajr].Format("15:04:05"),
				result[prayer.Sunrise].Format("15:04:05"),
				result[prayer.Zuhr].Format("15:04:05"),
				result[prayer.Asr].Format("15:04:05"),
				result[prayer.Maghrib].Format("15:04:05"),
				result[prayer.Isha].Format("15:04:05"),
			})

			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
