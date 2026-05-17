# internal/repository

## Purpose
Database access layer. All SQL queries live here.

## What goes here
- `team_repository.go` — implements `TeamRepository` interface (GetAll, GetByID)
- `match_repository.go` — implements `MatchRepository` interface (GetAll, GetByWeek, UpdateResult, GetUnplayed)
- Interfaces are defined here or in `domain` and implemented against a real PostgreSQL connection

## What needs to be done
- Write SQL queries using `database/sql`
- Return domain structs from query results
- Handle errors cleanly — no business logic here, only data in/out
