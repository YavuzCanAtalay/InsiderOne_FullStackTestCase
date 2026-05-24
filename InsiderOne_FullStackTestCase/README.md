# InsiderOne_FullStackTestCase

Go backend API for a 4-team football league simulation. The application stores teams and matches in PostgreSQL, simulates weekly results, calculates live standings, and estimates championship probabilities with Monte Carlo simulation.

## Features

- 4-team league with a 6-week double round-robin fixture list
- Weekly match simulation based on team strength
- Dynamic standings calculation using Premier League-style rules
- Play-next-week and play-full-league flows
- Championship predictions (with percentage probabilities) included in week responses from week 4 onward
- Manual match result editing — standings update immediately, no extra step needed
- Matches returned grouped as first leg (weeks 1–3) and second leg (weeks 4–6)
- Dockerized API and PostgreSQL setup

## Tech Stack

- Go 1.22
- PostgreSQL 15
- Docker / Docker Compose

## Project Structure

```text
cmd/api                  application entry point
internal/domain          shared data models
internal/repository      database access layer
internal/service         business logic
internal/handler         HTTP handlers
internal/simulator       match simulation logic
internal/prediction      Monte Carlo prediction engine
migrations               schema and seed SQL
docs                     reports and supporting documentation
```

## Prerequisites

Choose one of these:

- Docker and Docker Compose
- Go 1.22+ and PostgreSQL 15+

## Environment Variables

The project expects these variables:

```env
DB_USER=football_user
DB_PASSWORD=football_pass
DB_NAME=football_league
DB_HOST=db
DB_PORT=5432
API_PORT=8080
```

If you do not already have a `.env` file, create one from `.env.example`.

## Quick Start With Docker

From the project root:

```bash
docker compose up --build
```

This starts:

- `db`: PostgreSQL
- `api`: Go HTTP server on port `8080`

The SQL files under `migrations/` are mounted into the PostgreSQL container and run automatically on first startup.

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

### Stop The Project

```bash
docker compose down
```

## Run Without Docker

1. Start a local PostgreSQL instance.
2. Create a database matching your `.env` values.
3. Run the SQL files in `migrations/001_create_tables.sql` and `migrations/002_seed_data.sql`.
4. If the API runs on your machine instead of Docker, set `DB_HOST=localhost`.
5. Start the server:

```bash
go run ./cmd/api
```

## API Endpoints

### Health

```http
GET /health
```

### Teams

```http
GET /teams
```

### Matches

```http
GET /matches
GET /matches/week/{week}
PUT /matches/{id}
```

`GET /matches` returns matches grouped into two legs:

```json
{
  "first_leg":  [ { "week": 1, "matches": [...] }, ... ],
  "second_leg": [ { "week": 4, "matches": [...] }, ... ]
}
```

Edit a played match result:

```bash
curl -X PUT http://localhost:8080/matches/1 \
  -H "Content-Type: application/json" \
  -d '{"HomeGoals":2,"AwayGoals":1}'
```

Standings recalculate automatically — no additional request needed.

### Standings

```http
GET /standings
```

### League Simulation

```http
POST /weeks/next
POST /league/play-all
```

`POST /weeks/next` simulates the next unplayed week and returns the results, updated standings, and — from week 4 onward — championship predictions:

```bash
curl -X POST http://localhost:8080/weeks/next
```

`POST /league/play-all` simulates all remaining weeks at once. If all matches are already played it returns:

```json
{"message": "all matches are already played"}
```

### Predictions

```http
GET /predictions/current
```

Returns championship probability (formatted as a percentage, e.g. `"73.450%"`) and expected final position for each team. Available only after week 4 has been played. Predictions are also included automatically in the `POST /weeks/next` response from week 4 onward.

## Live Deployment

The project is deployed and publicly accessible on Railway:

**Base URL:** `https://insideronefullstacktestcase-production.up.railway.app`

```bash
curl https://insideronefullstacktestcase-production.up.railway.app/health
curl https://insideronefullstacktestcase-production.up.railway.app/teams
curl https://insideronefullstacktestcase-production.up.railway.app/matches
curl https://insideronefullstacktestcase-production.up.railway.app/standings
curl -X POST https://insideronefullstacktestcase-production.up.railway.app/weeks/next
curl https://insideronefullstacktestcase-production.up.railway.app/predictions/current
```

## Example Verification Flow

After starting the project locally, these requests are a good basic check:

```bash
curl http://localhost:8080/teams
curl http://localhost:8080/matches
curl http://localhost:8080/standings
curl -X POST http://localhost:8080/weeks/next
curl http://localhost:8080/predictions/current
```

## Notes

- Standings are calculated dynamically from match results. There is no standings table in the database.
- Team strength influences the simulator, but results are still randomized.
- The `predictions` SQL table exists in the schema, while current predictions are generated on demand by the API.
