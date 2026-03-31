package main

import (
	"math/rand"

	"github.com/rs/zerolog/log"
)

type Chromosome struct {
	positions    []int
	conflicts    []int
	conflictsSum int
	fitness      float64
}

func NewChromosome(positions []int) *Chromosome {
	conflicts := countConflicts(positions)
	conflictsSum := sumConflicts(conflicts)
	fitness := 0.0
	log.
		Debug().
		Int("conflictsSum", conflictsSum).
		Msg("chromosome conflicts sum")
	return &Chromosome{positions, conflicts, conflictsSum, fitness}
}

func (c *Chromosome) SetFitness(fitness float64) {
	c.fitness = fitness
}

func GenerateDistinctRandomValues(rng *rand.Rand, size int) []int {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	return rng.Perm(size)
}

func countConflicts(positions []int) []int {
	size := len(positions)
	conflicts := make([]int, size)
	diagonalCount := size*2 - 1
	mainDiagonalCounts := make([]int, diagonalCount)
	antiDiagonalCounts := make([]int, diagonalCount)
	diagonalOffset := size - 1

	for x, y := range positions {
		mainDiagonalCounts[x-y+diagonalOffset]++
		antiDiagonalCounts[x+y]++
	}

	for x, y := range positions {
		mainDiagonalConflicts := mainDiagonalCounts[x-y+diagonalOffset] - 1
		antiDiagonalConflicts := antiDiagonalCounts[x+y] - 1
		conflicts[x] = mainDiagonalConflicts + antiDiagonalConflicts
		if conflicts[x] > 0 {
			log.
				Trace().
				Int("x", x).
				Int("y", y).
				Int("conflicts", conflicts[x]).
				Msg("found conflicts")
		}
	}
	return conflicts
}

func sumConflicts(conflicts []int) int {
	conflictsSum := 0
	for _, c := range conflicts {
		conflictsSum += c
	}
	return conflictsSum / 2
}
