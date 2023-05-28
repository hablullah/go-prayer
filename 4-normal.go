package prayer

import (
	"math"
	"time"

	"github.com/hablullah/go-sampa"
)

func calcNormal(cfg Config, year int) ([]Schedule, int) {
	// Prepare location
	location := sampa.Location{
		Latitude:  cfg.Latitude,
		Longitude: cfg.Longitude,
		Elevation: cfg.Elevation,
	}

	// Prepare custom Sun events
	customEvents := []sampa.CustomSunEvent{{
		Name:          "dawn",
		BeforeTransit: true,
		Elevation:     func(sampa.SunPosition) float64 { return -18 },
	}, {
		Name:          "dusk",
		BeforeTransit: false,
		Elevation:     func(sampa.SunPosition) float64 { return -18 },
	}, {
		Name:          "fajr",
		BeforeTransit: true,
		Elevation:     func(sampa.SunPosition) float64 { return -cfg.TwilightConvention.FajrAngle },
	}, {
		Name:          "isha",
		BeforeTransit: false,
		Elevation:     func(sampa.SunPosition) float64 { return -cfg.TwilightConvention.IshaAngle },
	}, {
		Name:          "asr",
		BeforeTransit: false,
		Elevation: func(todayData sampa.SunPosition) float64 {
			a := getAsrCoefficient(cfg.AsrConvention)
			b := math.Abs(todayData.TopocentricDeclination - cfg.Latitude)
			elevation := acot(a + math.Tan(degToRad(b)))
			return radToDeg(elevation)
		},
	}}

	// Calculate schedules for each day in a year.
	start := time.Date(year, 1, 1, 0, 0, 0, 0, cfg.Timezone)
	limit := start.AddDate(1, 0, 0)
	nDays := int(limit.Sub(start).Hours() / 24)

	// Create slice to contain result
	schedules := make([]Schedule, nDays)

	// Calculate each day
	var idx int
	var nAbnormal int
	for dt := start; dt.Before(limit); dt = dt.AddDate(0, 0, 1) {
		// Calculate the events
		e, _ := sampa.GetSunEvents(dt, location, nil, customEvents...)

		// Create the prayer schedule
		s := Schedule{
			Date:    dt.Format("2006-01-02"),
			Fajr:    e.Others["fajr"].DateTime,
			Sunrise: e.Sunrise.DateTime,
			Zuhr:    e.Transit.DateTime,
			Asr:     e.Others["asr"].DateTime,
			Maghrib: e.Sunset.DateTime,
			Isha:    e.Others["isha"].DateTime,
		}

		// Check if schedule is normal
		dawn := e.Others["dawn"].DateTime
		dusk := e.Others["dusk"].DateTime
		hasNight := !e.Sunrise.IsZero() && !e.Sunset.IsZero()
		hasTwilight := !dawn.IsZero() && !dusk.IsZero()
		s.IsNormal = hasNight && hasTwilight

		// Save the schedule
		schedules[idx] = s
		if !s.IsNormal {
			nAbnormal++
		}
		idx++
	}

	return schedules, nAbnormal
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func acot(cotValue float64) float64 {
	return math.Atan(1 / cotValue)
}
