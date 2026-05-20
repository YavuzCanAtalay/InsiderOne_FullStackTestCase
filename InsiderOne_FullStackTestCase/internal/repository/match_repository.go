package repository

import (
	"database/sql"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
)

type MatchRepository interface {
	GetAll() ([]domain.Match, error)
	GetByWeek(week int) ([]domain.Match, error)
	GetUnplayed() ([]domain.Match, error)
	UpdateResult(matchID int, result domain.MatchResult) error
}

type postgresMatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) MatchRepository {
	return &postgresMatchRepository{db: db}
}

func (r *postgresMatchRepository) GetAll() ([]domain.Match, error) {
	rows, err := r.db.Query(
		"SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, is_played FROM matches ORDER BY week, id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMatches(rows)
}

func (r *postgresMatchRepository) GetByWeek(week int) ([]domain.Match, error) {
	rows, err := r.db.Query(
		"SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, is_played FROM matches WHERE week = $1 ORDER BY id",
		week,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMatches(rows)
}

func (r *postgresMatchRepository) GetUnplayed() ([]domain.Match, error) {
	rows, err := r.db.Query(
		"SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, is_played FROM matches WHERE is_played = FALSE ORDER BY week, id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMatches(rows)
}

func (r *postgresMatchRepository) UpdateResult(matchID int, result domain.MatchResult) error {
	_, err := r.db.Exec(
		"UPDATE matches SET home_goals = $1, away_goals = $2, is_played = TRUE WHERE id = $3",
		result.HomeGoals, result.AwayGoals, matchID,
	)
	return err
}

func scanMatches(rows *sql.Rows) ([]domain.Match, error) {
	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		if err := rows.Scan(&m.ID, &m.Week, &m.HomeTeamID, &m.AwayTeamID, &m.HomeGoals, &m.AwayGoals, &m.IsPlayed); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}
