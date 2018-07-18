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
	j := JewishDate{Shana: 5343, Chodesh: 1, Yom: 9}
	g := JewishToGregorian(j)
	fmt.Print(j, " -> ", g.Format("2006-01-02"))
	j = GregorianToJewish(g)
	fmt.Println("->", j)
	d := time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10000; i++ {
		d = d.AddDate(0, 0, rnd.Intn(365))
		j = GregorianToJewish(d)
		g = JewishToGregorian(j)
		switch g {
		case d:
			//fmt.Println(d, "<==>", j)
		default:
			t.Error(d, "->", j, "->", g)
		}

	}
	fmt.Println(d.Format("2006-01-02"), "->", j, "->", g.Format("2006-01-02"))
}
