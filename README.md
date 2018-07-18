# caleb
Jewish dates conversion tool

### Usage:
`go get github.com/simolev/caleb`

```go
package main

import (
      "fmt"
      "github.com/simolev/caleb"
      "time"
)

func main() {
      j1 := caleb.JewishDate{Shana: 5779, Chodesh: 7, Yom: 25} // 25 Adar II 5779
      g1 := caleb.JewishToGregorian(j1)                        //
      fmt.Println(j1, "=>", g1.Format("2006-01-02"))           // 25 Adar II 5779 => 2019-04-01
      g2 := time.Date(2018, 8, 11, 0, 0, 0, 0, time.UTC)       // 2018-08-11
      j2 := caleb.GregorianToJewish(g2)                        //
      fmt.Println(g2.Format("2006-01-02"), "=>", j2)           // 2018-08-11 => 30 Av 5778
}
```

### Disclaimer:
**a)** It is not accurate for dates before Gregorian **1582-10-15**. In the Gregorian calendar, there are 10 missing days between 1582-10-15 and 1582-10-04. Those days never occurred, but nonetheless they seem to exist in go's implementation, and the only solution I see at the moment would be to add manual correction to the algorithm.

**b)** This is quite new and would require more thorough testing before being relied upon.

**c)** Code could be more idiomatic and optimized.

### Credits:
Thanks to info@dafaweek.com.  
Converted from javascript: http://www.dafaweek.com/HebCal/HebCalSampleSource.php  
See also http://www.dafaweek.com/hebcal/hebcalvb6.php
