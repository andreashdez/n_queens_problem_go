package main

import (
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"
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
	ga := GeneticAlgorithm{mutationRate: 1, rng: rand.New(rand.NewSource(7))}
	genes := []int{0, 1, 2, 3, 4, 5}
	before := append([]int(nil), genes...)

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
		ga.rng = rand.New(rand.NewSource(int64(i + 1)))
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

func TestRunAlgorithmDeterministicWithSameSeed(t *testing.T) {
	seed := int64(1234)
	gaOne := BuildGeneticAlgorithm(8, 128, 150, 4, 8, 0.05, rand.New(rand.NewSource(seed)))
	gaTwo := BuildGeneticAlgorithm(8, 128, 150, 4, 8, 0.05, rand.New(rand.NewSource(seed)))

	bestOne := gaOne.RunAlgorithm()
	bestTwo := gaTwo.RunAlgorithm()

	if bestOne.conflictsSum != bestTwo.conflictsSum {
		t.Fatalf("best conflicts differ for same seed: %d vs %d", bestOne.conflictsSum, bestTwo.conflictsSum)
	}
	if !slices.Equal(bestOne.positions, bestTwo.positions) {
		t.Fatalf("best positions differ for same seed: %v vs %v", bestOne.positions, bestTwo.positions)
	}
	if !slices.Equal(bestOne.conflicts, bestTwo.conflicts) {
		t.Fatalf("best conflicts vector differ for same seed: %v vs %v", bestOne.conflicts, bestTwo.conflicts)
	}
}

func TestBuildGeneticAlgorithmCreatesPermutationPopulation(t *testing.T) {
	size := 12
	populationSize := 40
	ga := BuildGeneticAlgorithm(size, populationSize, 50, 2, 5, 0.05, rand.New(rand.NewSource(99)))

	if len(ga.population) != populationSize {
		t.Fatalf("population size = %d, want %d", len(ga.population), populationSize)
	}

	for i, chromosome := range ga.population {
		if len(chromosome.positions) != size {
			t.Fatalf("population[%d] chromosome size = %d, want %d", i, len(chromosome.positions), size)
		}
		if !isPermutationOfRange(chromosome.positions, size) {
			t.Fatalf("population[%d] is not a permutation: %v", i, chromosome.positions)
		}
	}
}

func TestMutateGenesPreservesPermutationAcrossSeeds(t *testing.T) {
	base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for seed := 1; seed <= 50; seed++ {
		ga := GeneticAlgorithm{mutationRate: 1, rng: rand.New(rand.NewSource(int64(seed)))}
		genes := append([]int(nil), base...)

		mutated := ga.mutateGenes(genes)
		if !mutated {
			t.Fatalf("seed %d: mutateGenes returned false, want true", seed)
		}
		if !isPermutationOfRange(genes, len(base)) {
			t.Fatalf("seed %d: mutation produced non-permutation: %v", seed, genes)
		}

		differentPositions := 0
		for i := range genes {
			if genes[i] != base[i] {
				differentPositions++
			}
		}
		if differentPositions != 2 {
			t.Fatalf("seed %d: mutation changed %d positions, want 2", seed, differentPositions)
		}
	}
}

func TestRunAlgorithmStopsAtMaxEpochsWithoutMating(t *testing.T) {
	ga := BuildGeneticAlgorithm(3, 32, 1, 0, 0, 0, rand.New(rand.NewSource(77)))
	initialBest := ga.getBestChromosome()

	done := make(chan Chromosome, 1)
	go func() {
		done <- ga.RunAlgorithm()
	}()

	select {
	case result := <-done:
		if result.conflictsSum == 0 {
			t.Fatalf("RunAlgorithm returned solved chromosome for size 3: %v", result.positions)
		}
		if result.conflictsSum != initialBest.conflictsSum {
			t.Fatalf("best conflicts changed without mating: got %d want %d", result.conflictsSum, initialBest.conflictsSum)
		}
		if !slices.Equal(result.positions, initialBest.positions) {
			t.Fatalf("best positions changed without mating: got %v want %v", result.positions, initialBest.positions)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("RunAlgorithm did not terminate at max epochs")
	}
}

func TestRunAlgorithmKeepsPopulationBounded(t *testing.T) {
	targetPopulation := 60
	ga := BuildGeneticAlgorithm(8, targetPopulation, 20, 8, 12, 0.05, rand.New(rand.NewSource(101)))

	_ = ga.RunAlgorithm()

	if len(ga.population) != targetPopulation {
		t.Fatalf("population size after RunAlgorithm = %d, want %d", len(ga.population), targetPopulation)
	}
}

func isPermutationOfRange(values []int, size int) bool {
	if len(values) != size {
		return false
	}
	seen := make([]bool, size)
	for _, value := range values {
		if value < 0 || value >= size || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}
