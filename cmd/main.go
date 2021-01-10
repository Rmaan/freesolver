package main

import (
	"freesolver"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(1)
	g := freesolver.GenerateGame()
	log.Printf("Solving:\n%s", g)
	start := time.Now()
	won := freesolver.NewSolver(g).Solve()
	duration := time.Since(start)
	log.Printf("Won!! full game:\n%s", won.FullGameString())
	log.Printf("took %s", duration)
}
