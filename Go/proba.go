package main

import (
	"fmt"

	"math"

	"time"
)

// Round float64 number and return int number
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Round float64 number to a given int precision
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func timeMeasure() func() float64 {
	startTime := time.Now()

	return func() float64 {

		//fmt.Print(toFixed(time.Since(startTime).Seconds(), 5))
		//fmt.Println(" seconds")
		return time.Since(startTime).Seconds()

	}
}

func vals() (int, func() float64) {
	a := timeMeasure()
	defer a()
	for i := 0; i < 100000; i++ {
		fmt.Println(i)
	}

	return 33, a
}

func main() {

	a, b := vals()
	fmt.Println(a)
	fmt.Println(b())

}
