package main

import (
	"flag"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RuntimeConfig struct {
	boardSize         int
	initialPopulation int
	maxEpochs         int
	minToMate         int
	maxToMate         int
	mutationRate      float64
	seed              int64
	logLevel          zerolog.Level
}

func parseFlags() RuntimeConfig {
	traceLogLevel := flag.Bool("trace", false, "sets log level to trace")
	debugLogLevel := flag.Bool("debug", false, "sets log level to debug")
	warnLogLevel := flag.Bool("warn", false, "sets log level to warn")
	errorLogLevel := flag.Bool("error", false, "sets log level to error")
	boardSize := flag.Int("size", 14, "number of queens and board size")
	initialPopulation := flag.Int("population", 40000, "initial population size")
	maxEpochs := flag.Int("max-epochs", 5000, "maximum number of epochs to run")
	minToMate := flag.Int("min-to-mate", 10, "minimum chromosomes to mate per epoch")
	maxToMate := flag.Int("max-to-mate", 50, "maximum chromosomes to mate per epoch")
	mutationRate := flag.Float64("mutation-rate", 0.03, "mutation probability per offspring (0.0 to 1.0)")
	seed := flag.Int64("seed", -1, "random seed; set negative value to use current time")
	flag.Parse()

	logLevel := zerolog.InfoLevel
	switch {
	case *traceLogLevel:
		logLevel = zerolog.TraceLevel
	case *debugLogLevel:
		logLevel = zerolog.DebugLevel
	case *warnLogLevel:
		logLevel = zerolog.WarnLevel
	case *errorLogLevel:
		logLevel = zerolog.ErrorLevel
	}

	return RuntimeConfig{
		boardSize:         *boardSize,
		initialPopulation: *initialPopulation,
		maxEpochs:         *maxEpochs,
		minToMate:         *minToMate,
		maxToMate:         *maxToMate,
		mutationRate:      *mutationRate,
		seed:              *seed,
		logLevel:          logLevel,
	}
}

func validateRuntimeConfig(config RuntimeConfig) {
	if config.boardSize <= 0 {
		log.Fatal().Int("size", config.boardSize).Msg("size must be greater than zero")
	}
	if config.initialPopulation <= 0 {
		log.Fatal().Int("population", config.initialPopulation).Msg("population must be greater than zero")
	}
	if config.maxEpochs <= 0 {
		log.Fatal().Int("maxEpochs", config.maxEpochs).Msg("max epochs must be greater than zero")
	}
	if config.minToMate <= 0 {
		log.Fatal().Int("minToMate", config.minToMate).Msg("minimum mates must be greater than zero")
	}
	if config.maxToMate < config.minToMate {
		log.
			Fatal().
			Int("minToMate", config.minToMate).
			Int("maxToMate", config.maxToMate).
			Msg("maximum mates must be greater than or equal to minimum mates")
	}
	if config.mutationRate < 0 || config.mutationRate > 1 {
		log.
			Fatal().
			Float64("mutationRate", config.mutationRate).
			Msg("mutation rate must be between 0 and 1")
	}
}

func configureLogger(logLevel zerolog.Level) {
	zerolog.SetGlobalLevel(logLevel)

	logFile := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	}
	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

func main() {
	config := parseFlags()
	validateRuntimeConfig(config)
	configureLogger(config.logLevel)

	seed := config.seed
	if seed < 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)

	log.Info().Msg("start n_queens_problem")
	log.
		Info().
		Int64("seed", seed).
		Float64("mutationRate", config.mutationRate).
		Msg("configured random seed and mutation rate")
	ga := BuildGeneticAlgorithm(
		config.boardSize,
		config.initialPopulation,
		config.maxEpochs,
		config.minToMate,
		config.maxToMate,
		config.mutationRate,
	)
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
