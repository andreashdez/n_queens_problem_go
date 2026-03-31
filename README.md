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
