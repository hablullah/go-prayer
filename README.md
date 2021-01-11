# Go-Prayer

[![Go Report Card][report-badge]][report-url]
[![Go Reference][doc-badge]][doc-url]

Go-Prayer is a Go package for calculating prayer/salah times for a specific location at a specific date. It 
should be accurate enough for most case and you can adjust the times if needed. However, you can alter the
result by using `TimeCorrections` field in config.

## Usage

For example, we want to get prayer times in Jakarta, Indonesia, at 4 September 2020 :

```go
package main

import (
	"fmt"
	"time"

	"github.com/RadhiFadlillah/go-prayer"
)

func main() {
	// Prepare configuration
	cfg := prayer.Config{
		Latitude:          -6.21,
		Longitude:         106.85,
		Elevation:         50,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
		TimeCorrections: prayer.TimeCorrections{
			Fajr:    2 * time.Minute,
			Sunrise: -time.Minute,
			Zuhr:    2 * time.Minute,
			Asr:     time.Minute,
			Maghrib: time.Minute,
			Isha:    time.Minute,
		},
	}

	// Specify the date that you want to calculate. This package will use timezone from your date,
	// so make sure to set it to the timezone of the location that you want to calculate.
	zone := time.FixedZone("WIB", 7*3600)
	date := time.Date(2020, 9, 4, 0, 0, 0, 0, zone)

	// Now we just need to calculate it
	result, _ := prayer.Calculate(cfg, date)
	fmt.Println(date.Format("2006-01-02"))
	fmt.Println("Fajr    =", result.Fajr.Format("15:04"))
	fmt.Println("Sunrise =", result.Sunrise.Format("15:04"))
	fmt.Println("Zuhr    =", result.Zuhr.Format("15:04"))
	fmt.Println("Asr     =", result.Asr.Format("15:04"))
	fmt.Println("Maghrib =", result.Maghrib.Format("15:04"))
	fmt.Println("Isha    =", result.Isha.Format("15:04"))
}
```

Which will give us following results :

```
2020-09-04
Fajr    = 04:36
Sunrise = 05:49
Zuhr    = 11:54
Asr     = 15:10
Maghrib = 17:54
Isha    = 19:02
```

## Calculation Result

There are five times that will be calculated by this package:

### Fajr

**Fajr** is the time when the sky begins to lighten (dawn) after previously completely dark. The exact time
is different between several conventions, however all of them agree that it occured within astronomical
twilight when the Sun is between 12 degrees and 18 degrees below the horizon.

### Sunrise

**Sunrise** is the moment when the upper limb of the Sun appears on the horizon in the morning. In theory,
the sunrise time is affected by the elevation of a location. However, the change is quite neglible, only
around 1 minute for every 1.5 km. Because of this, most calculators will ignore the elevation and treat the
earth as a simple spherical ball.

### Zuhr

**Zuhr** is the time when the Sun begins to decline after reaching the highest point in the sky, so a bit
after solar noon. However, there are some difference opinions on when azan for Zuhr should be commenced.

There is a hadith that forbade Muslim to pray exactly at the solar noon, and instead we should wait until
the Sun has descended a bit to the west before. From this, there are two different opinion on when azan
should be announced :

1. The azan should be announced after the Sun has descended a bit (noon + 1-2 minute).
2. The azan is announced right on solar noon since the prayer will be done later anyway (after iqamah).
 
### Asr

**Asr** is the time when the length of any object's shadow reaches a factor of the length of the object
itself plus the length of that object's shadow at noon. With that said, Asr time is calculated by measuring
the length of shadow of an object, relative to the height of the object itself.

### Maghrib

**Maghrib** is the time when the upper limb of the Sun disappears below the horizon, so Maghrib is equal with
sunset time. Like sunrise, the time for Maghrib might be affected by the elevation of a location. However,
since the change is neglible, most calculators will ignore the elevation and treat the earth as a simple
spherical ball.

### Isha

