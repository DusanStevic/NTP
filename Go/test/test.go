package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat"
)

func main() {
	xs := []float64{
		32.32, 56.98, 21.52, 44.32,
		55.63, 13.75, 43.47, 43.34,
		12.34,
	}

	fmt.Printf("data: %v\n", xs)

	// computes the weighted mean of the dataset.
	// we don't have any weights (ie: all weights are 1)
	// so we just pass a nil slice.
	mean := stat.Mean(xs, nil)

	variance := stat.Variance(xs, nil)
	stddev := math.Sqrt(variance)

	fmt.Printf("mean=     %v\n", mean)

	fmt.Printf("variance= %v\n", variance)
	fmt.Printf("std-dev=  %v\n", stddev)
	var s []float64
	x := []float64{
		10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
	}
	for i := 1; i < len(x); i++ {
		s = append(s, math.Log(x[i]/x[i-1]))
	}

	for _, num := range s {
		fmt.Println(num)
	}

}
