# Go Prayer [![Go Reference][doc-pkg-badge]][doc-pkg-url] [![Classic Go Reference][doc-godocs-badge]][doc-godocs-url]

Go Prayer is a Go package for calculating prayer/salah times for an entire year in the specified location. It uses [SPA][spa] algorithm from [`go-sampa`][go-sampa] package to calculate the Sun events which used to determine the prayer times.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installation](#installation)
- [Features](#features)
- [API](#api)
- [Calculation Result](#calculation-result)
- [Fajr and Isha Conventions](#fajr-and-isha-conventions)
- [Asr Conventions](#asr-conventions)
- [Higher Latitude Conventions](#higher-latitude-conventions)
- [FAQ](#faq)
- [License](#license)

## Installation

To use this package, make sure your project use Go 1.20 or above, then run the following command via terminal:

```
go get -u -v github.com/hablullah/go-prayer
```

## Features

- Seamlessly handle DST times.
- Should be mathematically accurate thanks to `go-sampa` package.
- Provides several twilight conventions for calculating Fajr and Isha time.
- Provides several adapter to calculate prayer times in area with higher latitude (>45°) where Sun may not rise or set for the entire day.

## API

You can check the Go documentations to see the available APIs. However, the main interest in this package is `Calculate` function, which calculate the prayer schedule for entire year.

For [example](example/main.go), here we want to get prayer times in Jakarta for 2023:

```go
package main

import (
	"time"

	"github.com/hablullah/go-prayer"
)

func main() {
	// Calculate prayer schedule in Jakarta for 2023
	asiaJakarta, _ := time.LoadLocation("Asia/Jakarta")
	schedules, _ := prayer.Calculate(prayer.Config{
		Latitude:            -6.14,
		Longitude:           106.81,
		Timezone:            asiaJakarta,
		TwilightConvention:  prayer.Kemenag(),
		AsrConvention:       prayer.Shafii,
		HighLatitudeAdapter: prayer.NearestLatitude(),
		PreciseToSeconds:    true,
	}, 2023)
}
```

You can also adjust the calculation result by specifying it in `Corrections` field in `Configuration`.

## Calculation Result

There are five times that will be calculated by this package:

1. **Fajr** is the time when the sky begins to lighten (dawn) after previously completely dark. The exact time is different between several conventions, however all of them agree that it occured within astronomical twilight when the Sun is between 12 degrees and 18 degrees below the horizon.

2. **Sunrise** is the moment when the upper limb of the Sun appears on the horizon in the morning. In theory, the sunrise time is affected by the elevation of a location. However, the change is quite neglible, only around 1 minute for every 1.5 km. Because of this, most calculators will ignore the elevation and treat the earth as a simple spherical ball.

3. **Zuhr** is the time when the Sun begins to decline after reaching the highest point in the sky, so a bit after solar noon. However, there are some difference opinions on when azan for Zuhr should be commenced.

   There is a hadith that forbade Muslim to pray exactly at the solar noon, and instead we should wait until the Sun has descended a bit to the west. From this, there are two different opinion on when azan should be announced. First, the azan should be announced after the Sun has descended a bit (noon + 1-2 minute). Second, the azan is announced right on solar noon since the prayer will be done later anyway (after iqamah).

   This package by default will use the second opinion. However, you can adjust the time by specifying it in `Corrections` field in `Configuration`.

4. **Asr** is the time when the length of any object's shadow reaches a factor of the length of the object itself plus the length of that object's shadow at noon. With that said, Asr time is calculated by measuring the length of shadow of an object, relative to the height of the object itself.

5. **Maghrib** is the time when the upper limb of the Sun disappears below the horizon, so Maghrib is equal with sunset time. Like sunrise, the time for Maghrib might be affected by the elevation of a location. However, since the change is neglible most calculators will ignore the elevation and treat the earth as a simple spherical ball.

6. **Isha** is the time at which darkness falls and after this point the sky is no longer illuminated. The exact time is different between several conventions. Most of them agree that it occured within astronomical twilight when the Sun is between 12 degrees and 18 degrees below the horizon. However there are also some conventions where the Isha time is started after fixed Maghrib duration.

## Fajr and Isha Conventions

Since there are so many Muslim from different cultures and locations, there are several conventions for calculating prayer times. For Fajr and Isha, all conventions agree that they occured within astronomical twilight, however there are differences in the value of Sun altitude. Special case for Isha, there are some conventions that uses fixed duration after Maghrib.

For these conventions, there is something called _Fajr angle_ and _Isha angle_ which basically the value of Sun altitude **below** the horizon. So, if Fajr angle is 18 degrees, then it means the Sun altitude is -18 degrees.

| No  |    Name     | Fajr Angle | Isha Angle | Maghrib Duration |                                                           Description                                                            |
| :-: | :---------: | :--------: | :--------: | :--------------: | :------------------------------------------------------------------------------------------------------------------------------: |
|  1  |     MWL     |     18     |     17     |                  | Calculation method from Muslim World League, usually used in Europe, Far East and parts of America. Default in most calculators. |
|  2  |    ISNA     |     15     |     15     |                  |                        Calculation method from Islamic Society of North America, used in Canada and USA.                         |
|  3  | Umm al-Qura |    18.5    |            |    90 minutes    |                       Calculation method from Umm al-Qura University in Makkah which used in Saudi Arabia.                       |
|  4  |    Gulf     |    19.5    |            |    90 minutes    |                       Calculation method that often used by countries in Gulf region like UAE and Kuwait.                        |
|  5  |  Algerian   |     18     |     17     |                  |                            Calculation method from Algerian Ministry of Religious Affairs and Wakfs.                             |
|  6  |   Karachi   |     18     |     18     |                  |                                 Calculation method from University of Islamic Sciences, Karachi.                                 |
|  7  |   Diyanet   |     18     |     17     |                  |                                   Calculation method from Turkey's Diyanet İşleri Başkanlığı.                                    |
|  8  |    Egypt    |    19.5    |    17.5    |                  |                                  Calculation method from Egyptian General Authority of Survey.                                   |
|  9  |  Egypt Bis  |     20     |     18     |                  |                              Another calculation method from Egyptian General Authority of Survey.                               |
| 10  |   Kemenag   |     20     |     18     |                  |                                  Calculation method from Kementerian Agama Republik Indonesia.                                   |
| 11  |    MUIS     |     20     |     18     |                  |                                      Calculation method from Majlis Ugama Islam Singapura.                                       |
| 12  |    JAKIM    |     20     |     18     |                  |                                     Calculation method from Jabatan Kemajuan Islam Malaysia.                                     |
| 13  |    UOIF     |     12     |     12     |                  |                              Calculation method from Union Des Organisations Islamiques De France.                               |
| 14  |  France15   |     15     |     15     |                  |                               Calculation method for France region with Fajr and Isha both at 15°.                               |
| 15  |  France18   |     18     |     18     |                  |                           Another calculation method for France region with Fajr and Isha both at 18°.                           |
| 16  |   Tunisia   |     18     |     18     |                  |                                 Calculation method from Tunisian Ministry of Religious Affairs.                                  |
| 17  |   Tehran    |    17.7    |     14     |                  |                             Calculation method from Institute of Geophysics at University of Tehran.                             |
| 18  |   Jafari    |     16     |     14     |                  |                     Calculation method from Shia Ithna Ashari that used in some Shia communities worldwide.                      |

These conventions are gatehered from various sources:

- [PrayTimes.org][angle-praytimes]
- [IslamicFinder.us][angle-islamicfinder]
- [MuslimPro.com][angle-muslimpro]

## Asr Conventions

As mentioned before, Asr time is determined by the length of shadow of an object. However there are different opinions on how long the shadow should be:

- In Hanafi school, Asr started when shadow length is **twice** the length of object + shadow length at noon.
- In Shafi'i school, Asr started when shadow length is **equal** the length of object + shadow length at noon.

## Higher Latitude Conventions

In locations at higher latitude, Sun might never rise or set for an entire day. In these abnormal periods, the determination of Fajr, Maghrib and Isha is not possible to calculate using the normal methods. This problem has been explained in detail by [PrayerTimes.dk][high-lat-introduction].

To overcome this issue, several solutions have been proposed by Muslim scholars:

1. **Follow schedules in Mecca**

   This convention based on Fatwa from Dar Al Iftah Al Misrriyah number 2806 dated at 2010-08-08. They propose that area with higher latitude to follows the schedule in Mecca when abnormal days occured, using transit time as the common point. Here the day is considered "abnormal" when there are no true night, or the day length is less than 4 hours.

   To prevent sudden schedule changes, this method uses transition period for maximum one month before and after the abnormal periods.

   This method doesn't require the sunrise and sunset to be exist in a day, so it's usable for area in extreme latitudes (>=65 degrees).

   For more detail, check out [PrayerTimes.dk][high-lat-mecca]. If you want to use this convention, you can do so by using `Mecca()` or `AlwaysMecca()` as `HighLatitudeAdapter` in config.

2. **Local Relative Estimation**

   "Local Relative Estimation" is method that created by cooperation between Fiqh Council of Muslim World League and Islamic Crescents' Observation Project (ICOP). In short, this method uses average percentage to calculate Fajr and Isha time for abnormal times.

   This method only estimates time for Isha and Fajr and require sunrise and sunset time. Therefore it's not suitable for area in extreme latitude (>=65 degrees).

   For more detail, check out [ICOP's site][high-lat-local-relative]. If you want to use this convention, you can do so by using `LocalRelativeEstimation()` as `HighLatitudeAdapter` in config.

3. **Use the schedule of last normal day before abnormal periods**

   In this method, the schedule for "abnormal" days will be taken from the schedule of the last "normal" day. This adapter doesn't require the sunrise and sunset to be exist in a day, so it's usable for area in extreme latitudes (>=65 degrees).

   For more detail, check out this [paper][high-lat-nearest-day]. If you want to use this convention, you can do so by using `NearestDay()` as `HighLatitudeAdapter` in config.

4. **Follow schedules in nearest latitude**

   In this method, the schedules will be estimated using percentage of schedule in location at the nearest 45 degrees latitude. For example, Amsterdam is located at 52°22'N / 4°54'E. Using this method, the schedule will be estimated using location at 45°00'N / 4°54'E. This method will change the schedule for entire year to prevent sudden changes in fasting time.

   This adapter only estimates time for Isha and Fajr and require sunrise and sunset time, therefore it's not suitable for area in extreme latitude (>=65 degrees). For those area, instead of using percentage, this method recommends to use schedule from the nearest latitude as it is without any change.

   For more detail, check out this article from [IslamOnline.net][high-lat-nearest-latitude]. This method also briefly mentioned in Islamicity's [paper][high-lat-nearest-day].

   If you want to use this convention, you can do so by using `NearestLatitude()` or `NearestLatitudeAsIs()` as `HighLatitudeAdapter` in config.

   Another alternative of this method is [proposed][high-lat-shari-normal-day] by Mohamed Nabeel Tarabishy, Ph.D. He proposes that a normal day is defined as day when the fasting period is between 10h17m and 17h36m. If the day is "abnormal" then the schedule is calculated using location at the nearest 45 degrees latitude.

   However, using his method there will be sudden changes in prayer schedules. To avoid this issue, the author has given suggestion to just use the schedule from 45° for entire year as has explained before.

   So, following that suggestion, we don't recommend you to use his method. If you still want to try, you can do so by using `ShariNormalDay()` as `HighLatitudeAdapter` in config.

5. **Isha and Fajr calculated using angle-based method**

   In this method, the night period is divided into several parts, depending on the value of twilight angle for Fajr and Isha. For example, let a be the twilight angle for Isha, and let t = a/60. The period between sunset and sunrise is divided into t parts. Isha begins after the first part. So, if the twilight angle for Isha is 15, then Isha begins at the end of the first quarter (15/60) of the night. Time for Fajr is calculated similarly.

   This adapter depends on sunrise and sunset time, so it might not be suitable for area in extreme latitudes (>=65 degrees).

   For more detail, check out this article by [PrayTimes.org][high-lat-angle-based]. If you want to use this convention, you can do so by using `AngleBased()` as `HighLatitudeAdapter` in config.

6. **Isha and Fajr at one-seventh of the night**

   In this method, the night period is divided into seven parts. Isha starts when the first seventh part ends, and Fajr starts when the last seventh part starts.

   This adapter depends on sunrise and sunset time, so it might not be suitable for area in extreme latitudes (>=65 degrees).

   For more detail, check out this article by [PrayTimes.org][high-lat-angle-based]. If you want to use this convention, you can do so by using `OneSeventhNight()` as `HighLatitudeAdapter` in config.

7. **Isha and Fajr in the middle of the night**

   In this method, the night period is divided into two halves. The first half is considered to be the "night" and the other half as "day break". Fajr and Isha in this method are assumed to be at mid-night during the abnormal periods.

   This adapter depends on sunrise and sunset time, so it might not be suitable for area in extreme latitudes (>=65 degrees).

   For more detail, check out this article by [PrayTimes.org][high-lat-angle-based]. If you want to use this convention, you can do so by using `MiddleNight()` as `HighLatitudeAdapter` in config.

## FAQ

1. **Does the elevation affects calculation result?**

   Yes, it will affect the result for sunrise and Maghrib time. If the elevation is very high, it can affect the times by a couple of minutes thanks to atmospheric refraction. However, most apps that I know prefer to set elevation to zero, which means every locations will be treated as located in sea level.

2. **Are the calculation results are really accurate up to seconds?**

   While the results of this package are in seconds, it's better to not expect it to be exactly accurate to seconds and instead treat it as minute rounding suggestions.

3. **Why are the Fajr, sunrise, Maghrib and Isha times occured in different day?**

   In this package, every times are connected to transit time (the time when Sun reach meridian). However, in area with higher latitude sometime the Sun will never rise nor set for the entire day. In this case, Fajr and sunrise might occur yesterday and the Maghrib and Isha might occur tomorrow.

4. **All of schedules are consistently missed by several minutes!**

   This package calculates the prayer times by utilizing Go SAMPA package to calculate Sun events in a day. That package uses several mathematic and astronomy formula, and has been tested with several astronomy applications and we found that it's pretty accurate within a minute. With that said, we hope the prayer times that calculated by this package should be mathematically accurate.

   However, in some locations the local scholar usually will adjust the schedule forward or backward for a few minutes. This is done for safety, to prevent the azan or prayer commenced before the actual time begin. And as mentioned before, in some country azan for Zuhr will be delayed for several minutes to follow the hadith that prohibit prayer when Sun reach its meridian.

   In case like this, you can easily adjust the time using `Corrections` field in configuration.

5. **I live in area with higher latitude, and no conventions provided in this package matches with the official schedule used in my area!**

   Since prayer in higher latitude is a matter of _ijtihad_, there are no definitive final texts pertaining to it. Therefore it's allowed for local Islamic bodies to specify their own conventions in order to save the Muslims living in that area from inconvenience and difficulty. Thanks to this, it's possible the convention that used in your area is not provided by this package.

   Fortunately, `HighLatitudeAdapter` is simply a function that defined like this:

   ```go
   type HighLatitudeAdapter func(cfg Config, year int, currentSchedules []Schedule) []Schedule
   ```

   So, if the convention is not available in this package but you know how the convention works, you can simply define it on your own. It would be even better if you open PR to add it to this package.

   If you know how the convention works but don't want to code it yourself, feel free to open an issue so we could add it to this package.

## License

Go-Prayer is distributed using [MIT] license.

[doc-pkg-badge]: https://img.shields.io/badge/-pkg.go-007d9c?logo=go&labelColor=gray&logoColor=white
[doc-pkg-url]: https://pkg.go.dev/github.com/hablullah/go-prayer
[doc-godocs-badge]: https://img.shields.io/badge/-godocs-375eab?logo=go&labelColor=gray&logoColor=white
[doc-godocs-url]: https://godocs.io/github.com/hablullah/go-prayer
[spa]: https://midcdmz.nrel.gov/spa/
[go-sampa]: https://github.com/hablullah/go-sampa
[angle-praytimes]: http://praytimes.org/wiki/Calculation_Methods
[angle-islamicfinder]: http://www.islamicfinder.us/index.php/api/index
[angle-muslimpro]: https://www.muslimpro.com/en/prayer-times
[high-lat-introduction]: https://www.prayertimes.dk/story.html
[high-lat-mecca]: https://www.prayertimes.dk/fatawa.html
[high-lat-local-relative]: https://www.astronomycenter.net/latitude.html?l=en
[high-lat-nearest-day]: https://www.islamicity.com/prayertimes/Salat.pdf
[high-lat-nearest-latitude]: https://fiqh.islamonline.net/en/praying-and-fasting-at-high-latitudes/
[high-lat-shari-normal-day]: https://www.astronomycenter.net/pdf/tarabishyshigh_2014.pdf
[high-lat-angle-based]: http://praytimes.org/calculation
[mit]: http://choosealicense.com/licenses/mit/
