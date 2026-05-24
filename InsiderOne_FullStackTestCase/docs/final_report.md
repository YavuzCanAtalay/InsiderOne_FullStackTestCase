# Football League Simulation — Final Project Report

---

## Project Overview

A Go backend REST API that simulates a 4-team football league.
The system plays matches week by week, maintains a live league table,
and predicts championship probabilities using Monte Carlo simulation.
All data is persisted in PostgreSQL. No frontend — tested via Postman or curl.

---

## Architecture

```
HTTP Request
    ↓
Handler         → parses request, calls service, returns JSON
    ↓
Service         → business logic, orchestrates repos + simulator
    ↓
Repository      → SQL queries, reads/writes PostgreSQL
    ↓
PostgreSQL DB

Simulator       → called by service to generate match scores
Prediction      → called by handler, runs 10,000 Monte Carlo seasons
Domain          → shared structs used by all layers
```

---

## File-by-File Reference

---

### cmd/api/main.go

**Entry point of the application.**

Runs at startup in this order:
1. Connects to PostgreSQL using env variables from `.env`
2. Creates repositories (teamRepo, matchRepo)
3. Creates the match simulator
4. Creates the league service (wires repos + simulator together)
5. Creates the prediction engine
6. Creates all HTTP handlers
7. Registers all routes on the HTTP mux
8. Starts the server on port 8080

| Route | Method | Handler |
|---|---|---|
| `/health` | GET | inline — returns `{"status":"ok"}` |
| `/teams` | GET | TeamHandler |
| `/matches` | GET | MatchHandler |
| `/matches/week/{n}` | GET | MatchHandler |
| `/matches/{id}` | PUT | MatchHandler |
| `/standings` | GET | StandingsHandler |
| `/weeks/next` | POST | WeekHandler |
| `/league/play-all` | POST | LeagueHandler |
| `/predictions/current` | GET | PredictionHandler |

---

### internal/domain/

Pure data types. No logic, no imports. Used by every other package.

**team.go**
```
Team
  ID       → database primary key
  Name     → team name (e.g. "Chelsea")
  Strength → integer 70–90, drives simulation outcome
```

**match.go**
```
Match
  ID         → database primary key
  Week       → which week (1–6) this match belongs to
  HomeTeamID → references teams.id
  AwayTeamID → references teams.id
  HomeGoals  → *int, nil until match is simulated
  AwayGoals  → *int, nil until match is simulated
  IsPlayed   → false until simulator runs

MatchResult
  HomeGoals  → plain int returned by simulator
  AwayGoals  → plain int returned by simulator
```

**standing.go**
```
Standing
  TeamID / TeamName  → which team
  Played             → matches played so far
  Won / Drawn / Lost → match outcomes
  GoalsFor           → total goals scored
  GoalsAgainst       → total goals conceded
  GoalDifference     → GoalsFor - GoalsAgainst
  Points             → Win=3, Draw=1, Loss=0
```

**prediction.go**
```
Prediction
  TeamID / TeamName          → which team
  ChampionshipProbability    → percentage (e.g. 45.2)
  ExpectedFinalPosition      → average finishing position from simulations
```

---

### migrations/

SQL files executed automatically by PostgreSQL on first Docker startup.

**001_create_tables.sql**

Creates three tables:
- `teams` — id, name, strength
- `matches` — id, week, home_team_id, away_team_id, home_goals, away_goals, is_played
- `predictions` — id, week, team_id, championship_probability, expected_final_position

Note: standings are NOT stored — they are always calculated dynamically from match results.

**002_seed_data.sql**

Inserts initial data:
- 4 teams: Chelsea (90), Arsenal (85), Manchester City (80), Liverpool (70)
- 12 matches: full double round-robin schedule across 6 weeks
  - Weeks 1–3: first leg fixtures
  - Weeks 4–6: reverse fixtures

---

### internal/repository/

The only layer that talks to the database. All SQL lives here.

**team_repository.go**
```
TeamRepository (interface)
  GetAll()      → SELECT all teams → []domain.Team
  GetByID(id)   → SELECT one team  → domain.Team

NewTeamRepository(db) → returns interface implementation
```

**match_repository.go**
```
MatchRepository (interface)
  GetAll()                    → SELECT all 12 matches
  GetByWeek(week)             → SELECT matches for one week
  GetUnplayed()               → SELECT matches WHERE is_played = FALSE
  UpdateResult(id, result)    → UPDATE goals + set is_played = TRUE

scanMatches(rows)             → shared helper, maps SQL rows → []domain.Match

NewMatchRepository(db) → returns interface implementation
```

---

### internal/simulator/

Generates realistic match scores based on team strength. No database access.

