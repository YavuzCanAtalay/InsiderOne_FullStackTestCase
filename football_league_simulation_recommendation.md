# Football League Simulation Project — Recommended Tools and Progress Roadmap

## 1. Project Goal

This project is a **Go backend football league simulator and prediction system**.  
The case expects the system to:

- Simulate a league with **4 football teams**
- Show **weekly match results**
- Update the **league table** after each week
- Generate a **final league table prediction after Week 4**
- Be implemented in **GoLang**
- Use **interface-based design** and **struct composition**
- Expose clear **REST endpoints** that can be tested via Postman
- Include **SQL schema and queries**
- Ideally include documentation, deployment, and bonus features

---

## 2. Recommended Tools

### Core Development Tools

| Tool | Purpose |
|---|---|
| **GoLang** | Mandatory programming language for the case |
| **Go Modules** | Dependency and package management |
| **VS Code / GoLand** | Development environment |
| **Git + GitHub** | Version control and code handover |
| **Postman** | Testing REST endpoints |
| **PostgreSQL or MySQL** | Relational database for SQL schema and queries |
| **Docker + Docker Compose** | Reproducible setup and deployment |

### Recommended Stack

```text
Language: Go
API: net/http or Gin
Database: PostgreSQL
Database Access: database/sql
Testing: Go testing package
API Testing: Postman
Deployment: Docker Compose
```

---

## 3. Recommended League Structure

Because the task examples show the league progressing through **Week 4** and **Week 5**, the cleanest structure is:

- **4 teams**
- **Double round-robin format**
- Each team plays every other team **twice**
- Total:
  - **6 weeks**
  - **2 matches per week**
  - **12 matches overall**

### Suggested Teams

You may use real or fictional teams. A practical option:

| Team | Strength Rating |
|---|---:|
| Chelsea | 90 |
| Arsenal | 85 |
| Manchester City | 80 |
| Liverpool | 70 |

The `strength_rating` can influence match simulation and prediction results.

---

## 4. League Table Rules

Use Premier League-style table calculations:

1. **Points**
2. **Goal Difference**
3. **Goals Scored**
4. Optional tie-breakers if you want to be more complete:
   - Head-to-head points
   - Head-to-head away goals

### Points Rule

| Result | Points |
|---|---:|
| Win | 3 |
| Draw | 1 |
| Loss | 0 |

---

## 5. Recommended System Design

The case explicitly asks for:

- **Interface-based design**
- **Struct composition**

So the architecture should avoid putting all logic into one file or one service.

### Core Domain Models

```go
type Team struct {
    ID             int
    Name           string
    StrengthRating int
}

type Match struct {
    ID         int
    Week       int
    HomeTeamID int
    AwayTeamID int
    HomeGoals  *int
    AwayGoals  *int
    IsPlayed   bool
}

type Standing struct {
    TeamID        int
    TeamName      string
    Played        int
    Won           int
    Drawn         int
    Lost          int
    GoalsFor      int
    GoalsAgainst  int
    GoalDifference int
    Points        int
}

type Prediction struct {
    TeamID                  int
    TeamName                string
    ChampionshipProbability float64
    ExpectedFinalPosition   float64
}
```

---

## 6. Recommended Interfaces

```go
type MatchSimulator interface {
    Simulate(match Match) MatchResult
}

type StandingsCalculator interface {
    Calculate(matches []Match, teams []Team) []Standing
}

type PredictionEngine interface {
    PredictFinalTable(currentWeek int) PredictionResult
}

type MatchRepository interface {
    GetByWeek(week int) ([]Match, error)
    GetAll() ([]Match, error)
    UpdateResult(matchID int, result MatchResult) error
}

type TeamRepository interface {
    GetAll() ([]Team, error)
}
```

This makes the project modular, testable, and aligned with the case requirements.

---

## 7. Suggested Folder Structure

```text
/cmd/api
/internal/domain
/internal/repository
/internal/service
/internal/simulator
/internal/prediction
/internal/handler
/migrations
/docs
```

### Example Responsibility Breakdown