**Isha** is the time at which darkness falls and after this point the sky is no longer illuminated. The
exact time is different between several conventions. Most of them agree that it occured within astronomical
twilight when the Sun is between 12 degrees and 18 degrees below the horizon. However there are also some
conventions where the Isha time is started after fixed Maghrib duration.

## Conventions

Since there are so many Muslim from different cultures and locations, there are several conventions for
calculating prayer times. Depending on their purpose, those conventions can be classified into three
categories:

### Fajr and Isha Conventions

These conventions are used to specify the time for Fajr and Isha. All conventions agree that Fajr is occured
within astronomical twilight, however there are differences in the value of Sun altitude. The same for Isha,
except for Isha there are some conventions that uses fixed duration after Maghrib.

For these conventions, there is something called *Fajr angle* and *Isha angle* which basically the value of
Sun altitude **below** the horizon. So, if Fajr angle is 18 degrees, then it means the Sun altitude is -18
degrees.

| No |     Name    | Fajr Angle | Isha Angle | Maghrib Duration |                                                            Description                                                           |
|:--:|:-----------:|:----------:|:----------:|:----------------:|:--------------------------------------------------------------------------------------------------------------------------------:|
| 1  | MWL         | 18         | 17         |                  | Calculation method from Muslim World League, usually used in Europe, Far East and parts of America. Default in most calculators. |
| 2  | ISNA        | 15         | 15         |                  | Calculation method from Islamic Society of North America, used in Canada and USA.                                                |
| 3  | Umm al-Qura | 18.5       |            | 90 minutes       | Calculation method from Umm al-Qura University in Makkah which used in Saudi Arabia.                                             |
| 4  | Gulf        | 19.5       |            | 90 minutes       | Calculation method that often used by countries in Gulf region like UAE and Kuwait.                                              |
| 5  | Algerian    | 18         | 17         |                  | Calculation method from Algerian Ministry of Religious Affairs and Wakfs.                                                        |
| 6  | Karachi     | 18         | 18         |                  | Calculation method from University of Islamic Sciences, Karachi.                                                                 |
| 7  | Diyanet     | 18         | 17         |                  | Calculation method from Turkey's Diyanet İşleri Başkanlığı.                                                                      |
| 8  | Egypt       | 19.5       | 17.5       |                  | Calculation method from Egyptian General Authority of Survey.                                                                    |
| 9  | Egypt Bis   | 20         | 18         |                  | Another calculation method from Egyptian General Authority of Survey.                                                            |
| 10 | Kemenag     | 20         | 18         |                  | Calculation method from Kementerian Agama Republik Indonesia.                                                                    |
| 11 | MUIS        | 20         | 18         |                  | Calculation method from Majlis Ugama Islam Singapura.                                                                            |
| 12 | JAKIM       | 20         | 18         |                  | Calculation method from Jabatan Kemajuan Islam Malaysia.                                                                         |
| 13 | UOIF        | 12         | 12         |                  | Calculation method from Union Des Organisations Islamiques De France.                                                            |
| 14 | France15    | 15         | 15         |                  | Calculation method for France region with Fajr and Isha both at 15°.                                                             |
| 15 | France18    | 18         | 18         |                  | Another calculation method for France region with Fajr and Isha both at 18°.                                                     |
| 16 | Tunisia     | 18         | 18         |                  | Calculation method from Tunisian Ministry of Religious Affairs.                                                                  |
| 17 | Tehran      | 17.7       | 14         |                  | Calculation method from Institute of Geophysics at University of Tehran.                                                         |
| 18 | Jafari      | 16         | 14         |                  | Calculation method from Shia Ithna Ashari that used in some Shia communities worldwide.                                          |

### Asr Conventions

As mentioned before, Asr time is determined by the length of shadow of an object. However there are
different opinions on how long the shadow should be:

- In Hanafi school, Asr started when shadow length is **twice** the length of object + shadow length at noon.
- In Shafi'i school, Asr started when shadow length is **equal** the length of object + shadow length at noon.

### Higher Latitude Conventions

