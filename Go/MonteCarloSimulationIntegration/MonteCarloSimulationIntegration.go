package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Lower bound of Integral.
	lowerBound = 1.0
	// Upper bound of Integral.
	upperBound = 2.0
	// The area under the graph of a function can be found by adding slices that approach zero in width.
	sliceSize = 0.01
)

// mcs stands for Monte Carlo Simulation

// MonteCarloSimulationIntegration is a structure to calculate integral using Monte Carlo simulation.
type MonteCarloSimulationIntegration struct {
	numberOfProcesses int
	parallelFlag      bool
	experimentFlag    bool
}

/* The function f(x) to be integrated is called the integrand.
The function we are integrating must be non-negative continuous function between lower bound and upper bound
Non-negative function: is a function when it attain non negative values only. A function would be called a
positive function if its values are positive for all arguments of its domain, or a non-negative function
if all of its values are non-negative. The function graph sits above or on the x-axis.
Continuous function: is a function with no holes, jumps or vertical asymptotes
(where the function heads up/down towards infinity). A vertical asymptote between lower bound and
upper bound affects the definite integral. */
func (MonteCarloSimulationIntegration *MonteCarloSimulationIntegration) function(x float64) float64 {
	return 2 * x
}

func (MonteCarloSimulationIntegration *MonteCarloSimulationIntegration) simulationIntegration(numberOfSimulations int, channel chan float64) {
	if MonteCarloSimulationIntegration.experimentFlag == true {
		// Points under the graph of a function.
		below := 0
		lowerBoundInterval := lowerBound
		upperBoundInterval := upperBound
		// Define the interval between the lower and upper bound.
		var x []float64
		// Function Values
		var y []float64
		// Maximum of the function f(x) on the interval[lower_bound, upper_bound]
		fMax := MonteCarloSimulationIntegration.function(lowerBound)
		for lowerBoundInterval < upperBoundInterval {
			x = append(x, lowerBoundInterval)
			t := MonteCarloSimulationIntegration.function(lowerBoundInterval)
			y = append(y, t)
			if t > fMax {
				fMax = t
			}
			lowerBoundInterval += sliceSize

		}

		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		for i := 0; i < numberOfSimulations; i++ {
			xRand := lowerBound + (upperBound-lowerBound)*r.Float64()
			yRand := 0 + fMax*r.Float64()
			if yRand < MonteCarloSimulationIntegration.function(xRand) {
				below++
			}
		}
		// Rectangle area that surrounds the area under the graph of a function.
		a := upperBound - lowerBound
		b := fMax - 0
		rectangleArea := a * b
		// bellow = Points under the graph of a function.
		// number_of_simulations = Total number of points = Points inside rectangle
		proportion := float64(below) / float64(numberOfSimulations)
		integral := float64(proportion) * float64(rectangleArea)
		channel <- integral
	} else {
		//One row in the output file.
		var row strings.Builder
		// Points under the graph of a function.
		below := 0
		lowerBoundInterval := lowerBound
		upperBoundInterval := upperBound
		// Define the interval between the lower and upper bound.
		var x []float64
		// Function Values
		var y []float64
		// Maximum of the function f(x) on the interval[lower_bound, upper_bound]
		fMax := MonteCarloSimulationIntegration.function(lowerBound)
		for lowerBoundInterval < upperBoundInterval {
			x = append(x, lowerBoundInterval)
			t := MonteCarloSimulationIntegration.function(lowerBoundInterval)
			y = append(y, t)
			if t > fMax {
				fMax = t
			}
			lowerBoundInterval += sliceSize

		}
		// Source code for random number generator https://play.golang.org/p/ZdFpbahgC1
		// The default number generator is deterministic, so it'll
		// produce the same sequence of numbers each time by default.
		// To produce varying sequences, give it a seed that changes.
		// Note that this is not safe to use for random numbers you
		// intend to be secret, use `crypto/rand` for those.
		// Seeding - Go provides a method, Seed(seed int64), that allows you
		// to initialize this default sequence. Implementation is slow
		// to make it faster rand.Seed(time.Now().UnixNano()) is added.
		// Seed is the current time, converted to int64 by UnixNano.
		// Gives constantly changing numbers.

		// Seed
		s := rand.NewSource(time.Now().UnixNano())
		// Randomly changing numbers.
		r := rand.New(s)
		for i := 0; i < numberOfSimulations; i++ {
			// Call the resulting `rand.Rand` just like the
			// functions on the `rand` package.
			xRand := lowerBound + (upperBound-lowerBound)*r.Float64()
			yRand := 0 + fMax*r.Float64()
			xRandString := fmt.Sprintf("%.2f", xRand)
			yRandString := fmt.Sprintf("%.2f", yRand)
			row.WriteString(xRandString + " " + yRandString + "\r\n")
			if yRand < MonteCarloSimulationIntegration.function(xRand) {
				below++
			}
		}
		// Rectangle area that surrounds the area under the graph of a function.
		a := upperBound - lowerBound
		b := fMax - 0
		rectangleArea := a * b
		// bellow = Points under the graph of a function.
		// number_of_simulations = Total number of points = Points inside rectangle
		proportion := float64(below) / float64(numberOfSimulations)
		integral := float64(proportion) * float64(rectangleArea)
		MonteCarloSimulationIntegration.exportIntegrationFile(row.String())
		channel <- integral

	}

}
func (MonteCarloSimulationIntegration *MonteCarloSimulationIntegration) mcsIntegrationSerial(numberOfSimulations int) (float64, float64) {
	startTime := time.Now()
	MonteCarloSimulationIntegration.parallelFlag = false
	channel := make(chan float64)
	go MonteCarloSimulationIntegration.simulationIntegration(numberOfSimulations, channel)
	integral := <-channel
	executionTime := time.Since(startTime).Seconds()
	return integral, executionTime
}

