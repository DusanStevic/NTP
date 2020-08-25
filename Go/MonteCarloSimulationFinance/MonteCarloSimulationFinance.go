package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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
	data              []float64
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

// Differencing time series = Shifting and lagging time series
func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculatePeriodicDailyReturn() {
	for i := 1; i < len(monteCarloSimulationFinance.timeSeries); i++ {
		monteCarloSimulationFinance.data = append(monteCarloSimulationFinance.data,
			math.Log(monteCarloSimulationFinance.timeSeries[i]/monteCarloSimulationFinance.timeSeries[i-1]))
	}
}

// https://github.com/gonum/gonum
// https://godoc.org/gonum.org/v1/gonum/stat#Mean
func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateAverageDailyReturn() float64 {
	/* 	computes the weighted mean of the dataset.
	   	we don't have any weights (ie: all weights are 1)
	   	so we just pass a nil slice. */
	return stat.Mean(monteCarloSimulationFinance.data, nil)
}

// https://github.com/gonum/gonum
// https://godoc.org/gonum.org/v1/gonum/stat#Variance
func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateVariance() float64 {
	/* 	computes the weighted variance of the dataset.
	   	we don't have any weights (ie: all weights are 1)
	   	so we just pass a nil slice. */
	return stat.Variance(monteCarloSimulationFinance.data, nil)

}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateStandardDeviation() float64 {
	return math.Sqrt(monteCarloSimulationFinance.calculateVariance())

}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateDrift() float64 {
	return monteCarloSimulationFinance.calculateAverageDailyReturn() - monteCarloSimulationFinance.calculateVariance()/2
}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) calculateRandomValue() float64 {
	return monteCarloSimulationFinance.calculateStandardDeviation() * monteCarloSimulationFinance.calculateZScore()
}

func main() {
	monteCarloSimulationFinance := MonteCarloSimulationFinance{
		tickerSymbol:      "AAPL",
		startDate:         "2000-01-01",
		endDate:           "2020-01-01",
		numberOfProcesses: 4}
	monteCarloSimulationFinance.dataAcquisition()
	monteCarloSimulationFinance.calculatePeriodicDailyReturn()

	f, err := os.Create("C:\\Users\\Dule\\Desktop\\telep.txt") // creating...
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return
	}
	defer f.Close()
	for i := 0; i < len(monteCarloSimulationFinance.timeSeries); i++ { // Generating...
		_, err = f.WriteString(fmt.Sprintf("%f\r\n", monteCarloSimulationFinance.timeSeries[i])) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	fmt.Println(monteCarloSimulationFinance.calculateAverageDailyReturn())
}
