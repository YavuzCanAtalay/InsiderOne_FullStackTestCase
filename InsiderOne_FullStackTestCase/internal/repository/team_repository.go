package repository

import (
	"database/sql"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
)

type TeamRepository interface {
	GetAll() ([]domain.Team, error)
	GetByID(id int) (domain.Team, error)
}

type postgresTeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &postgresTeamRepository{db: db}
}

func (r *postgresTeamRepository) GetAll() ([]domain.Team, error) {
	rows, err := r.db.Query("SELECT id, name, strength FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Strength); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *postgresTeamRepository) GetByID(id int) (domain.Team, error) {
	var t domain.Team
	err := r.db.QueryRow("SELECT id, name, strength FROM teams WHERE id = $1", id).
		Scan(&t.ID, &t.Name, &t.Strength)
	return t, err
}
