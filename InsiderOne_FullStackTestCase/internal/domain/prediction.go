package domain

type Prediction struct {
	TeamID                  int
	TeamName                string
	ChampionshipProbability string
	ExpectedFinalPosition   float64
}
