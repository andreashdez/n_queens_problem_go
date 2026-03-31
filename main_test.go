package main

import (
	"flag"
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func runSetLogLevelFromArgs(t *testing.T, args ...string) zerolog.Level {
	t.Helper()

	oldCommandLine := flag.CommandLine
	oldArgs := os.Args
	oldLevel := zerolog.GlobalLevel()

	programName := "test"
	if len(oldArgs) > 0 {
		programName = oldArgs[0]
	}

	flag.CommandLine = flag.NewFlagSet(programName, flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = append([]string{programName}, args...)

	setLogLevelFromFlags()
	got := zerolog.GlobalLevel()

	flag.CommandLine = oldCommandLine
	os.Args = oldArgs
	zerolog.SetGlobalLevel(oldLevel)

	return got
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
			got := runSetLogLevelFromArgs(t, tt.args...)
			if got != tt.want {
				t.Fatalf("setLogLevelFromFlags(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
