-- 002_seed_data.sql

INSERT INTO teams (name, strength) VALUES
    ('Chelsea',          90),
    ('Arsenal',          85),
    ('Manchester City',  80),
    ('Liverpool',        70)
ON CONFLICT (name) DO NOTHING;

-- Double round-robin: 4 teams, 6 weeks, 2 matches per week, 12 matches total
-- Team IDs: Chelsea=1, Arsenal=2, Manchester City=3, Liverpool=4

INSERT INTO matches (week, home_team_id, away_team_id)
SELECT v.week, v.home_team_id, v.away_team_id
FROM (VALUES
    (1, 1, 2),
    (1, 3, 4),
    (2, 1, 3),
    (2, 2, 4),
    (3, 1, 4),
    (3, 2, 3),
    (4, 2, 1),
    (4, 4, 3),
    (5, 3, 1),
    (5, 4, 2),
    (6, 4, 1),
    (6, 3, 2)
) AS v(week, home_team_id, away_team_id)
WHERE NOT EXISTS (SELECT 1 FROM matches);
