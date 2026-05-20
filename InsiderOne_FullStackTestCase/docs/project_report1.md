# Football League Simulation — Project Report

---

## internal/domain/ — Core Data Types

### team.go
| Field | Description |
|---|---|
| `ID` | Database row identifier |
| `Name` | "Chelsea", "Arsenal", etc. |
| `Strength` | 70–90, controls how often the team wins |

### match.go
| Field | Description |
|---|---|
| `ID` | Database row identifier |
| `Week` | Which week (1–6) this match belongs to |
| `HomeTeamID` | Foreign key → teams.id |
| `AwayTeamID` | Foreign key → teams.id |
| `HomeGoals` | `*int` pointer — nil until match is played |
| `AwayGoals` | `*int` pointer — nil until match is played |
| `IsPlayed` | false until simulator runs |

**MatchResult**
| Field | Description |
|---|---|
| `HomeGoals` | Plain int, returned by simulator |
| `AwayGoals` | Plain int, returned by simulator |

### standing.go
| Field | Description |
|---|---|
| `TeamID / TeamName` | Which team this row belongs to |
| `Played` | Matches played so far |
| `Won / Drawn / Lost` | Match outcomes |
| `GoalsFor` | Total goals scored |
| `GoalsAgainst` | Total goals conceded |
| `GoalDifference` | GoalsFor - GoalsAgainst |
| `Points` | Win=3, Draw=1, Loss=0 |

### prediction.go
| Field | Description |
|---|---|
| `TeamID / TeamName` | Which team |
| `ChampionshipProbability` | e.g. 45.2 (%) |
| `ExpectedFinalPosition` | e.g. 1.3 (average finishing position from Monte Carlo) |

---

## migrations/ — Database Setup

### 001_create_tables.sql
| Table | Purpose |
|---|---|
| `teams` | Stores 4 teams with strength ratings |
| `matches` | Stores all 12 scheduled matches; goals filled in after simulation |
| `predictions` | Stores championship % per team per week |

### 002_seed_data.sql
| Data | Detail |
|---|---|
| 4 teams | Chelsea (90), Arsenal (85), Manchester City (80), Liverpool (70) |
| 12 matches | Full double round-robin across 6 weeks. Weeks 1–3: first leg. Weeks 4–6: reverse fixtures |

---

## internal/repository/ — Database Access

### team_repository.go

**TeamRepository (interface)**
| Function | SQL | Returns |
|---|---|---|
| `GetAll()` | SELECT all 4 teams | `[]domain.Team` |
| `GetByID(id)` | SELECT one team by id | `domain.Team` |

- `postgresTeamRepository` — holds a `*sql.DB` connection and executes queries
- `NewTeamRepository(db)` — wires the implementation to the interface

### match_repository.go

**MatchRepository (interface)**
| Function | SQL | Returns |
|---|---|---|
| `GetAll()` | SELECT all 12 matches | `[]domain.Match` |
| `GetByWeek(week)` | SELECT matches WHERE week = ? | `[]domain.Match` |
| `GetUnplayed()` | SELECT matches WHERE is_played = FALSE | `[]domain.Match` |
| `UpdateResult(id, res)` | UPDATE home_goals, away_goals, is_played=TRUE | `error` |

- `scanMatches(rows)` — shared helper that reads SQL rows into `[]domain.Match`
- `NewMatchRepository(db)` — wires the implementation to the interface

---

## internal/simulator/ — Match Score Generator

### match_simulator.go

**MatchSimulator (interface)**
```
Simulate(home, away Team) → MatchResult
```

**BasicMatchSimulator**
| Field | Value | Purpose |
|---|---|---|
| `HomeAdvantage` | 0.1 | Small boost for home team |

**Simulate(home, away) logic**
```
strengthDiff = (home.Strength - away.Strength) / 100
homeExpected = 1.2 + strengthDiff + 0.1
awayExpected = 1.2 - strengthDiff
→ Poisson(homeExpected) → HomeGoals
→ Poisson(awayExpected) → AwayGoals
```

**poisson(lambda)**
Knuth algorithm — generates a random integer following a Poisson distribution.
Used because real football goals follow a Poisson distribution.

---

## internal/service/ — Business Logic

### standings_service.go

**CalculateStandings(teams, matches) → []domain.Standing**
- Loops through all played matches
- Tallies W/D/L, goals, and points for each team
- Sorts by Premier League rules: Points → Goal Difference → Goals Scored
- Never touches the database — pure calculation from match data

### league_service.go

**LeagueService struct**
| Field | Purpose |
|---|---|
| `teamRepo` | Fetches teams from DB |
| `matchRepo` | Fetches and updates matches in DB |
| `simulator` | Generates match scores |

**Functions**
| Function | What it does |
|---|---|
| `PlayNextWeek()` | Finds next unplayed week → simulates matches → saves to DB → returns results + table |
| `PlayAll()` | Loops PlayNextWeek() until no matches remain → returns all results by week + final standings |
| `GetStandings()` | Fetches all matches → calculates and returns current league table |

---

## Still To Write

| Phase | Files | Purpose |
|---|---|---|
| 6 | `prediction/prediction_engine.go` | Monte Carlo championship probability |
| 7 | 6 handler files | HTTP endpoints |
| 8 | `cmd/api/main.go` (update) | Wire everything together |
