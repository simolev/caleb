package caleb

import (
	"fmt"
	"math"
	"time"
)

func main() {
	j := jewishDate{5738, 11, 17}
	g := JewishToGregorian(j)
	fmt.Println(j, "->", g)
	j = GregorianToJewish(g)
	fmt.Println(g, "->", j)
}

func MonthsSinceFirstMolad(shana int) (n int) {
	return int(math.Floor(((float64(shana) * 235) - 234) / 19))
}

func IsMehubberet(shana int) (mehubberet bool) {
	// Shana Mehubberet is year 3, 6, 8, 11, 14, 17, 19 in a 19 year cycle
	switch shana % 19 {
	case 3, 6, 8, 11, 14, 17, 0:
		return true
	}
	return false
}

func RoshHashana(shana int) (roshHashana time.Time) {
	var nMonthsSinceFirstMolad int
	var nChalakim int
	var nHours int
	var nDays int
	var nDayOfWeek int

	nMonthsSinceFirstMolad = MonthsSinceFirstMolad(shana)
	nChalakim = 793 * nMonthsSinceFirstMolad
	nChalakim += 204
	// carry the excess Chalakim over to the hours
	nHours = int(math.Floor(float64(nChalakim) / 1080))
	nChalakim = nChalakim % 1080

	nHours += nMonthsSinceFirstMolad * 12
	nHours += 5

	// carry the excess hours over to the days
	nDays = int(math.Floor(float64(nHours) / 24))
	nHours = nHours % 24

	nDays += 29 * nMonthsSinceFirstMolad
	nDays += 2

	// figure out which day of the week the molad occurs.
	// Sunday = 1, Moday = 2 ..., Shabbos = 0
	nDayOfWeek = nDays % 7

	if !IsMehubberet(shana) && nDayOfWeek == 3 && (nHours*1080)+nChalakim >= (9*1080)+204 {
		// This prevents the year from being 356 days. We have to push
		// Rosh Hashanah off two days because if we pushed it off only
		// one day, Rosh Hashanah would comes out on a Wednesday. Check
		// the Hebrew year 5745 for an example.
		nDayOfWeek = 5
		nDays += 2
	} else if IsMehubberet(shana-1) && nDayOfWeek == 2 && (nHours*1080)+nChalakim >= (15*1080)+589 {
		// This prevents the previous year from being 382 days. Check
		// the Hebrew Year 5766 for an example. If Rosh Hashanah was not
		// pushed off a day then 5765 would be 382 days
		nDayOfWeek = 3
		nDays += 1
	} else {
		// see rule 2 above. Check the Hebrew year 5765 for an example
		if nHours >= 18 {
			nDayOfWeek += 1
			nDayOfWeek = nDayOfWeek % 7
			nDays += 1
		}
		// see rule 1 above. Check the Hebrew year 5765 for an example
		switch nDayOfWeek {
		case 1, 4, 6:
			nDayOfWeek += 1
			nDayOfWeek = nDayOfWeek % 7
			nDays += 1
		}
	}
	nDays -= 2067025
	roshHashana = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC) // 2067025 days after creation
	roshHashana = roshHashana.AddDate(0, 0, nDays)
	return
}

func LengthOfYear(shana int) int {
	return int(math.Round(RoshHashana(shana+1).Sub(RoshHashana(shana)).Hours() / 24))
}

type jewishDate struct {
	yom     int
	chodesh int
	shana   int
}

func (t jewishDate) String() string {
	return fmt.Sprintf("%02d-%02d-%04d", t.yom, t.chodesh, t.shana)
}
func (t jewishDate) Full() string {
	chodashim := [13]string{"Tishri", "Cheshvan", "Kislev", "Tevet", "Shevat", "Adar", "Adar II", "Nisan", "Yiar", "Sivan", "Tamuz", "Av", "Elul"}
	return fmt.Sprintf("%02d %s %04d", t.yom, chodashim[t.chodesh-1], t.shana)
}
func (t jewishDate) Serialize() (int, int, int) {
	return t.yom, t.chodesh, t.shana
}

func JewishToGregorian(j jewishDate) (gregorian time.Time) {
	var nLengthOfYear int
	var bLeap bool
	var dGreg time.Time
	var nMonth int
	var nMonthLen int
	var bHaser bool
	var bShalem bool

	yom, chodesh, shana := j.Serialize()
	bLeap = IsMehubberet(shana)
	nLengthOfYear = LengthOfYear(shana)

	// The regular length of a non-leap year is 354 days.
	// The regular length of a leap year is 384 days.
	// On regular years, the length of the months are as follows
	//   Tishrei (1)   30
	//   Cheshvan(2)   29
	//   Kislev  (3)   30
	//   Teves   (4)   29
	//   Shevat  (5)   30
	//   Adar A  (6)   30     (only valid on leap years)
	//   Adar    (7)   29     (Adar B for leap years)
	//   Nisan   (8)   30
	//   Iyar    (9)   29
	//   Sivan   (10)  30
	//   Tamuz   (11)  29
	//   Av      (12)  30
	//   Elul    (13)  29
	// If the year is shorter by one less day, it is called a haser
	// year. Kislev on a haser year has 29 days. If the year is longer
	// by one day, it is called a shalem year. Cheshvan on a shalem
	// year is 30 days.

	bHaser = (nLengthOfYear == 353 || nLengthOfYear == 383)
	bShalem = (nLengthOfYear == 355 || nLengthOfYear == 385)

	// get the date for Tishrei 1
	dGreg = RoshHashana(shana)

	// Now count up days within the year
	for nMonth = 1; nMonth <= chodesh-1; nMonth++ {
		switch nMonth {
		case 1, 5, 8, 10, 12:
			nMonthLen = 30
		case 4, 7, 9, 11, 13:
			nMonthLen = 29
		case 6:
			if bLeap {
				nMonthLen = 30
			} else {
				nMonthLen = 0
			}
		case 2:
			if bShalem {
				nMonthLen = 30
			} else {
				nMonthLen = 29
			}
		case 3:
			if bHaser {
				nMonthLen = 29
			} else {
				nMonthLen = 30
			}
		}
		dGreg = dGreg.AddDate(0, 0, nMonthLen)
	}
	dGreg = dGreg.AddDate(0, 0, yom-1)
	return dGreg
}

