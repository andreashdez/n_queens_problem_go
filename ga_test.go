package main

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

func TestCalcFitnessAssignsUniformValuesWhenDiffIsZero(t *testing.T) {
	positions := []int{0, 1, 2, 3}
	ga := GeneticAlgorithm{
		population: []Chromosome{
			*NewChromosome(positions),
			*NewChromosome(positions),
		},
	}

	ga.calcFitness()

	for i, chromosome := range ga.population {
		if chromosome.fitness != 1 {
			t.Fatalf("population[%d].fitness = %v, want 1", i, chromosome.fitness)
		}
	}
}

func TestSelectRandomChromosomeFallsBackForInvalidFitnessSum(t *testing.T) {
	ga := GeneticAlgorithm{
		population: []Chromosome{
			{positions: []int{0}, fitness: 0},
			{positions: []int{1}, fitness: 0},
		},
	}

	tests := []struct {
		name       string
		fitnessSum float64
	}{
		{name: "zero", fitnessSum: 0},
		{name: "nan", fitnessSum: math.NaN()},
		{name: "positive inf", fitnessSum: math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected := ga.selectRandomChromosome(tt.fitnessSum)
			if len(selected.positions) == 0 {
				t.Fatalf("selectRandomChromosome(%v) returned empty chromosome", tt.fitnessSum)
			}
			if selected.positions[0] != 0 && selected.positions[0] != 1 {
				t.Fatalf("selectRandomChromosome(%v) returned value outside population", tt.fitnessSum)
			}
		})
	}
}

func TestSelectRandomChromosomeReturnsPopulationValue(t *testing.T) {
	ga := GeneticAlgorithm{
		population: []Chromosome{
			{positions: []int{0}, fitness: 10},
			{positions: []int{1}, fitness: 0},
		},
	}

	selected := ga.selectRandomChromosome(10)
	if len(selected.positions) == 0 || selected.positions[0] != 0 {
		t.Fatalf("selectRandomChromosome returned %v, want first chromosome", selected.positions)
	}
}

func TestMutateGenesRateZeroDoesNotMutate(t *testing.T) {
	ga := GeneticAlgorithm{mutationRate: 0}
	genes := []int{0, 1, 2, 3}
	before := append([]int(nil), genes...)

	mutated := ga.mutateGenes(genes)

	if mutated {
		t.Fatalf("mutateGenes returned true, want false")
	}
	if !slices.Equal(genes, before) {
		t.Fatalf("genes changed from %v to %v with zero mutation rate", before, genes)
	}
}

func TestMutateGenesRateOneSwapsTwoValues(t *testing.T) {
	ga := GeneticAlgorithm{mutationRate: 1}
	genes := []int{0, 1, 2, 3, 4, 5}
	before := append([]int(nil), genes...)
	rand.Seed(7)

	mutated := ga.mutateGenes(genes)

	if !mutated {
		t.Fatalf("mutateGenes returned false, want true")
	}
	if slices.Equal(genes, before) {
		t.Fatalf("genes did not change after mutation: %v", genes)
	}
	slices.Sort(genes)
	slices.Sort(before)
	if !slices.Equal(genes, before) {
		t.Fatalf("mutation changed gene set, got %v want permutation of %v", genes, before)
	}
}

func TestPMXReturnsPermutation(t *testing.T) {
	ga := GeneticAlgorithm{}
	parentOne := []int{0, 1, 2, 3, 4, 5, 6, 7}
	parentTwo := []int{7, 6, 5, 4, 3, 2, 1, 0}
	want := append([]int(nil), parentOne...)
	slices.Sort(want)

	for i := 0; i < 100; i++ {
		rand.Seed(int64(i + 1))
		child := ga.pmx(parentOne, parentTwo)
		if len(child) != len(parentOne) {
			t.Fatalf("pmx child len = %d, want %d", len(child), len(parentOne))
		}
		got := append([]int(nil), child...)
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Fatalf("pmx child is not a permutation: %v", child)
		}
	}
}
