package prayer_test

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

// These data must be identical with the one in test/generator.go
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

func Test_Calculator(t *testing.T) {
	for city, calc := range calculators {
		testData, err := openTestData(city)
		if err != nil {
			t.Fatal(err)
		}

		for _, data := range testData {
			result := calc.SetDate(data.Date).Calculate()
			diff := data.Fajr.Sub(result[prayer.Fajr]).Seconds()
			if math.Abs(diff) > 0 {
				t.Errorf("%s, %s, Fajr: want %s got %s (%v)\n",
					strings.ToUpper(city),
					data.Date.Format("2006-01-02"),
					data.Fajr.Format("15:04:05"),
					result[prayer.Fajr].Format("15:04:05"),
					diff)
			}
		}
	}
}

type TestData struct {
	Date    time.Time
	Fajr    time.Time
	Sunrise time.Time
	Zuhr    time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
}

func openTestData(city string) ([]TestData, error) {
	// Open test file
	f, err := os.Open("test/" + city + ".csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Prepare timezone
	timezone := time.FixedZone("", timezones[city]*3600)

	// Parse test file
	dataList := []TestData{}
	csvReader := csv.NewReader(f)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		date, err := time.ParseInLocation("2006-01-02", record[0], timezone)
		if err != nil {
			return nil, err
		}

		fajr, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[1], timezone)
		if err != nil {
			return nil, err
		}

		sunrise, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[2], timezone)
		if err != nil {
			return nil, err
		}

		zuhr, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[3], timezone)
		if err != nil {
			return nil, err
		}

		asr, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[4], timezone)
		if err != nil {
			return nil, err
		}

		maghrib, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[5], timezone)
		if err != nil {
			return nil, err
		}

		isha, err := time.ParseInLocation("2006-01-02 15:04:05", record[0]+" "+record[6], timezone)
		if err != nil {
			return nil, err
		}

		dataList = append(dataList, TestData{
			Date:    date,
			Fajr:    fajr,
			Sunrise: sunrise,
			Zuhr:    zuhr,
			Asr:     asr,
			Maghrib: maghrib,
			Isha:    isha,
		})
	}

	return dataList, nil
}
