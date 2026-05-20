package domain

type Prediction struct {
	TeamID                  int
	TeamName                string
	ChampionshipProbability float64
	ExpectedFinalPosition   float64
}
