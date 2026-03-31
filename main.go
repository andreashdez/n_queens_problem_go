package main

import (
	"flag"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func setLogLevelFromFlags() {
	trace := flag.Bool("trace", false, "sets log level to trace")
	debug := flag.Bool("debug", false, "sets log level to debug")
	warn := flag.Bool("warn", false, "sets log level to warn")
	error := flag.Bool("error", false, "sets log level to error")
	flag.Parse()

	logLevel := zerolog.InfoLevel
	switch {
	case *trace:
		logLevel = zerolog.TraceLevel
	case *debug:
		logLevel = zerolog.DebugLevel
	case *warn:
		logLevel = zerolog.WarnLevel
	case *error:
		logLevel = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(logLevel)
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
