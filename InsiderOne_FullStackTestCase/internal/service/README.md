# internal/service

## Purpose
Business logic layer. Orchestrates repositories, simulator, and prediction engine.

## What goes here
- `league_service.go` — coordinates playing a week, recalculating standings, triggering predictions
- `standings_service.go` — dynamically calculates the league table from match results (do NOT store standings in DB)

## What needs to be done
- Implement `POST /weeks/next` logic: find next unplayed week → simulate → save → return results
- Implement `POST /league/play-all` logic: loop through all remaining weeks
- Implement `PUT /matches/{id}` logic: update result → recalculate standings + predictions
- Standings must be computed dynamically from stored match results, not from a cached table
