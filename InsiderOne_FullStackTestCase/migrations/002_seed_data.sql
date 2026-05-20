-- 002_seed_data.sql

INSERT INTO teams (name, strength) VALUES
    ('Chelsea',          90),
    ('Arsenal',          85),
    ('Manchester City',  80),
    ('Liverpool',        70);

-- Double round-robin: 4 teams, 6 weeks, 2 matches per week, 12 matches total
-- Team IDs: Chelsea=1, Arsenal=2, Manchester City=3, Liverpool=4

INSERT INTO matches (week, home_team_id, away_team_id) VALUES
    -- Week 1
    (1, 1, 2),  -- Chelsea vs Arsenal
    (1, 3, 4),  -- Manchester City vs Liverpool

    -- Week 2
    (2, 1, 3),  -- Chelsea vs Manchester City
    (2, 2, 4),  -- Arsenal vs Liverpool

    -- Week 3
    (3, 1, 4),  -- Chelsea vs Liverpool
    (3, 2, 3),  -- Arsenal vs Manchester City

    -- Week 4 (reverse fixtures)
    (4, 2, 1),  -- Arsenal vs Chelsea
    (4, 4, 3),  -- Liverpool vs Manchester City

    -- Week 5
    (5, 3, 1),  -- Manchester City vs Chelsea
    (5, 4, 2),  -- Liverpool vs Arsenal

    -- Week 6
    (6, 4, 1),  -- Liverpool vs Chelsea
    (6, 3, 2);  -- Manchester City vs Arsenal
