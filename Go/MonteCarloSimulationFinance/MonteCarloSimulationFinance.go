package main

import (
	"fmt"
	"math"
	"reflect"

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

	fmt.Println(monteCarloSimulationFinance.timeSeries)
}
