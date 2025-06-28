package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	debug := flag.Bool("debug", false, "sets log level to debug")
	trace := flag.Bool("trace", false, "sets log level to trace")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *trace {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Info().Msg("start n_queens_problem")
	ga := BuildGeneticAlgorithm(10, 40000)
	log.Info().Msg("done building genetic algorithm")
	bestChromosome := ga.RunAlgorithm()

	log.
		Info().
		Int("conflictsSum", bestChromosome.conflictsSum).
		Msg("best chromosome conflicts sum")
	positions := bestChromosome.positions
	conflicts := bestChromosome.conflicts

	DrawBoard(positions, conflicts)
	log.Info().Msg("done n_queens_problem")
}
