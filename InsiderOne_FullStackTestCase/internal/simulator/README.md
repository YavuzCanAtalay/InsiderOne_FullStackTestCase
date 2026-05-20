# internal/simulator

## Purpose
Match simulation logic. Generates realistic match scores based on team strengths.

## What goes here
- `match_simulator.go` — implements the `MatchSimulator` interface
- Uses home/away team strength ratings and a home advantage bonus to produce a random but weighted scoreline

## What needs to be done
- Stronger teams should win more often, but not always
- Add a small home team advantage
- Generate goal counts using weighted random logic (e.g. Poisson-style distribution)
- Keep this stateless — it receives a match and returns a result, nothing else