func (MonteCarloSimulationIntegration *MonteCarloSimulationIntegration) mcsIntegrationParallel(numberOfSimulations int) (float64, float64) {
	startTime := time.Now()
	MonteCarloSimulationIntegration.parallelFlag = true
	numberOfSimulationsPerProcess := numberOfSimulations / MonteCarloSimulationIntegration.numberOfProcesses
	/* 	Buffered channels are useful when you know how many goroutines you have launched,
	   	want to limit the number of goroutines you will launch, or want to limit
	   	the amount of work that is queued up. */
	channel := make(chan float64, MonteCarloSimulationIntegration.numberOfProcesses)
	// partial result per process
	for i := 0; i < MonteCarloSimulationIntegration.numberOfProcesses; i++ {
		go MonteCarloSimulationIntegration.simulationIntegration(numberOfSimulationsPerProcess, channel)
	}

	var integralPerProcesses float64
	// cumulative result, aggregating partial results
	for i := 0; i < MonteCarloSimulationIntegration.numberOfProcesses; i++ {
		integralPerProcesses += <-channel
	}
	integral := float64(integralPerProcesses) / float64(MonteCarloSimulationIntegration.numberOfProcesses)
	executionTime := time.Since(startTime).Seconds()
	return integral, executionTime

}
func (MonteCarloSimulationIntegration *MonteCarloSimulationIntegration) exportIntegrationFile(simulations string) {
	var path string
	if MonteCarloSimulationIntegration.parallelFlag == false {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Integration\\GolangIntegrationSerial.txt"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Integration\\GolangIntegrationParallel.txt"
	}
	f, err := os.Create(path) // creating...
	if err != nil {
		fmt.Printf("Error while creating a file: %v", err)
		return
	}
	defer f.Close()
	// Writing to file
	_, err = f.WriteString(simulations)
	if err != nil {
		fmt.Printf("Error while writing a file: %v", err)
	}

}
func exportScalingFile(scaling string, strongFlag bool) {
	var path string
	if strongFlag == true {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Integration\\GolangIntegrationStrongScaling.csv"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Integration\\GolangIntegrationWeakScaling.csv"
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
	numberOfSimulations := 10000000
	numberOfProcessesSerial := 1
	monteCarloSimulationIntegrationSerial := MonteCarloSimulationIntegration{numberOfProcesses: numberOfProcessesSerial}
	monteCarloSimulationIntegrationSerial.experimentFlag = true
	fmt.Println("Integral approximation by using the Monte Carlo simulation serial version")
	serialIntegration, serialExecutionTime := monteCarloSimulationIntegrationSerial.mcsIntegrationSerial(numberOfSimulations)
	fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", numberOfSimulations, numberOfProcessesSerial, serialIntegration)
	fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime)
	for numberOfProcessesParallel := 2; numberOfProcessesParallel < 14; numberOfProcessesParallel++ {
		monteCarloSimulationIntegrationParallel := MonteCarloSimulationIntegration{numberOfProcesses: numberOfProcessesParallel}
		monteCarloSimulationIntegrationParallel.experimentFlag = true
		fmt.Println("Integral approximation by using the Monte Carlo simulation parallel version")
		parallelIntegration, parallelExecutionTime := monteCarloSimulationIntegrationParallel.mcsIntegrationParallel(numberOfSimulations)
		fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", numberOfSimulations, numberOfProcessesParallel, parallelIntegration)
		fmt.Printf("Execution time (duration): %f seconds.\r\n", parallelExecutionTime)
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
	numberOfSimulations := 10000000
	for numberOfProcesses := 2; numberOfProcesses < 14; numberOfProcesses++ {
		increasedNumberOfSimulations := numberOfSimulations * numberOfProcesses
		monteCarloSimulationIntegration := MonteCarloSimulationIntegration{numberOfProcesses: numberOfProcesses}
		monteCarloSimulationIntegration.experimentFlag = true
		fmt.Println("Integral approximation by using the Monte Carlo simulation serial version")
		serialIntegration, serialExecutionTime := monteCarloSimulationIntegration.mcsIntegrationSerial(increasedNumberOfSimulations)
		fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", increasedNumberOfSimulations, 1, serialIntegration)
		fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime)
		fmt.Println("Integral approximation by using the Monte Carlo simulation parallel version")
		parallelIntegration, parallelExecutionTime := monteCarloSimulationIntegration.mcsIntegrationParallel(increasedNumberOfSimulations)
		fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", increasedNumberOfSimulations, numberOfProcesses, parallelIntegration)
		fmt.Printf("Execution time (duration): %f seconds.\r\n", parallelExecutionTime)
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
	/* 	numberOfSimulationsSerial := 1000
	   	numberOfProcessesSerial := 1
	   	monteCarloSimulationIntegrationSerial := MonteCarloSimulationIntegration{numberOfProcesses: numberOfProcessesSerial}
	   	monteCarloSimulationIntegrationSerial.experimentFlag = false
	   	fmt.Println("Integral Approximation by using the Monte Carlo simulation serial version")
	   	serialIntegration, serialExecutionTime := monteCarloSimulationIntegrationSerial.mcsIntegrationSerial(numberOfSimulationsSerial)
	   	fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", numberOfSimulationsSerial, numberOfProcessesSerial, serialIntegration)
	   	fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime) */

	/* 	numberOfSimulationsParallel := 1000
	   	numberOfProcessesParallel := 4
	   	monteCarloSimulationIntegrationParallel := MonteCarloSimulationIntegration{numberOfProcesses: numberOfProcessesParallel}
	   	monteCarloSimulationIntegrationParallel.experimentFlag = false
	   	fmt.Println("Integral Approximation by using the Monte Carlo simulation parallel version")
	   	parallelIntegration, parallelExecutionTime := monteCarloSimulationIntegrationParallel.mcsIntegrationParallel(numberOfSimulationsParallel)
	   	fmt.Printf("Integral(n = %d, p = %d) = %f\r\n", numberOfSimulationsParallel, numberOfProcessesParallel, parallelIntegration)
	   	fmt.Printf("Execution time (duration): %f seconds", parallelExecutionTime) */

	strongScaling()
	weakScaling()
}
