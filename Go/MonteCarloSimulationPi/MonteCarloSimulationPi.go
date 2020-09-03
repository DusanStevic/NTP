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

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiSerial(numberOfSimulations int) (float64, func() float64) {
	executionTime := calculateExecutionTime()
	defer executionTime()
	monteCarloSimulationPi.parallelFlag = false
	channel := make(chan int)
	go monteCarloSimulationPi.simulationPi(numberOfSimulations, channel)
	inside := <-channel
	pi := 4 * float64(inside) / float64(numberOfSimulations)
	return pi, executionTime
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiParallel(numberOfSimulations int) float64 {
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

	return 4 * float64(inside) / float64(numberOfSimulations)
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

func calculateExecutionTime() func() float64 {
	startTime := time.Now()
	//These kind of functions are called anonymous functions since they do not have a name.
	return func() float64 {
		return time.Since(startTime).Seconds()
	}
}

func main() {
	monteCarloSimulationPi := MonteCarloSimulationPi{numberOfProcesses: 4}
	monteCarloSimulationPi.experimentFlag = true
	//start := time.Now()
	fmt.Println(monteCarloSimulationPi.mcsPiSerial(100000))
	//fmt.Println(monteCarloSimulationPi.mcsPiParallel(100000))
	//duration := time.Since(start)
	//fmt.Println("duration:", duration)
	serialPi, serialExecutionTime := monteCarloSimulationPi.mcsPiSerial(100000)
	fmt.Println("pi:", serialPi)
	fmt.Println("duration:", serialExecutionTime())

}
