# internal/domain

## Purpose
Core data types shared across the entire application.

## What goes here
- `team.go` — `Team` struct (ID, Name, StrengthRating)
- `match.go` — `Match` struct (ID, Week, HomeTeamID, AwayTeamID, HomeGoals, AwayGoals, IsPlayed)
- `standing.go` — `Standing` struct (Played, Won, Drawn, Lost, GoalsFor, GoalsAgainst, GoalDifference, Points)
- `prediction.go` — `Prediction` struct (TeamID, ChampionshipProbability, ExpectedFinalPosition)

## What needs to be done
- Define all structs with correct field types
- No business logic here — pure data models only
