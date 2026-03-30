package main

import (
	"flag"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strconv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setLogLevelFromFlags() {
	trace := flag.Bool("trace", false, "sets log level to trace")
	debug := flag.Bool("debug", false, "sets log level to debug")
	warn := flag.Bool("warn", false, "sets log level to warn")
	error := flag.Bool("error", false, "sets log level to error")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *trace {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *warn {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *error {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func configureLogger() {
	setLogLevelFromFlags()

	logFile := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

func main() {
	configureLogger()

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
