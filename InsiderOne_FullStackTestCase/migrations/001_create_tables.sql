-- 001_create_tables.sql

CREATE TABLE IF NOT EXISTS teams (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL UNIQUE,
    strength INT NOT NULL
);

CREATE TABLE IF NOT EXISTS matches (
    id           SERIAL PRIMARY KEY,
    week         INT NOT NULL,
    home_team_id INT NOT NULL REFERENCES teams(id),
    away_team_id INT NOT NULL REFERENCES teams(id),
    home_goals   INT,
    away_goals   INT,
    is_played    BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS predictions (
    id                       SERIAL PRIMARY KEY,
    week                     INT NOT NULL,
    team_id                  INT NOT NULL REFERENCES teams(id),
    championship_probability DECIMAL(5,2) NOT NULL,
    expected_final_position  DECIMAL(4,2) NOT NULL
);
