package domain

type Match struct {
	ID         int
	Week       int
	HomeTeamID int
	AwayTeamID int
	HomeGoals  *int
	AwayGoals  *int
	IsPlayed   bool
}

type MatchResult struct {
	HomeGoals int
	AwayGoals int
}
