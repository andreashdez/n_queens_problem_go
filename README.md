# N-Queens Problem (Go + Genetic Algorithm)

This project solves the N-Queens puzzle with a genetic algorithm.

- A chromosome is a permutation of row positions (one queen per column).
- Fitness is based on diagonal conflicts (lower is better, `0` is solved).
- The solver uses roulette selection, PMX crossover, optional mutation, and bounded population trimming.

## Requirements

- Go `1.26+`

## Quick Start

```bash
go run .
```

Run with custom parameters:

```bash
go run . -size 8 -population 1000 -max-epochs 1000 -min-to-mate 5 -max-to-mate 20 -mutation-rate 0.05 -seed 1234
```

## CLI Flags

| Flag             | Default | Description                                         |
| ---------------- | ------- | --------------------------------------------------- |
| `-size`          | `14`    | Number of queens and board size                     |
| `-population`    | `40000` | Initial population size                             |
| `-max-epochs`    | `5000`  | Maximum epochs before stopping                      |
| `-min-to-mate`   | `10`    | Minimum offspring count per epoch                   |
| `-max-to-mate`   | `50`    | Maximum offspring count per epoch                   |
| `-mutation-rate` | `0.03`  | Mutation probability per offspring (`0.0` to `1.0`) |
| `-seed`          | `-1`    | RNG seed (`< 0` uses current time)                  |
| `-trace`         | `false` | Set log level to trace                              |
| `-debug`         | `false` | Set log level to debug                              |
| `-warn`          | `false` | Set log level to warn                               |
| `-error`         | `false` | Set log level to error                              |

## Parameter Tuning Guide

Start with these presets, then adjust from there:

- Small boards (`N=8` to `N=16`)
  - `-population 500` to `2000`
  - `-max-epochs 500` to `2000`
  - `-min-to-mate 4` to `8`
  - `-max-to-mate 10` to `30`
  - `-mutation-rate 0.02` to `0.08`
- Larger boards (`N>=32`)
  - `-population 10000` to `50000`
  - `-max-epochs 3000` to `10000`
  - `-min-to-mate 10` to `30`
  - `-max-to-mate 30` to `120`
  - `-mutation-rate 0.03` to `0.12`

Rules of thumb:

- If the best conflict score plateaus early, increase `-mutation-rate` slightly.
- If progress is noisy, reduce `-mutation-rate` slightly and increase `-population`.
- If runs are too slow, reduce `-population` first, then `-max-epochs`.
- Keep `-max-to-mate` noticeably higher than `-min-to-mate` for diversity.
- For difficult sizes, run multiple seeds and keep the best result.

Example commands:

```bash
# Fast local run for 8x8
go run . -size 8 -population 1000 -max-epochs 1200 -min-to-mate 5 -max-to-mate 20 -mutation-rate 0.05 -seed 42

# Stronger run for 32x32
go run . -size 32 -population 20000 -max-epochs 6000 -min-to-mate 15 -max-to-mate 80 -mutation-rate 0.06 -seed 42
```

## Reproducible Runs

Use a fixed `-seed` to make runs deterministic.

```bash
go run . -size 14 -seed 2026
```

The configured seed and mutation rate are written to logs at startup.

## Output

- The terminal prints the final board.
- Each queen cell shows its conflict count (`00` means no conflicts).
- Structured logs are written to `app.log`.

Example solved `8x8` board:

```
╔════╤════╤════╤════╤════╤════╤════╤════╗
║    │    │    │    │    │ 00 │    │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║ 00 │    │    │    │    │    │    │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │    │    │    │ 00 │    │    │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │ 00 │    │    │    │    │    │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │    │    │    │    │    │    │ 00 ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │    │ 00 │    │    │    │    │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │    │    │    │    │    │ 00 │    ║
╟────┼────┼────┼────┼────┼────┼────┼────╢
║    │    │    │ 00 │    │    │    │    ║
╚════╧════╧════╧════╧════╧════╧════╧════╝
```

## Development

Run tests:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -run '^$' -bench PMX -benchmem ./...
go test -run '^$' -bench CountConflicts -benchmem ./...
```
