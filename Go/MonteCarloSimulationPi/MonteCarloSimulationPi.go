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
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) simulationPi(numberOfSimulations int, channel chan int) {
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

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiSerial(numberOfSimulations int) float64 {
	monteCarloSimulationPi.parallelFlag = false
	channel := make(chan int)
	go monteCarloSimulationPi.simulationPi(numberOfSimulations, channel)
	inside := <-channel
	return 4 * float64(inside) / float64(numberOfSimulations)
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
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Pharo\\GolangPiSerial.txt"
	} else {
		path = "C:\\Users\\Dule\\Desktop\\NAPREDNE TEHNIKE PROGRAMIRANJA\\PROJEKAT\\NTP\\Pharo\\GolangPiParallel.txt"
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

func main() {
	monteCarloSimulationPi := MonteCarloSimulationPi{numberOfProcesses: 4}
	start := time.Now()
	//fmt.Println(monteCarloSimulationPi.mcsPiSerial(1000000000))
	fmt.Println(monteCarloSimulationPi.mcsPiParallel(10000))
	duration := time.Since(start)
	fmt.Println("duration:", duration)

}
