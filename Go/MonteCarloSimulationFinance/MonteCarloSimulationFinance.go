package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

// https://github.com/markcheno/go-quote
// https://godoc.org/github.com/markcheno/go-quote#NewQuoteFromYahoo
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

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) simulationFinance(numberOfSimulations int, predictionWindowSize int, channel chan [][]float64) {
	var prediction []float64
	var predictions [][]float64
	for i := 0; i < numberOfSimulations; i++ {
		// today’s price
		prediction = append(prediction, monteCarloSimulationFinance.timeSeries[len(monteCarloSimulationFinance.timeSeries)-1])
		for j := 0; j < predictionWindowSize; j++ {
			// Next Day’s Price=Today’s Price × e^(Drift+Random Value)
			prediction = append(prediction,
				prediction[len(prediction)-1]*math.Pow(math.E, (monteCarloSimulationFinance.calculateDrift()+monteCarloSimulationFinance.calculateRandomValue())))
		}
		predictions = append(predictions, prediction)
		/* 	Setting the slice to nil is the best way to clear a slice.
		   	nil slices in go are perfectly well behaved and setting the slice to nil
		   	will release the underlying memory to the garbage collector. */
		prediction = nil
	}

	channel <- predictions

}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) mcsFinanceSerial(numberOfSimulations int, predictionWindowSize int) ([][]float64, float64) {
	startTime := time.Now()
	channel := make(chan [][]float64)
	go monteCarloSimulationFinance.simulationFinance(numberOfSimulations, predictionWindowSize, channel)
	predictions := <-channel
	executionTime := time.Since(startTime).Seconds()
	return predictions, executionTime
}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) mcsFinanceParallel(numberOfSimulations int, predictionWindowSize int) ([][][]float64, float64) {
	startTime := time.Now()
	numberOfSimulationsPerProcess := numberOfSimulations / monteCarloSimulationFinance.numberOfProcesses
	/* 	Buffered channels are useful when you know how many goroutines you have launched,
	   	want to limit the number of goroutines you will launch, or want to limit
	   	the amount of work that is queued up. */
	channel := make(chan [][]float64, monteCarloSimulationFinance.numberOfProcesses)
	for i := 0; i < monteCarloSimulationFinance.numberOfProcesses; i++ {
		go monteCarloSimulationFinance.simulationFinance(numberOfSimulationsPerProcess, predictionWindowSize, channel)
	}

	var predictions [][][]float64
	for i := 0; i < monteCarloSimulationFinance.numberOfProcesses; i++ {
		prediction := <-channel
		predictions = append(predictions, prediction)

	}
	executionTime := time.Since(startTime).Seconds()
	return predictions, executionTime
}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) exportFinanceFileSerial(serialPredictions [][]float64) {
	//One row in the output file.
	var row strings.Builder
	var serialNumber int64 = 1
	var predictions [][]float64 = serialPredictions
	f, err := os.Create("C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Finance\\GolangFinanceSerial.txt") // creating...
	if err != nil {
		fmt.Printf("Error while creating a file: %v", err)
		return
	}
	defer f.Close()
	for i := 0; i < len(predictions); i++ {
		for j := 0; j < len(predictions[i]); j++ {
			if j == 0 {
				// Write a serial number of simulations at the beginning of a file.
				var serialNumberOfSimulation string = strconv.FormatInt(serialNumber, 10)
				row.WriteString(serialNumberOfSimulation + "," + " ")
			}

			//A component in one row of the output file.
			var component string = strconv.FormatFloat(predictions[i][j], 'f', 7, 64)
			row.WriteString(component)
			if j == (len(predictions[i]) - 1) {
				row.WriteString("\r\n")
				serialNumber++
			} else {
				row.WriteString("," + " ")
			}
		}
	}
	// Writing to file
	_, err = f.WriteString(row.String())
	if err != nil {
		fmt.Printf("Error while writing a file: %v", err)
	}

}

func (monteCarloSimulationFinance *MonteCarloSimulationFinance) exportFinanceFileParallel(parallelPredictions [][][]float64) {
	//One row in the output file.
	var row strings.Builder
	var serialNumber int64 = 1
	var predictions [][][]float64 = parallelPredictions
	// Creating file
	f, err := os.Create("C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Finance\\GolangFinanceParallel.txt")
	if err != nil {
		fmt.Printf("Error while creating a file: %v", err)
		return
	}
	defer f.Close()
	for i := 0; i < len(predictions); i++ {
		for j := 0; j < len(predictions[i]); j++ {
			for k := 0; k < len(predictions[i][j]); k++ {
				if k == 0 {
					// Write a serial number of simulations at the beginning of a file.
					var serialNumberOfSimulation string = strconv.FormatInt(serialNumber, 10)
					row.WriteString(serialNumberOfSimulation + "," + " ")
				}

				//A component in one row of the output file.
				var component string = strconv.FormatFloat(predictions[i][j][k], 'f', 7, 64)
				row.WriteString(component)
				if k == (len(predictions[i][j]) - 1) {
					row.WriteString("\r\n")
					serialNumber++
				} else {
					row.WriteString("," + " ")
				}

			}
		}
	}
	// Writing to file
	_, err = f.WriteString(row.String())
	if err != nil {
		fmt.Printf("Error while writing a file: %v", err)
	}

}

