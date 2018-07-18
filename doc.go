// Package caleb provides a JewishDate type and conversion to and from time.Time
//
//package main
//
//import (
//	"fmt"
//	"github.com/simolev/caleb"
//	"time"
//)
//
//func main() {
//	j1 := caleb.JewishDate{Shana: 5779, Chodesh: 7, Yom: 25} // 25 Adar II 5779
//	g1 := caleb.JewishToGregorian(j1)                        //
//	fmt.Println(j1, "=>", g1)                                // 25 Adar II 5779 => 2019-04-01 00:00:00 +0000 UTC
//	g2 := time.Date(2018, 8, 11, 0, 0, 0, 0, time.UTC)       // 2018-08-11
//	j2 := caleb.GregorianToJewish(g2)                        //
//	fmt.Println(g2, "=>", j2)                                // 2018-08-11 00:00:00 +0000 UTC => 30 Av 5778
//}
package caleb
