# migrations

## Purpose
SQL schema and seed data. Required deliverable for the case submission.

## What goes here
- `001_create_tables.sql` — creates `teams`, `matches`, and `predictions` tables
- `002_seed_data.sql` — inserts the 4 teams and schedules all 12 matches (6 weeks × 2 matches)

## What needs to be done
- Define `teams` table: id, name, strength_rating
- Define `matches` table: id, week, home_team_id, away_team_id, home_goals, away_goals, is_played
- Define `predictions` table: id, week, team_id, championship_probability
- Seed 4 teams with strength ratings
- Seed 12 pre-scheduled matches covering all 6 weeks of the double round-robin
