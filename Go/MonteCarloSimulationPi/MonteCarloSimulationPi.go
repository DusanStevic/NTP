package main

import (
	"fmt"
	"math/rand"
	"time"
)

// mcs stands for Monte Carlo Simulation
// pi=3.1415926535

// MonteCarloSimulationPi is a structure to calculate pi using Monte Carlo simulation.
type MonteCarloSimulationPi struct {
	numberOfProcesses int
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) simulationPi(numberOfSimulations int, channel chan int) {
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
		if (x*x + y*y) < 1 {
			inside++
		}
	}

	channel <- inside
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiSerial(numberOfSimulations int) float64 {
	channel := make(chan int)
	go monteCarloSimulationPi.simulationPi(numberOfSimulations, channel)
	inside := <-channel
	return 4 * float64(inside) / float64(numberOfSimulations)
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) mcsPiParallel(numberOfSimulations int) float64 {

	numberOfSimulationsPerProcess := numberOfSimulations / monteCarloSimulationPi.numberOfProcesses
	channel := make(chan int, monteCarloSimulationPi.numberOfProcesses)

	for j := 0; j < monteCarloSimulationPi.numberOfProcesses; j++ {
		go monteCarloSimulationPi.simulationPi(numberOfSimulationsPerProcess, channel)
	}

	var inside int
	for i := 0; i < monteCarloSimulationPi.numberOfProcesses; i++ {
		inside += <-channel
	}

	return 4 * float64(inside) / float64(numberOfSimulations)
}

func main() {
	monteCarloSimulationPi := MonteCarloSimulationPi{numberOfProcesses: 4}
	start := time.Now()
	fmt.Println(monteCarloSimulationPi.mcsPiSerial(1000000000))
	//fmt.Println(monteCarloSimulationPi.mcsPiParallel(1000000000))
	duration := time.Since(start)
	fmt.Println("duration:", duration)

}
