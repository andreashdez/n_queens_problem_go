package main

import (
	"math"
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
			if selected == nil {
				t.Fatalf("selectRandomChromosome(%v) returned nil", tt.fitnessSum)
			}
			if selected != &ga.population[0] && selected != &ga.population[1] {
				t.Fatalf("selectRandomChromosome(%v) returned pointer outside population", tt.fitnessSum)
			}
		})
	}
}

func TestSelectRandomChromosomeReturnsPopulationPointer(t *testing.T) {
	ga := GeneticAlgorithm{
		population: []Chromosome{
			{positions: []int{0}, fitness: 10},
			{positions: []int{1}, fitness: 0},
		},
	}

	selected := ga.selectRandomChromosome(10)
	if selected != &ga.population[0] {
		t.Fatalf("selectRandomChromosome returned %p, want %p", selected, &ga.population[0])
	}
}
