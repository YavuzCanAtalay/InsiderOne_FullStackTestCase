package domain

type Standing struct {
	TeamID         int
	TeamName       string
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Points         int
}