**match_simulator.go**
```
MatchSimulator (interface)
  Simulate(home, away Team) → MatchResult

BasicMatchSimulator
  HomeAdvantage = 0.1

Simulate logic:
  strengthDiff = (home.Strength - away.Strength) / 100
  homeExpected = 1.2 + strengthDiff + 0.1
  awayExpected = 1.2 - strengthDiff
  HomeGoals = Poisson(homeExpected)
  AwayGoals = Poisson(awayExpected)

poisson(lambda)
  Knuth algorithm — generates random goals following Poisson distribution
  Produces realistic football scorelines (0-0, 1-0, 2-1, etc.)

NewMatchSimulator() → returns ready-to-use simulator
```

---

### internal/service/

Business logic layer. Orchestrates repositories and simulator.

**standings_service.go**
```
CalculateStandings(teams, matches) → []domain.Standing
  Loops through all played matches
  Tallies wins, draws, losses, goals for each team
  Sorts by Premier League rules:
    1. Points (descending)
    2. Goal Difference (descending)
    3. Goals Scored (descending)
  Pure calculation — never touches the database
```

**league_service.go**
```
LeagueService
  teamRepo  → fetch teams from DB
  matchRepo → fetch and update matches in DB
  simulator → generate scores

PlayNextWeek()
  1. Fetch all unplayed matches
  2. Identify the lowest unplayed week number
  3. Fetch both teams for each match in that week
  4. Call simulator.Simulate() → get score
  5. Call matchRepo.UpdateResult() → save to DB
  6. Recalculate standings from all matches
  7. Return: week's match results + updated standings

PlayAll()
  Loops PlayNextWeek() until no matches remain
  Returns: all results grouped by week + final standings

GetStandings()
  Fetches all matches from DB
  Returns: current league table via CalculateStandings()
```

---

### internal/prediction/

Monte Carlo engine that estimates championship probabilities.
Only activated after Week 4 has been played.

**prediction_engine.go**
```
PredictionEngine (interface)
  Predict(currentWeek) → []domain.Prediction

MonteCarloPredictionEngine
  teamRepo  → get all teams
  matchRepo → get played + unplayed matches
  simulator → simulate remaining matches

Predict(currentWeek)
  Runs 10,000 simulated season completions:
    Each run:
      → Keeps already-played matches fixed
      → Simulates all remaining matches randomly
      → Calculates final standings
      → Records who finished 1st
  After 10,000 runs:
    championship_probability = (times team finished 1st / 10000) × 100
    expected_final_position  = average finishing position across all runs

simulateSeason(played, unplayed, teamMap, sim)
  Helper — combines fixed results with randomly simulated remaining matches
  Returns a complete set of 12 matches for standings calculation

NewPredictionEngine(teamRepo, matchRepo, sim) → returns engine
```

---

### internal/handler/

HTTP layer. Parses requests, calls services, returns JSON responses.
No business logic here.

**helpers.go**
```
writeJSON(w, status, data)
  Sets Content-Type: application/json
  Writes HTTP status code
  Encodes data as JSON to response body
  Used by all handlers
```

**team_handler.go**
```
GET /teams
  Calls teamRepo.GetAll()
  Returns all 4 teams as JSON array
```

**match_handler.go**
```
GET /matches
  Calls matchRepo.GetAll()
  Returns all 12 matches

GET /matches/week/{n}
  Parses week number from URL
  Calls matchRepo.GetByWeek(n)
  Returns matches for that week

PUT /matches/{id}
  Parses match id from URL
  Decodes {home_goals, away_goals} from request body
  Calls matchRepo.UpdateResult()
  Returns {"status": "updated"}
```

**standings_handler.go**
```
GET /standings
  Calls leagueService.GetStandings()
  Returns current league table sorted by points
```

**week_handler.go**
```
POST /weeks/next
  Calls leagueService.PlayNextWeek()
  Returns:
    - week number
    - match results for that week
    - updated standings table
  If all weeks are played: returns a message saying so
```

**league_handler.go**
```
POST /league/play-all
  Calls leagueService.PlayAll()
  Returns:
    - all match results grouped by week
    - final league standings
```

**prediction_handler.go**
```
GET /predictions/current
  Checks how many weeks have been played
  If < 4 weeks played: returns message asking to play more weeks
  If >= 4 weeks played:
    Calls predictionEngine.Predict(currentWeek)
    Returns championship probabilities and expected positions for all teams
```

---

## Docker Setup

**Dockerfile**
Two-stage build:
1. Builder stage — uses `golang:1.22-alpine`, downloads deps, compiles binary
2. Final stage — copies only the binary into `alpine:latest` (small image)

**docker-compose.yml**
Two services:
- `db` — PostgreSQL 15, runs migration SQL files on first start, health-checked
- `api` — Go binary, waits for DB to be healthy before starting

**`.env`**
```
DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT, API_PORT
```
Never committed to git (listed in `.gitignore`).

---

## Key Design Decisions

| Decision | Reason |
|---|---|
| Standings calculated dynamically | Editing a past result automatically corrects the whole table |
| Interfaces for repository and simulator | Enforces modular design, makes testing possible |
| Poisson distribution for goals | Statistically accurate model for football scoring |
| Monte Carlo with 10,000 runs | Balances prediction accuracy with performance |
| No standings table in DB | Single source of truth — match results only |
