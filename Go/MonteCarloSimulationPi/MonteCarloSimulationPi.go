package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// mcs stands for Monte Carlo Simulation
// pi=3.1415926535

// MonteCarloSimulationPi is a structure to calculate pi using Monte Carlo simulation.
type MonteCarloSimulationPi struct {
	numberOfProcesses int
	parallelFlag      bool
	experimentFlag    bool
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) simulationPi(numberOfSimulations int, channel chan int) {
	if monteCarloSimulationPi.experimentFlag == true {
		inside := 0
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		for i := 0; i < numberOfSimulations; i++ {
			x := r.Float64()
			y := r.Float64()
			if (x*x + y*y) < 1 {
				inside++
			}
		}
		channel <- inside
	} else {
		//One row in the output file.
		var row strings.Builder
		inside := 0
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
			x := r.Float64()
			y := r.Float64()

			// Pharo for Data Visualization. Circle of radius 250 centered at the point(250, 250).
			// To create a Rectangle in Pharo you must provide the top left and the bottom right points.
			row.WriteString(strconv.FormatInt(int64(x*500), 10) + " " + strconv.FormatInt(int64(y*500), 10) + "\r\n")
			if (x*x + y*y) < 1 {
				inside++
			}
		}
		monteCarloSimulationPi.exportPiFile(row.String())
		channel <- inside
	}

}

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiSerial(numberOfSimulations int) (float64, float64) {
	startTime := time.Now()
	monteCarloSimulationPi.parallelFlag = false
	channel := make(chan int)
	go monteCarloSimulationPi.simulationPi(numberOfSimulations, channel)
	inside := <-channel
	pi := 4 * float64(inside) / float64(numberOfSimulations)
	executionTime := time.Since(startTime).Seconds()
	return pi, executionTime
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiParallel(numberOfSimulations int) (float64, float64) {
	startTime := time.Now()
	monteCarloSimulationPi.parallelFlag = true
	numberOfSimulationsPerProcess := numberOfSimulations / monteCarloSimulationPi.numberOfProcesses
	/* 	Buffered channels are useful when you know how many goroutines you have launched,
	   	want to limit the number of goroutines you will launch, or want to limit
	   	the amount of work that is queued up. */
	channel := make(chan int, monteCarloSimulationPi.numberOfProcesses)

	for i := 0; i < monteCarloSimulationPi.numberOfProcesses; i++ {
		go monteCarloSimulationPi.simulationPi(numberOfSimulationsPerProcess, channel)
	}

	var inside int
	for i := 0; i < monteCarloSimulationPi.numberOfProcesses; i++ {
		inside += <-channel
	}
	pi := 4 * float64(inside) / float64(numberOfSimulations)
	executionTime := time.Since(startTime).Seconds()
	return pi, executionTime

}

func (monteCarloSimulationPi *MonteCarloSimulationPi) exportPiFile(simulations string) {
	var path string
	if monteCarloSimulationPi.parallelFlag == false {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Pi\\GolangPiSerial.txt"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Execution Results\\Pi\\GolangPiParallel.txt"
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
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Pi\\GolangPiStrongScaling.csv"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Scaling Results\\Pi\\GolangPiWeakScaling.csv"
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
	numberOfSimulations := 100000000
	numberOfProcessesSerial := 1
	monteCarloSimulationPiSerial := MonteCarloSimulationPi{numberOfProcesses: numberOfProcessesSerial}
	monteCarloSimulationPiSerial.experimentFlag = true
	fmt.Println("Approximation of Pi by using the Monte Carlo simulation serial version")
	serialPi, serialExecutionTime := monteCarloSimulationPiSerial.mcsPiSerial(numberOfSimulations)
	fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", numberOfSimulations, numberOfProcessesSerial, serialPi)
	fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime)
	for numberOfProcessesParallel := 2; numberOfProcessesParallel < 14; numberOfProcessesParallel++ {
		monteCarloSimulationPiParallel := MonteCarloSimulationPi{numberOfProcesses: numberOfProcessesParallel}
		monteCarloSimulationPiParallel.experimentFlag = true
		fmt.Println("Approximation of Pi by using the Monte Carlo simulation parallel version")
		parallelPi, parallelExecutionTime := monteCarloSimulationPiParallel.mcsPiParallel(numberOfSimulations)
		fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", numberOfSimulations, numberOfProcessesParallel, parallelPi)
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
	numberOfSimulations := 100000000
	for numberOfProcesses := 2; numberOfProcesses < 14; numberOfProcesses++ {
		increasedNumberOfSimulations := numberOfSimulations * numberOfProcesses
		monteCarloSimulationPi := MonteCarloSimulationPi{numberOfProcesses: numberOfProcesses}
		monteCarloSimulationPi.experimentFlag = true
		fmt.Println("Approximation of Pi by using the Monte Carlo simulation serial version")
		serialPi, serialExecutionTime := monteCarloSimulationPi.mcsPiSerial(increasedNumberOfSimulations)
		fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", increasedNumberOfSimulations, 1, serialPi)
		fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime)
		fmt.Println("Approximation of Pi by using the Monte Carlo simulation parallel version")
		parallelPi, parallelExecutionTime := monteCarloSimulationPi.mcsPiParallel(increasedNumberOfSimulations)
		fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", increasedNumberOfSimulations, numberOfProcesses, parallelPi)
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
	numberOfSimulationsSerial := 1000
	numberOfProcessesSerial := 1
	monteCarloSimulationPiSerial := MonteCarloSimulationPi{numberOfProcesses: numberOfProcessesSerial}
	monteCarloSimulationPiSerial.experimentFlag = false
	fmt.Println("Approximation of Pi by using the Monte Carlo simulation serial version")
	serialPi, serialExecutionTime := monteCarloSimulationPiSerial.mcsPiSerial(numberOfSimulationsSerial)
	fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", numberOfSimulationsSerial, numberOfProcessesSerial, serialPi)
	fmt.Printf("Execution time (duration): %f seconds\r\n", serialExecutionTime)

	numberOfSimulationsParallel := 4000
	numberOfProcessesParallel := 4
	monteCarloSimulationPiParallel := MonteCarloSimulationPi{numberOfProcesses: numberOfProcessesParallel}
	monteCarloSimulationPiParallel.experimentFlag = false
	fmt.Println("Approximation of Pi by using the Monte Carlo simulation parallel version")
	parallelPi, parallelExecutionTime := monteCarloSimulationPiParallel.mcsPiParallel(numberOfSimulationsParallel)
	fmt.Printf("Pi(n = %d, p = %d) = %f\r\n", numberOfSimulationsParallel, numberOfProcessesParallel, parallelPi)
	fmt.Printf("Execution time (duration): %f seconds", parallelExecutionTime)

	//strongScaling()
	//weakScaling()

}