func GregorianToJewish(dGreg time.Time) jewishDate {
	var nYearH int
	var nMonthH int
	var nDateH int
	var nOneMolad float64
	var nAvrgYear float64
	var nDays int
	var dTishrei1 time.Time
	var nLengthOfYear int
	var bLeap bool
	var bHaser bool
	var bShalem bool
	var nMonthLen int
	var bWhile = true
	d1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	// The basic algorythm to get Hebrew date for the Gregorian date dGreg.
	// 1) Find out how many days dGreg is after creation.
	// 2) Based on those days, estimate the Hebrew year
	// 3) Now that we a good estimate of the Hebrew year, use brute force to
	//    find the Gregorian date for Tishrei 1 prior to or equal to dGreg
	// 4) Add to Tishrei 1 the amount of days dGreg is after Tishrei 1

	// Figure out how many days are in a month.
	// 29 days + 12 hours + 793 chalakim
	nOneMolad = 29.0 + (12.0 / 24.0) + (793.0 / (1080.0 * 24.0))
	// Figure out the average length of a year. The hebrew year has exactly
	// 235 months over 19 years.
	nAvrgYear = nOneMolad * (235.0 / 19.0)
	// Get how many days dGreg is after creation. See note as to why I
	// use 1/1/1900 and add 2067025
	nDays = int(math.Round((dGreg.Sub(d1900).Hours() / 24)))
	nDays += 2067025 // 2067025 days after creation
	// Guess the Hebrew year. This should be a pretty accurate guess.
	nYearH = int(math.Floor(float64(nDays)/nAvrgYear) + 1)
	// Use brute force to find the exact year nYearH. It is the Tishrei 1 in
	// the year <= dGreg.
	dTishrei1 = RoshHashana(nYearH)

	if SameDate(dTishrei1, dGreg) {
		// If we got lucky and landed on the exact date, we can stop here
		nMonthH = 1
		nDateH = 1
	} else {
		// Here is the brute force.  Either count up or count down nYearH
		// until Tishrei 1 is <= dGreg.
		if dTishrei1.Sub(dGreg).Hours() < 0 {
			// If Tishrei 1, nYearH is less than dGreg, count nYearH up.
			for RoshHashana(nYearH+1).Sub(dGreg).Hours() <= 0 {
				nYearH += 1
			}
		} else {
			// If Tishrei 1, nYearH is greater than dGreg, count nYearH down.
			nYearH -= 1
			for RoshHashana(nYearH).Sub(dGreg).Hours() > 0 {
				nYearH -= 1
			}
		}

		// Subtract Tishrei 1, nYearH from dGreg. That should leave us with
		// how many days we have to add to Tishrei 1
		nDays = int(math.Round((dGreg.Sub(RoshHashana(nYearH)).Hours() / 24)))
		// Find out what type of year it is so that we know the length of the
		// months
		nLengthOfYear = LengthOfYear(nYearH)
		bHaser = nLengthOfYear == 353 || nLengthOfYear == 383
		bShalem = nLengthOfYear == 355 || nLengthOfYear == 385
		bLeap = IsMehubberet(nYearH)

		// Add nDays to Tishrei 1.
		nMonthH = 1
		for bWhile {
			switch nMonthH {
			case 1, 5, 6, 8, 10, 12:
				nMonthLen = 30
			case 4, 7, 9, 11, 13:
				nMonthLen = 29
			case 2: // Cheshvan, see note above
				if bShalem {
					nMonthLen = 30
				} else {
					nMonthLen = 29
				}
			case 3: // Kislev, see note above
				if bHaser {
					nMonthLen = 29
				} else {
					nMonthLen = 30
				}
			}

			if nDays >= nMonthLen {
				bWhile = true
				if bLeap || nMonthH != 5 {
					nMonthH++
				} else {
					// We can skip Adar A (6) if its not a leap year
					nMonthH += 2
				}
				nDays -= nMonthLen
			} else {
				bWhile = false
			}
		}
		//Add the remaining days to Date
		nDateH = nDays + 1
	}
	if nMonthH == 7 && bLeap == false {
		nMonthH = 6
	}
	return jewishDate{shana: nYearH, chodesh: nMonthH, yom: nDateH}
}

func SameDate(d1, d2 time.Time) bool {
	return d1.Day() == d2.Day() && d1.Month() == d2.Month() && d1.Year() == d2.Year()
}
