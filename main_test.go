package main

import (
	"flag"
	"io"
	"math"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func runParseFlagsFromArgs(t *testing.T, args ...string) RuntimeConfig {
	t.Helper()

	oldCommandLine := flag.CommandLine
	oldArgs := os.Args

	programName := "test"
	if len(oldArgs) > 0 {
		programName = oldArgs[0]
	}

	flag.CommandLine = flag.NewFlagSet(programName, flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = append([]string{programName}, args...)

	config := parseFlags()

	flag.CommandLine = oldCommandLine
	os.Args = oldArgs

	return config
}

func TestSetLogLevelFromFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want zerolog.Level
	}{
		{name: "default info", args: nil, want: zerolog.InfoLevel},
		{name: "trace", args: []string{"-trace"}, want: zerolog.TraceLevel},
		{name: "debug", args: []string{"-debug"}, want: zerolog.DebugLevel},
		{name: "warn", args: []string{"-warn"}, want: zerolog.WarnLevel},
		{name: "error", args: []string{"-error"}, want: zerolog.ErrorLevel},
		{name: "precedence", args: []string{"-warn", "-error"}, want: zerolog.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := runParseFlagsFromArgs(t, tt.args...)
			got := config.logLevel
			if got != tt.want {
				t.Fatalf("setLogLevelFromFlags(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestRuntimeFlags(t *testing.T) {
	config := runParseFlagsFromArgs(
		t,
		"-size", "8",
		"-population", "1000",
		"-max-epochs", "120",
		"-min-to-mate", "5",
		"-max-to-mate", "15",
		"-mutation-rate", "0.05",
		"-seed", "1234",
	)

	if config.boardSize != 8 {
		t.Fatalf("boardSize = %d, want 8", config.boardSize)
	}
	if config.initialPopulation != 1000 {
		t.Fatalf("initialPopulation = %d, want 1000", config.initialPopulation)
	}
	if config.maxEpochs != 120 {
		t.Fatalf("maxEpochs = %d, want 120", config.maxEpochs)
	}
	if config.minToMate != 5 {
		t.Fatalf("minToMate = %d, want 5", config.minToMate)
	}
	if config.maxToMate != 15 {
		t.Fatalf("maxToMate = %d, want 15", config.maxToMate)
	}
	if math.Abs(config.mutationRate-0.05) > 1e-9 {
		t.Fatalf("mutationRate = %f, want 0.05", config.mutationRate)
	}
	if config.seed != 1234 {
		t.Fatalf("seed = %d, want 1234", config.seed)
	}
}
