# internal/prediction

## Purpose
Championship probability engine. Activated after Week 4.

## What goes here
- `prediction_engine.go` — implements the `PredictionEngine` interface using Monte Carlo simulation

## What needs to be done
- Read current standings and remaining unplayed matches
- Run at least 1,000 simulated season completions
- Count how often each team finishes first
- Return championship probability (%) for each team
- This should reuse the same `MatchSimulator` used for real weeks — do not duplicate simulation logic
