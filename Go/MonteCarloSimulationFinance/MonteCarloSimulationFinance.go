package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"

	"github.com/chobie/go-gaussian"
	"github.com/markcheno/go-quote"
	"gonum.org/v1/gonum/stat"
)

// MonteCarloSimulationFinance is a structure to predict stock market prices using Monte Carlo simulation.
type MonteCarloSimulationFinance struct {
	numberOfProcesses int
	timeSeries        []float64
	startDate         string
	endDate           string
	tickerSymbol      string
}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) dataAcquisition() {
	stock, _ := quote.NewQuoteFromYahoo(monteCarloSimulationFinance.tickerSymbol,
		monteCarloSimulationFinance.startDate, monteCarloSimulationFinance.endDate, quote.Daily, true)
	monteCarloSimulationFinance.timeSeries = stock.Close
}

/*  A z-table, also called the standard normal table, is a mathematical table that allows us to know
the percentage of values below (to the left) a z-score in a standard normal distribution (SND).
A z-score, also known as a standard score, indicates the number of standard deviations
a raw score lays above or below the mean. When the mean of the z-score is calculated it is always 0,
and the standard deviation (variance) is always in increments of 1. */
func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateZScore() float64 {
	normalDistribution := gaussian.NewGaussian(0, 1)
	/* 	Source code for random number generator https://play.golang.org/p/ZdFpbahgC1
	   	The default number generator is deterministic, so it'll
	   	produce the same sequence of numbers each time by default.
	   	To produce varying sequences, give it a seed that changes.
	   	Note that this is not safe to use for random numbers you
	   	intend to be secret, use `crypto/rand` for those.
	   	Seeding - Go provides a method, Seed(seed int64), that allows you
	   	to initialize this default sequence. Implementation is slow
	   	to make it faster rand.Seed(time.Now().UnixNano()) is added.
	   	Seed is the current time, converted to int64 by UnixNano.
	   	Gives constantly changing numbers. */

	// Seed
	s := rand.NewSource(time.Now().UnixNano())
	// Randomly changing numbers.
	r := rand.New(s)
	/* 	Call the resulting `rand.Rand` just like the
	   	functions on the `rand` package. */
	return normalDistribution.Ppf(r.Float64())

}

func main() {
	monteCarloSimulationFinance := MonteCarloSimulationFinance{
		tickerSymbol:      "spy",
		startDate:         "2016-01-01",
		endDate:           "2016-04-01",
		numberOfProcesses: 4}
	monteCarloSimulationFinance.dataAcquisition()

	stock, _ := quote.NewQuoteFromYahoo("spy", "2016-01-01", "2016-04-01", quote.Daily, true)
	fmt.Print(stock.CSV())

	xs := []float64{
		32.32, 56.98, 21.52, 44.32,
		55.63, 13.75, 43.47, 43.34,
		12.34,
	}
	xs = monteCarloSimulationFinance.timeSeries

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
	fmt.Println(reflect.TypeOf(stock.Close))

	fmt.Println(monteCarloSimulationFinance.calculateZScore())
}
