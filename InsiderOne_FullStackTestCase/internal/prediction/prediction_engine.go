package prediction

import (
	"fmt"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/simulator"
)

const simulations = 10000

type PredictionEngine interface {
	Predict(currentWeek int) ([]domain.Prediction, error)
}

type MonteCarloPredictionEngine struct { // asks repositories for data
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
	teams, err := e.teamRepo.GetAll() // loads all teams
	if err != nil {
		return nil, err
	}

	unplayed, err := e.matchRepo.GetUnplayed() //loads unplayed matches
	if err != nil {
		return nil, err
	}

	playedMatches, err := e.matchRepo.GetAll() // loads all matches
	if err != nil {
		return nil, err
	}

	teamMap := make(map[int]domain.Team) // builds a map of teams; key is team ID and value is team struct, used for quick lookups during simulation
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	// championCount[teamID] = how many times this team won the league across all simulations
	championCount := make(map[int]int)
	positionSum := make(map[int]float64)

	for i := 0; i < simulations; i++ { // onte carlo simulation 
		// simulate season
		simMatches := simulateSeason(playedMatches, unplayed, teamMap, e.simulator)
		standings := service.CalculateStandings(teams, simMatches)
		// update stats
		for pos, s := range standings {
			positionSum[s.TeamID] += float64(pos + 1)
			if pos == 0 {
				championCount[s.TeamID]++
			}
		}
	}
	//predict championship probabilities and expected final positions based on simulation results
	predictions := make([]domain.Prediction, 0, len(teams))
	for _, t := range teams {
		predictions = append(predictions, domain.Prediction{
			TeamID:                  t.ID,
			TeamName:                t.Name,
			ChampionshipProbability: fmt.Sprintf("%.3f%%", float64(championCount[t.ID])/float64(simulations)*100),
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
