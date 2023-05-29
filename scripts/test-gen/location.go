package main

import (
	"time"

	"github.com/hablullah/go-prayer"
)

type Location struct {
	Name      string
	Timezone  string
	Latitude  float64
	Longitude float64
}

var testLocations = []Location{
	{ // Tromso (Norway) is representation for location in North Frigid area
		Name:      "Tromso",
		Timezone:  "CET",
		Latitude:  69.682778,
		Longitude: 18.942778,
	}, { // London (UK) is representation for location in North Temperate area
		Name:      "London",
		Timezone:  "Europe/London",
		Latitude:  51.507222,
		Longitude: -0.1275,
	}, { // Jakarta (Indonesia) is representation for location in Torrid area
		Name:      "Jakarta",
		Timezone:  "Asia/Jakarta",
		Latitude:  -6.175,
		Longitude: 106.825,
	}, { // Wellington (New Zealand) is representation for location in South Temperate area
		Name:      "Wellington",
		Timezone:  "Pacific/Auckland",
		Latitude:  -41.288889,
		Longitude: 174.777222,
	},
}

func getSchedules(loc Location) []prayer.Schedule {
	tz, _ := time.LoadLocation(loc.Timezone)
	schedules, _ := prayer.Calculate(prayer.Config{
		Latitude:            loc.Latitude,
		Longitude:           loc.Longitude,
		Timezone:            tz,
		TwilightConvention:  prayer.AstronomicalTwilight(),
		AsrConvention:       prayer.Shafii,
		HighLatitudeAdapter: prayer.NearestLatitude(),
		PreciseToSeconds:    true,
	}, 2023)
	return schedules
}
