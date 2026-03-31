package main

import (
	"math/rand"
	"slices"
	"testing"
)

func TestRandomGeneration(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	result := GenerateDistinctRandomValues(rng, 8)
	if !slices.Contains(result, 0) {
		t.Fatalf(`random values missing value '0', error`)
	}
	if !slices.Contains(result, 1) {
		t.Fatalf(`random values missing value '1', error`)
	}
	if !slices.Contains(result, 2) {
		t.Fatalf(`random values missing value '2', error`)
	}
	if !slices.Contains(result, 3) {
		t.Fatalf(`random values missing value '3', error`)
	}
	if !slices.Contains(result, 4) {
		t.Fatalf(`random values missing value '4', error`)
	}
	if !slices.Contains(result, 5) {
		t.Fatalf(`random values missing value '5', error`)
	}
	if !slices.Contains(result, 6) {
		t.Fatalf(`random values missing value '6', error`)
	}
	if !slices.Contains(result, 7) {
		t.Fatalf(`random values missing value '7', error`)
	}
}

func TestCountConflicts(t *testing.T) {
	positions1 := []int{0, 2, 4, 6, 1, 3, 5, 7}
	chromosome1 := NewChromosome(positions1)
	conflictsSum1 := chromosome1.conflictsSum
	want1 := 1
	if conflictsSum1 != want1 {
		t.Fatalf(`%v conflictsSum = '%d', want '%d', error`, positions1, conflictsSum1, want1)
	}
	positions2 := []int{2, 4, 1, 7, 5, 0, 6, 3}
	chromosome2 := NewChromosome(positions2)
	conflictsSum2 := chromosome2.conflictsSum
	want2 := 2
	if conflictsSum2 != want2 {
		t.Fatalf(`%v conflictsSum = '%d', want '%d', error`, positions2, conflictsSum2, want2)
	}
}
