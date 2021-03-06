package main

import (
	"freesolver"
	"log"
	"math/rand"
	"time"
)

func solveHardcodedCase() {
	g, err := freesolver.GenerateGameFromString(`
8C JS TH 8S AC 5D 9D
QD 7D 6S 5C 7S KC 8H
7C 6C TD QS JD AH 3H
2H 5H AD KS QC JC 8D
9S 6D TC 2C KD 9C
9H 3D TS JH 5S 4H
2D 3S 7H 2S 4C 6H
KH 3C AS 4D 4S QH
`)
	if err != nil {
		log.Fatalf("Can't parse: %s", err)
	}
	solver := freesolver.NewSolver(*g)
	solver.Debug = true
	won := solver.Solve()
	log.Printf("Won!! full game:\n%s", won.FullGameString())
}

func solveAllSeeds() {
	totalDuration := time.Duration(0)

	for s, count := int64(1), 1; ; s++ {
		rand.Seed(s)
		g := freesolver.GenerateGame()
		start := time.Now()
		solver := freesolver.NewSolver(g)
		solver.Solve()
		duration := time.Since(start)
		totalDuration += duration
		log.Printf("Seed %d won! average solve time: %v", s, totalDuration/time.Duration(count))
		if duration > time.Second {
			log.Printf("LONG SEED %d TOOK %s", s, duration)
		}
		count++
	}
}

func solveSpecificSeed(specificSeed int64) {
	rand.Seed(specificSeed)
	g := freesolver.GenerateGame()
	log.Printf("Solving:\n%s", g)
	start := time.Now()
	solver := freesolver.NewSolver(g)
	solver.Debug = true
	won := solver.Solve()
	duration := time.Since(start)
	log.Printf("Won!! full game:\n%s", won.FullGameString())
	log.Printf("Seed %d won! took: %v", specificSeed, duration)
}

func main() {
	solveAllSeeds()
	//solveHardcodedCase()
	//solveSpecificSeed(53)
}
