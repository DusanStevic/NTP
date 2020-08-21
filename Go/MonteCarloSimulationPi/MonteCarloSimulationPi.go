package main

import (
	"fmt"
	"math/rand"
	"time"
)

// MonteCarloSimulationPi is a structure to calculate pi using Monte Carlo simulation.
type MonteCarloSimulationPi struct {
	numberOfProcesses int
}

func (monteCarloSimulationPi *MonteCarloSimulationPi) simulationPi(numberOfSimulations int) int {
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

	return inside
}

func main() {
	monteCarloSimulationPi := MonteCarloSimulationPi{numberOfProcesses: 5}
	fmt.Println(monteCarloSimulationPi.simulationPi(10000000))
	monteCarloSimulationPi.simulationPi(1000)
}