func exportScalingFile(scaling string, strongFlag bool) {
	var path string
	if strongFlag == true {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Finance\\GolangPiStrongScaling.csv"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Finance\\GolangPiWeakScaling.csv"
	}
	f, err := os.Create(path) // creating...
	if err != nil {
		fmt.Printf("Error while creating a file: %v", err)
		return
	}
	defer f.Close()
	// Writing to file
	_, err = f.WriteString(scaling)
	if err != nil {
		fmt.Printf("Error while writing a file: %v", err)
	}

}

var serialPart float64 = 0.0
var parallelPart float64 = 1.0

func calculateAmdahlSpeedup(numberOfProcesses int) float64 {
	return 1.0 / (serialPart + parallelPart/float64(numberOfProcesses))
}

func calculateGustafsonSpeedup(numberOfProcesses int) float64 {
	return serialPart + parallelPart*float64(numberOfProcesses)
}

func strongScaling() {
	fmt.Println("=======================")
	fmt.Println("Start strong scaling:")
	fmt.Println("=======================")
	//One row in the output file.
	var row strings.Builder
	row.WriteString("number_of_processes,achieved_speedup,theoretical_maximum_speedup\r\n")
	numberOfSimulations := 100
	predictionWindowSize := 10
	numberOfProcessesSerial := 1
	monteCarloSimulationFinanceSerial := MonteCarloSimulationFinance{
		tickerSymbol:      "AAPL",
		startDate:         "2000-01-01",
		endDate:           "2020-01-01",
		numberOfProcesses: numberOfProcessesSerial}
	monteCarloSimulationFinanceSerial.dataAcquisition()
	monteCarloSimulationFinanceSerial.calculatePeriodicDailyReturn()
	fmt.Println("Stock market price predictions using the Monte Carlo simulation serial version")
	_, serialExecutionTime := monteCarloSimulationFinanceSerial.mcsFinanceSerial(numberOfSimulations, predictionWindowSize)
	fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
		numberOfSimulations, numberOfProcessesSerial,
		predictionWindowSize, serialExecutionTime)

	for numberOfProcessesParallel := 2; numberOfProcessesParallel < 14; numberOfProcessesParallel++ {

		monteCarloSimulationFinanceParallel := MonteCarloSimulationFinance{
			tickerSymbol:      "AAPL",
			startDate:         "2000-01-01",
			endDate:           "2020-01-01",
			numberOfProcesses: numberOfProcessesParallel}
		monteCarloSimulationFinanceParallel.dataAcquisition()
		monteCarloSimulationFinanceParallel.calculatePeriodicDailyReturn()
		fmt.Println("Stock market price predictions using the Monte Carlo simulation parallel version")
		_, parallelExecutionTime := monteCarloSimulationFinanceParallel.mcsFinanceParallel(numberOfSimulations, predictionWindowSize)

		fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
			numberOfSimulations, numberOfProcessesParallel,
			predictionWindowSize, parallelExecutionTime)

		achievedSpeedup := serialExecutionTime / parallelExecutionTime
		theoreticalMaximumSpeedup := calculateAmdahlSpeedup(numberOfProcessesParallel)
		fmt.Printf("Achieved speedup is: %f times.\r\n", achievedSpeedup)
		fmt.Printf("Theoretical maximum speedup according to Amdahl’s law is: %f times.\r\n", theoreticalMaximumSpeedup)
		row.WriteString(strconv.FormatInt(int64(numberOfProcessesParallel), 10) + "," + strconv.FormatFloat(achievedSpeedup, 'f', -1, 64) + "," + strconv.FormatInt(int64(theoreticalMaximumSpeedup), 10) + "\r\n")
	}
	exportScalingFile(row.String(), true)
	fmt.Println("End strong scaling.")

}
func weakScaling() {
	fmt.Println("=======================")
	fmt.Println("Start weak scaling:")
	fmt.Println("=======================")
	//One row in the output file.
	var row strings.Builder
	row.WriteString("number_of_processes,achieved_speedup,theoretical_maximum_speedup\r\n")
	numberOfSimulations := 100
	predictionWindowSize := 10
	for numberOfProcesses := 2; numberOfProcesses < 14; numberOfProcesses++ {
		increasedNumberOfSimulations := numberOfSimulations * numberOfProcesses
		monteCarloSimulationFinanceSerial := MonteCarloSimulationFinance{
			tickerSymbol:      "AAPL",
			startDate:         "2000-01-01",
			endDate:           "2020-01-01",
			numberOfProcesses: numberOfProcesses}
		monteCarloSimulationFinanceSerial.dataAcquisition()
		monteCarloSimulationFinanceSerial.calculatePeriodicDailyReturn()
		fmt.Println("Stock market price predictions using the Monte Carlo simulation serial version")
		_, serialExecutionTime := monteCarloSimulationFinanceSerial.mcsFinanceSerial(increasedNumberOfSimulations, predictionWindowSize)
		fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
			increasedNumberOfSimulations, numberOfProcesses,
			predictionWindowSize, serialExecutionTime)

		monteCarloSimulationFinanceParallel := MonteCarloSimulationFinance{
			tickerSymbol:      "AAPL",
			startDate:         "2000-01-01",
			endDate:           "2020-01-01",
			numberOfProcesses: numberOfProcesses}
		monteCarloSimulationFinanceParallel.dataAcquisition()
		monteCarloSimulationFinanceParallel.calculatePeriodicDailyReturn()
		fmt.Println("Stock market price predictions using the Monte Carlo simulation parallel version")
		_, parallelExecutionTime := monteCarloSimulationFinanceParallel.mcsFinanceParallel(increasedNumberOfSimulations, predictionWindowSize)

		fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
			increasedNumberOfSimulations, numberOfProcesses,
			predictionWindowSize, parallelExecutionTime)

		achievedSpeedup := serialExecutionTime / parallelExecutionTime
		theoreticalMaximumSpeedup := calculateAmdahlSpeedup(numberOfProcesses)
		fmt.Printf("Achieved speedup is: %f times.\r\n", achievedSpeedup)
		fmt.Printf("Theoretical maximum speedup according to Amdahl’s law is: %f times.\r\n", theoreticalMaximumSpeedup)
		row.WriteString(strconv.FormatInt(int64(numberOfProcesses), 10) + "," + strconv.FormatFloat(achievedSpeedup, 'f', -1, 64) + "," + strconv.FormatInt(int64(theoreticalMaximumSpeedup), 10) + "\r\n")
	}
	exportScalingFile(row.String(), false)
	fmt.Println("End weak scaling.")

}