| Folder | Responsibility |
|---|---|
| `domain` | Structs and core types |
| `repository` | Database access |
| `service` | Business logic |
| `simulator` | Match generation logic |
| `prediction` | Monte Carlo simulation engine |
| `handler` | HTTP endpoints |
| `migrations` | SQL schema |
| `docs` | API examples / diagrams / Postman notes |

---

## 8. Database Design

### `teams` Table

```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    strength_rating INT NOT NULL
);
```

### `matches` Table

```sql
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    week INT NOT NULL,
    home_team_id INT NOT NULL REFERENCES teams(id),
    away_team_id INT NOT NULL REFERENCES teams(id),
    home_goals INT,
    away_goals INT,
    is_played BOOLEAN NOT NULL DEFAULT FALSE
);
```

### `predictions` Table

```sql
CREATE TABLE predictions (
    id SERIAL PRIMARY KEY,
    week INT NOT NULL,
    team_id INT NOT NULL REFERENCES teams(id),
    championship_probability DECIMAL(5,2) NOT NULL,
    expected_final_position DECIMAL(4,2) NOT NULL
);
```

---

## 9. Important Data Design Recommendation

Do **not** treat a stored standings table as the source of truth.

Instead:

- Store:
  - teams
  - match schedule
  - match results
- Dynamically calculate standings from match results

### Why?

This makes it much easier to:

- Edit past results
- Recalculate standings correctly
- Re-run prediction logic
- Avoid inconsistent table data

This is especially helpful for the bonus requirement:  
**“Edit the results of the matches and calculate the edited results based on the modified standings.”**

---

## 10. Recommended API Endpoints

### Basic Reading Endpoints

```http
GET /teams
GET /matches
GET /matches/week/{week}
GET /standings
```

### Progress the League Week by Week

```http
POST /weeks/next
```

This endpoint should:

1. Find the next unplayed week
2. Simulate both matches of that week
3. Save results to the database
4. Recalculate standings
5. Return updated match results and league table
6. If the week is `>= 4`, return predictions as well

### Prediction Endpoints

```http
GET /predictions/current
GET /predictions/week/{week}
```

### Bonus Endpoint — Play the Rest of the League

```http
POST /league/play-all
```

This should:

- Automatically simulate all remaining weeks
- Return weekly match results
- Return the final standings
- Return updated predictions if relevant

### Bonus Endpoint — Edit Match Results

```http
PUT /matches/{id}
```

Example request body:

```json
{
  "home_goals": 2,
  "away_goals": 1
}
```

After editing, the system should:

- Update the match result
- Recalculate league standings
- Recalculate predictions if the match belongs to Week 4 or later

---

## 11. Match Simulation Logic

The case allows match results to depend on team strengths.

### Recommended Simple Model

Each team has a strength score:

```text
Chelsea: 90
Arsenal: 85
Manchester City: 80
Liverpool: 70
```

Then when simulating a match:

- The stronger team should be more likely to win
- The home team may receive a small bonus
- Goals can be produced through weighted random generation

### Recommended Inputs for the Simulator

```text
home_strength
away_strength
home_bonus
randomness_factor
```

### Example Design

```go
type BasicMatchSimulator struct {
    HomeAdvantage int
}
```

This type can implement:

```go
func (s BasicMatchSimulator) Simulate(match Match) MatchResult
```

---

## 12. Prediction Engine Recommendation

The prediction after Week 4 should be implemented using **Monte Carlo Simulation**.

### Why Monte Carlo?

It makes the project more realistic and demonstrates actual prediction logic instead of returning a hardcoded table.

### Suggested Flow

After Week 4:

1. Read current standings
2. Read remaining unplayed matches
3. Simulate the rest of the league many times:
   - 1,000 runs minimum
   - 10,000 runs if performance is fine
4. Count:
   - How often each team becomes champion
   - Average final position
   - Optional: average final points

### Example Prediction Response

