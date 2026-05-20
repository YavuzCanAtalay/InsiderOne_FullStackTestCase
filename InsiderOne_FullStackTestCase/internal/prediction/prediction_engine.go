package prediction

import (
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/simulator"
)

const simulations = 10000

type PredictionEngine interface {
	Predict(currentWeek int) ([]domain.Prediction, error)
}

type MonteCarloPredictionEngine struct {
	teamRepo  repository.TeamRepository
	matchRepo repository.MatchRepository
	simulator simulator.MatchSimulator
}

func NewPredictionEngine(
	teamRepo repository.TeamRepository,
	matchRepo repository.MatchRepository,
	sim simulator.MatchSimulator,
) PredictionEngine {
	return &MonteCarloPredictionEngine{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
		simulator: sim,
	}
}

func (e *MonteCarloPredictionEngine) Predict(currentWeek int) ([]domain.Prediction, error) {
	teams, err := e.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	unplayed, err := e.matchRepo.GetUnplayed()
	if err != nil {
		return nil, err
	}

	playedMatches, err := e.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	teamMap := make(map[int]domain.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	// championCount[teamID] = how many times this team won the league across all simulations
	championCount := make(map[int]int)
	positionSum := make(map[int]float64)

	for i := 0; i < simulations; i++ {
		simMatches := simulateSeason(playedMatches, unplayed, teamMap, e.simulator)
		standings := service.CalculateStandings(teams, simMatches)

		for pos, s := range standings {
			positionSum[s.TeamID] += float64(pos + 1)
			if pos == 0 {
				championCount[s.TeamID]++
			}
		}
	}

	predictions := make([]domain.Prediction, 0, len(teams))
	for _, t := range teams {
		predictions = append(predictions, domain.Prediction{
			TeamID:                  t.ID,
			TeamName:                t.Name,
			ChampionshipProbability: float64(championCount[t.ID]) / float64(simulations) * 100,
			ExpectedFinalPosition:   positionSum[t.ID] / float64(simulations),
		})
	}

	return predictions, nil
}

// simulateSeason takes already-played matches and simulates all remaining ones
func simulateSeason(
	played []domain.Match,
	unplayed []domain.Match,
	teamMap map[int]domain.Team,
	sim simulator.MatchSimulator,
) []domain.Match {
	all := make([]domain.Match, len(played))
	copy(all, played)

	for _, m := range unplayed {
		result := sim.Simulate(teamMap[m.HomeTeamID], teamMap[m.AwayTeamID])
		hg, ag := result.HomeGoals, result.AwayGoals
		m.HomeGoals = &hg
		m.AwayGoals = &ag
		m.IsPlayed = true
		all = append(all, m)
	}

	return all
}
