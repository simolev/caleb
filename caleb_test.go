package caleb

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestJewishToGregorian(t *testing.T) {
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	j := jewishDate{shana: 5343, chodesh: 1, yom: 9}
	g := JewishToGregorian(j)
	fmt.Println(j, "->", g)
	j = GregorianToJewish(g)
	fmt.Println(g, "->", j)
	d := time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	o := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	//n := o.Sub(d).Hours() / 24
	nDays := 2067025
	o = o.AddDate(0, 0, -nDays-10)
	fmt.Println(o, GregorianToJewish(o))
	for i := 0; i < 10000; i++ {
		d = d.AddDate(0, 0, rnd.Intn(365))
		j = GregorianToJewish(d)
		g = JewishToGregorian(j)
		switch g {
		case d:
			//fmt.Println(d, "<==>", j.Full())
		default:
			t.Error(d, "==>", j.Full(), "==>", g)
		}

	}
	fmt.Println(d, "<==>", j.Full())
}