In locations at higher latitude, Sun might never rise or set for an entire day. In these abnormal periods,
the determination of Fajr, Maghrib and Isha is not possible to calculate using the normal methods. To
overcome this problem, several solutions have been proposed by Muslim scholars:

- Angle-based method

	This method is created after conference between Muslim scholars in Brussels, 25-26 May 2009. It's one of
	the most common method and used by some recent prayer time calculators. Let a be the twilight angle for
	Isha, and let t = a/60. The period between sunset and sunrise is divided into t parts. Isha begins after
	the first part. For example, if the twilight angle for Isha is 15, then Isha begins at the end of the first
	quarter (15/60) of the night. Time for Fajr is calculated similarly.

- One-seventh of the night

	In this method, the period between sunset and sunrise is divided into seven parts. Isha begins after the
	first one-seventh part, and Fajr is at the beginning of the seventh part. 

- Middle of the night

	In this method, the period from sunset to sunrise is divided into two halves. The first half is 
	considered to be the "night" and the other half as "day break". Fajr and Isha in this method are assumed
	to be at mid-night during the abnormal periods. 

Those three conventions are the common conventions for calculating prayer times in higher latitude. However,
as you can see, all of them requires sunrise and sunset to be exist. So, these methods only suitable for
area between 48.6 and 66.6 latitude.

To fix this issue, Mohamed Nabeel Tarabishy, Ph.D (2014) created a new method to calculate the prayer times.
In his method, if in a day the fasting time (duration between Fajr and Sunset) is too short (less than 10h
17m) or too long (more than 17h 36m), then that day is considered abnormal. In those abnormal days, the
prayer times is calculated by setting the latitude into 45. With that said, this method only used in area
above 45 latitude. This method is named `NormalRegion` in this package.

However, by using this method there will be sudden changes in the length of the day of fasting. So, he also
proposed for area above 45 latitude to just calculate the prayer times as if the latitude is 45 degrees,
no matter if the day is abnormal or not. If you want to use this method, it's named as `ForcedNormalRegion`
in this package.

## Resources

1. Anugraha, R. 2012. _Mekanika Benda Langit_. ([PDF][pdf-rinto-anugraha])
2. Tarabishy, MN. 2014. _Salat / Fasting Time in Northern Regions_. ([PDF][pdf-tarabishy])
3. _Calculation_ by PrayTimes ([web][web-pray-times-calc])
4. _Calculation Methods_ by PrayTimes ([web][web-pray-times-calc-method])
5. _Solar Calculation Details_ by NOAA ([web][web-noaa])
6. _Prayer Times in High-Latitude Areas_ by International Astronomical Center ([web][web-iac])
7. List of Fajr and Isha conventions from IslamicFinder ([web][web-islamic-finder])
8. List of Fajr and Isha conventions from MuslimPro ([web][web-muslim-pro])

## License

Go-Prayer is distributed using [MIT] license.

[report-badge]: https://goreportcard.com/badge/github.com/hablullah/go-prayer
[report-url]: https://goreportcard.com/report/github.com/hablullah/go-prayer
[doc-badge]: https://pkg.go.dev/badge/github.com/hablullah/go-prayer.svg
[doc-url]: https://pkg.go.dev/github.com/hablullah/go-prayer

[pdf-rinto-anugraha]: https://simpan.ugm.ac.id/s/GcxKuyZWn8Rshnn
[pdf-tarabishy]: http://www.astronomycenter.net/pdf/tarabishyshigh_2014.pdf
[web-pray-times-calc]: http://praytimes.org/calculation
[web-pray-times-calc-method]: http://praytimes.org/wiki/Calculation_Methods
[web-noaa]: https://www.esrl.noaa.gov/gmd/grad/solcalc/calcdetails.html
[web-iac]: http://www.astronomycenter.net/latitude.html?l=en
[web-islamic-finder]: http://www.islamicfinder.us/index.php/api/index
[web-muslim-pro]: https://www.muslimpro.com/en/prayer-times

[MIT]: http://choosealicense.com/licenses/mit/