func main() {
	/*
		numberOfSimulationsSerial := 10
		predictionWindowSizeSerial := 100
		numberOfProcessesSerial := 1
		monteCarloSimulationFinanceSerial := MonteCarloSimulationFinance{
			tickerSymbol:      "AAPL",
			startDate:         "2000-01-01",
			endDate:           "2020-01-01",
			numberOfProcesses: numberOfProcessesSerial}
		monteCarloSimulationFinanceSerial.dataAcquisition()
		monteCarloSimulationFinanceSerial.calculatePeriodicDailyReturn()
		fmt.Println("Stock market price predictions using the Monte Carlo simulation serial version")
		serialPredictions, serialExecutionTime := monteCarloSimulationFinanceSerial.mcsFinanceSerial(numberOfSimulationsSerial, predictionWindowSizeSerial)
		fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
			numberOfSimulationsSerial, numberOfProcessesSerial,
			predictionWindowSizeSerial, serialExecutionTime)
		monteCarloSimulationFinanceSerial.exportFinanceFileSerial(serialPredictions)
	*/
	/*
		numberOfSimulationsParallel := 10
		predictionWindowSizeParallel := 100
		numberOfProcessesParallel := 4
		monteCarloSimulationFinanceParallel := MonteCarloSimulationFinance{
			tickerSymbol:      "AAPL",
			startDate:         "2000-01-01",
			endDate:           "2020-01-01",
			numberOfProcesses: numberOfProcessesParallel}
		monteCarloSimulationFinanceParallel.dataAcquisition()
		monteCarloSimulationFinanceParallel.calculatePeriodicDailyReturn()
		fmt.Println("Stock market price predictions using the Monte Carlo simulation parallel version")
		parallelPredictions, parallelExecutionTime := monteCarloSimulationFinanceParallel.mcsFinanceParallel(numberOfSimulationsParallel, predictionWindowSizeParallel)

		fmt.Printf("Execution time(n = %d, p = %d, w = %d) = %f seconds\r\n",
			numberOfSimulationsParallel, numberOfProcessesParallel,
			predictionWindowSizeParallel, parallelExecutionTime)

		monteCarloSimulationFinanceParallel.exportFinanceFileParallel(parallelPredictions)
	*/
	strongScaling()
	weakScaling()

}
