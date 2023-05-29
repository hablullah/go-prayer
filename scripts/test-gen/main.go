package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Clean up dst dir
	dstDir := "internal/datatest"
	os.RemoveAll(dstDir)
	os.MkdirAll(dstDir, os.ModePerm)

	// Generate common files
	err := genCommonFiles(dstDir)
	checkError(err)

	// Generate test data for each location
	for _, location := range testLocations {
		err = genTestData(location, dstDir)
		checkError(err)
	}
}

func genCommonFiles(dstDir string) error {
	// Write package header and imports
	var sb strings.Builder
	sb.WriteString("package datatest\n")
	sb.WriteString(`import (
		"time"
		"regexp"
		"strconv"

		"github.com/hablullah/go-prayer"
	)` + "\n\n")

	sb.WriteString("var rxDT = regexp.MustCompile(`" +
		`(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})` +
		"`)\n\n")

	// Put struct for test data
	sb.WriteString(`type TestData struct {
		Name      string
		Latitude  float64
		Longitude float64
		Timezone  *time.Location
		Schedules []prayer.Schedule
	}` + "\n\n")

	// Put helper function
	sb.WriteString(`func schedule(date string, fajr, sunrise, zuhr, asr, maghrib, isha string, tz *time.Location) prayer.Schedule {
		return prayer.Schedule {
			Date:    date,
			Fajr:    st(fajr, tz),
			Sunrise: st(sunrise, tz),
			Zuhr:    st(zuhr, tz),
			Asr:     st(asr, tz),
			Maghrib: st(maghrib, tz),
			Isha:    st(isha, tz),
		}
	}` + "\n\n")

	sb.WriteString(`
	func st(s string, tz *time.Location) time.Time {
		parts := rxDT.FindStringSubmatch(s)
		if len(parts) != 7 {
			return time.Time{}
		}
	
		year, _ := strconv.Atoi(parts[1])
		month, _ := strconv.Atoi(parts[2])
		day, _ := strconv.Atoi(parts[3])
		hour, _ := strconv.Atoi(parts[4])
		minute, _ := strconv.Atoi(parts[5])
		second, _ := strconv.Atoi(parts[6])
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, tz)
	}` + "\n\n")

	// Format code
	bt, err := format.Source([]byte(sb.String()))
	if err != nil {
		return err
	}

	// Save to file
	dstPath := filepath.Join(dstDir, "common.go")
	return os.WriteFile(dstPath, bt, os.ModePerm)
}

func genTestData(loc Location, dstDir string) error {
	// Write package header and imports
	var sb strings.Builder
	sb.WriteString("package datatest\n")
	sb.WriteString(`import (
		"time"
		
		"github.com/hablullah/go-prayer"
	)` + "\n")

	// Put the variable for timezone
	tzName := fmt.Sprintf("tz%s", loc.Name)
	sb.WriteString(fmt.Sprintf(
		"var %s, _ = time.LoadLocation(%q)\n\n",
		tzName, loc.Timezone))

	// Put the variable for location
	sb.WriteString(fmt.Sprintf(""+
		"var %s = TestData {\n"+
		"Name: %q,\n"+
		"Latitude: %f,\n"+
		"Longitude: %f,\n"+
		"Timezone: %s,\n",
		loc.Name, loc.Name, loc.Latitude, loc.Longitude, tzName))

	// Calculate and put schedules
	sb.WriteString("Schedules: []prayer.Schedule{\n")
	for _, s := range getSchedules(loc) {
		sb.WriteString(fmt.Sprintf(
			"schedule(%q,%q,%q,%q,%q,%q,%q,%s),\n",
			s.Date,
			strTime(s.Fajr),
			strTime(s.Sunrise),
			strTime(s.Zuhr),
			strTime(s.Asr),
			strTime(s.Maghrib),
			strTime(s.Isha),
			tzName,
		))
	}
	sb.WriteString("},\n")
	sb.WriteString("}\n")

	// Format code
	bt, err := format.Source([]byte(sb.String()))
	if err != nil {
		return err
	}

	// Save to file
	dstPath := filepath.Join(dstDir, strings.ToLower(loc.Name)+".go")
	return os.WriteFile(dstPath, bt, os.ModePerm)
}

func strTime(t time.Time) string {
	if t.IsZero() {
		return strings.Repeat(" ", 19)
	} else {
		return t.Format("2006-01-02 15:04:05")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