```json
{
  "week": 4,
  "championship_probabilities": [
    {
      "team": "Chelsea",
      "probability": 45.2
    },
    {
      "team": "Arsenal",
      "probability": 28.4
    },
    {
      "team": "Manchester City",
      "probability": 20.1
    },
    {
      "team": "Liverpool",
      "probability": 6.3
    }
  ]
}
```

---

## 13. Recommended Progress Roadmap

# Phase 1 — Project Setup

### Goal
Create a clean backend skeleton.

### Tasks

- Install Go
- Initialize project with `go mod init`
- Create GitHub repository
- Create folder structure
- Add `.env` support if needed
- Set up database connection
- Create Docker Compose for backend + database

---

# Phase 2 — Domain Models and SQL Schema

### Goal
Create the core data layer.

### Tasks

- Define:
  - `Team`
  - `Match`
  - `Standing`
  - `Prediction`
- Create SQL migration files
- Add seed data:
  - 4 teams
  - 12 scheduled matches
- Implement basic repositories:
  - TeamRepository
  - MatchRepository

---

# Phase 3 — Standings Calculator

### Goal
Correctly compute the league table.

### Tasks

- Calculate:
  - Played
  - Wins
  - Draws
  - Losses
  - Goals For
  - Goals Against
  - Goal Difference
  - Points
- Sort table by league rules
- Write unit tests for standings calculation

### Priority
Very high.  
This is one of the most important parts of the project.

---

# Phase 4 — Weekly Simulation

### Goal
Play the league week by week.

### Tasks

- Implement match simulator
- Add:

```http
POST /weeks/next
```

- Simulate the next unplayed week
- Save results
- Return:
  - Match results
  - Updated table

---

# Phase 5 — Prediction Engine

### Goal
Produce final league estimates after Week 4.

### Tasks

- Build Monte Carlo simulation engine
- Run many season-completion simulations
- Compute:
  - Championship probability
  - Expected final position
- Add prediction response after Week 4
- Add:

```http
GET /predictions/current
```

---

# Phase 6 — Clean REST API Layer

### Goal
Make the project easy to evaluate.

### Tasks

- Implement endpoint handlers
- Add consistent JSON responses
- Add useful HTTP status codes
- Add structured error messages
- Create Postman collection
- Document sample requests and responses

---

# Phase 7 — Bonus Features

### Bonus 1: Play All Remaining Weeks

```http
POST /league/play-all
```

### Bonus 2: Edit Match Results

```http
PUT /matches/{id}
```

### Tasks

- Update match score
- Recalculate table
- Recalculate predictions
- Ensure edited results persist in the DB

---

# Phase 8 — Documentation and Deployment

### Goal
Submit a polished, reviewer-friendly project.

### Tasks

- Write a clear `README.md`
- Add:
  - Setup instructions
  - Database initialization instructions
  - Docker run instructions
  - Endpoint descriptions
  - Sample requests/responses
- Include:
  - SQL schema
  - Seed queries
  - Postman collection
- Optional:
  - Deploy the backend and share a live URL

---

## 14. Priority Order

If time becomes limited, follow this order:

1. Correct **league table calculations**
2. Weekly simulation endpoint
3. Prediction engine after Week 4
4. Interface-based Go architecture
5. SQL schema and seed queries
6. Postman collection
7. Play All endpoint
8. Edit Match Result endpoint
9. Docker and deployment

Do **not** start with deployment or bonus polish before the core system works correctly.

---

## 15. Recommended Final Deliverables

Your final handover should include:

```text
1. Go backend source code
2. SQL schema and seed queries
3. README.md with setup instructions
4. Postman collection
5. Docker Compose setup
6. Optional deployment URL
7. Unit tests for key business logic
```

---

## 16. Final Recommendation

The best version of this project is not the largest one. It is the one that is:

- Correct
- Modular
- Easy to test
- Easy to understand
- Clearly aligned with the case requirements

The features that will make your submission stand out most are:

1. **A clean Go architecture using interfaces**
2. **Correct, dynamically calculated standings**
3. **Monte Carlo-based championship predictions**
4. **Play All and Edit Result bonus endpoints**
5. **Excellent README + Postman + Docker handover**
