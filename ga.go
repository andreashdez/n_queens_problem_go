package main

import (
	"math"
	"math/rand"
	"slices"

	"github.com/rs/zerolog/log"
)

type GeneticAlgorithm struct {
	population []Chromosome
}

func (ga GeneticAlgorithm) getBestChromosome() Chromosome {
	bestChromosome := ga.population[0]
	for _, chromosome := range ga.population {
		if chromosome.conflictsSum < bestChromosome.conflictsSum {
			bestChromosome = chromosome
		}
	}
	return bestChromosome
}

func (ga GeneticAlgorithm) getWorstChromosome() Chromosome {
	bestChromosome := ga.population[0]
	for _, chromosome := range ga.population {
		if chromosome.conflictsSum > bestChromosome.conflictsSum {
			bestChromosome = chromosome
		}
	}
	return bestChromosome
}

func (ga GeneticAlgorithm) calcFitness() {
	mostConflicts := float64(ga.getWorstChromosome().conflictsSum)
	leastConflicts := float64(ga.getBestChromosome().conflictsSum)
	diffConflicts := mostConflicts - leastConflicts
	if diffConflicts == 0 {
		log.Debug().Msg("all chromosomes have equal conflicts; assigning uniform fitness")
		for i := range ga.population {
			ga.population[i].fitness = 1
		}
		return
	}
	log.
		Debug().
		Float64("mostConflicts", mostConflicts).
		Float64("leastConflicts", leastConflicts).
		Float64("diffConflicts", diffConflicts).
		Msg("calculating fitness")
	for i, c := range ga.population {
		conflictsSum := float64(c.conflictsSum)
		fitness := math.Pow(mostConflicts-conflictsSum, 3.0) / math.Pow(diffConflicts, 3.0)
		if math.IsNaN(fitness) || math.IsInf(fitness, 0) {
			fitness = 0
		}
		ga.population[i].fitness = fitness
		log.
			Trace().
			Float64("conflictsSum", conflictsSum).
			Float64("fitness", fitness).
			Msg("calculating fitness for chromosome")
	}
}

func (ga *GeneticAlgorithm) mateRandomChromosomes(minToMate int, maxToMate int) {
	if len(ga.population) == 0 {
		return
	}
	mateAmount := rand.Intn(maxToMate-minToMate) + minToMate
	fitnessSum := 0.0
	for _, v := range ga.population {
		fitnessSum += v.fitness
	}
	log.
		Debug().
		Int("mateAmount", mateAmount).
		Float64("fitnessSum", fitnessSum).
		Msg("mate random chromosomes")
	for range mateAmount {
		parentOne := ga.selectRandomChromosome(fitnessSum)
		parentTwo := ga.selectRandomChromosome(fitnessSum)
		child := ga.mateChromosomes(*parentOne, *parentTwo)
		ga.population = append(ga.population, *child)
	}
}

func (ga GeneticAlgorithm) selectRandomChromosome(fitnessSum float64) *Chromosome {
	if len(ga.population) == 0 {
		return nil
	}
	if fitnessSum <= 0 || math.IsNaN(fitnessSum) || math.IsInf(fitnessSum, 0) {
		randomIndex := rand.Intn(len(ga.population))
		return &ga.population[randomIndex]
	}
	rouletteSpin := rand.Float64() * fitnessSum
	selectionRank := 0.0
	for i := range ga.population {
		selectionRank += ga.population[i].fitness
		if selectionRank > rouletteSpin {
			return &ga.population[i]
		}
	}
	return &ga.population[0]
}

func (ga GeneticAlgorithm) mateChromosomes(parentOne Chromosome, parentTwo Chromosome) *Chromosome {
	log.
		Trace().
		Ints("parentOne", parentOne.positions).
		Ints("parentTwo", parentTwo.positions).
		Msg("mate random chromosomes")
	childGenes := ga.pmx(parentOne.positions, parentTwo.positions)
	child := NewChromosome(childGenes)
	return child
}

func (ga GeneticAlgorithm) pmx(parentOne []int, parentTwo []int) []int {
	chromosomeSize := len(parentOne)
	chromosomeHalfSize := chromosomeSize / 2
	pointOne := rand.Intn(chromosomeHalfSize)
	pointTwo := rand.Intn(chromosomeSize-chromosomeHalfSize) + chromosomeHalfSize
	log.
		Trace().
		Int("pointOne", pointOne).
		Int("pointTwo", pointTwo).
		Msg("partially mapped crossover")
	var childGenes = make([]int, chromosomeSize)
	for i := range chromosomeSize {
		if i >= pointOne && i < pointTwo {
			childGenes[i] = parentOne[i]
		} else {
			childGenes[i] = -1
		}
	}
	log.
		Trace().
		Ints("childGenes", childGenes).
		Msg("generating child (step 1)")
	for i := pointOne; i < pointTwo; i++ {
		if !slices.Contains(childGenes, parentTwo[i]) {
			position := findPosition(i, parentOne, parentTwo, childGenes)
			childGenes[position] = parentTwo[i]
		}
	}
	log.
		Trace().
		Ints("childGenes", childGenes).
		Msg("generating child (step 2)")
	for i := range chromosomeSize {
		if childGenes[i] == -1 {
			childGenes[i] = parentTwo[i]
		}
	}
	log.
		Trace().
		Ints("childGenes", childGenes).
		Msg("generating child (step 3)")
	return childGenes
}

func findPosition(index int, parentOne []int, parentTwo []int, child []int) int {
	position := -1
	for i := range parentOne {
		if parentTwo[i] == parentOne[index] {
			position = i
			break
		}
	}
	if child[position] != -1 {
		return findPosition(position, parentOne, parentTwo, child)
	}
	return position
}

func (ga GeneticAlgorithm) RunAlgorithm() Chromosome {
	ga.calcFitness()
	epochCounter := 0
	for {
		epochCounter += 1
		ga.mateRandomChromosomes(10, 50)
		ga.calcFitness()
		bestConflictsSum := ga.getBestChromosome().conflictsSum
		log.
			Info().
			Int("populationSize", len(ga.population)).
			Int("epochCounter", epochCounter).
			Int("bestConflictsSum", bestConflictsSum).
			Msg("running ga epoch")
		if bestConflictsSum == 0 {
			return ga.getBestChromosome()
		}
		if epochCounter > 5000 {
			return ga.getBestChromosome()
		}
	}
}

func BuildGeneticAlgorithm(size int, initialPopulation int) GeneticAlgorithm {
	population := make([]Chromosome, initialPopulation)
	for i := range initialPopulation {
		positions := GenerateDistinctRandomValues(size)
		chromosome := NewChromosome(positions)
		population[i] = *chromosome
	}
	return GeneticAlgorithm{population}
}
