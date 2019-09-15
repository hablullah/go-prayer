package prayer

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	fp "path/filepath"
	"testing"
	"time"
)

func TestGetTimes(t *testing.T) {
	// All geography data taken from https://dateandtime.info
	// while test file taken from https://muslimpro.com
	scenarios := []struct {
		name              string
		latitude          float64
		longitude         float64
		elevation         float64
		calculationMethod CalculationMethod
		asrJuristicMethod AsrJuristicMethod
		timezone          *time.Location
		testFile          string
	}{{
		name:              "San Fransisco, ISNA",
		latitude:          37.7749300,
		longitude:         -122.4194200,
		elevation:         28,
		calculationMethod: ISNA,
		asrJuristicMethod: Shafii,
		timezone:          time.FixedZone("PDT", -7*3600),
		testFile:          "san-fransisco-isna.csv",
	}}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			// Create config
			cfg := Config{
				Latitude:          tt.latitude,
				Longitude:         tt.longitude,
				Elevation:         tt.elevation,
				CalculationMethod: tt.calculationMethod,
				AsrJuristicMethod: tt.asrJuristicMethod,
				PreciseToSeconds:  true,
			}

			// Open test file
			filePath := fp.Join("test", tt.testFile)
			testFile, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer testFile.Close()

			// Parse test file
			csvReader := csv.NewReader(testFile)
			for {
				record, err := csvReader.Read()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						t.Fatalf("failed to parse test file: %v\n", err)
					}
				}

				// Read row value
				strDate := record[0]
				date, _ := time.ParseInLocation("2006-01-02", strDate, tt.timezone)

				expFajr := parseCSVTime(date, record[1])
				expSunrise := parseCSVTime(date, record[2])
				expZuhr := parseCSVTime(date, record[3])
				expAsr := parseCSVTime(date, record[4])
				expMaghrib := parseCSVTime(date, record[5])
				expIsha := parseCSVTime(date, record[6])

				// Calculate prayer time
				adhan, _ := GetTimes(date, cfg)

				// Calculate diff
				diffFajr := math.Abs(adhan.Fajr.Sub(expFajr).Seconds())
				diffSunrise := math.Abs(adhan.Sunrise.Sub(expSunrise).Seconds())
				diffZuhr := math.Abs(adhan.Zuhr.Sub(expZuhr).Seconds())
				diffAsr := math.Abs(adhan.Asr.Sub(expAsr).Seconds())
				diffMaghrib := math.Abs(adhan.Maghrib.Sub(expMaghrib).Seconds())
				diffIsha := math.Abs(adhan.Isha.Sub(expIsha).Seconds())

				// Print test result
				strExpFajr := expFajr.Format("15:04")
				strExpSunrise := expSunrise.Format("15:04")
				strExpZuhr := expZuhr.Format("15:04")
				strExpAsr := expAsr.Format("15:04")
				strExpMaghrib := expMaghrib.Format("15:04")
				strExpIsha := expIsha.Format("15:04")

				strFajr := adhan.Fajr.Format("15:04:05")
				strSunrise := adhan.Sunrise.Format("15:04:05")
				strZuhr := adhan.Zuhr.Format("15:04:05")
				strAsr := adhan.Asr.Format("15:04:05")
				strMaghrib := adhan.Maghrib.Format("15:04:05")
				strIsha := adhan.Isha.Format("15:04:05")

				// Difference must be less than 60 seconds
				switch {
				case diffFajr >= 60:
					t.Errorf("fajr %s, expected %s got %s (%f)", strDate, strExpFajr, strFajr, diffFajr)
				case diffSunrise >= 60:
					t.Errorf("sunrise %s, expected %s got %s (%f)", strDate, strExpSunrise, strSunrise, diffSunrise)
				case diffZuhr >= 60:
					t.Errorf("zuhr %s, expected %s got %s (%f)", strDate, strExpZuhr, strZuhr, diffZuhr)
				case diffAsr >= 60:
					t.Errorf("asr %s, expected %s got %s (%f)", strDate, strExpAsr, strAsr, diffAsr)
				case diffMaghrib >= 60:
					t.Errorf("maghrib %s, expected %s got %s (%f)", strDate, strExpMaghrib, strMaghrib, diffMaghrib)
				case diffIsha >= 60:
					t.Errorf("isha %s, expected %s got %s (%f)", strDate, strExpIsha, strIsha, diffIsha)
				}
			}
		})
	}
}

func parseCSVTime(date time.Time, csvRecord string) time.Time {
	csvTime, _ := time.Parse("15:04", csvRecord)
	minutes := csvTime.Hour()*60 + csvTime.Minute()
	return date.Add(time.Minute * time.Duration(minutes))
}
