package service

import (
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/simulator"
)

type LeagueService struct { // playing next week, whole season, getting standings
	teamRepo  repository.TeamRepository
	matchRepo repository.MatchRepository
	simulator simulator.MatchSimulator
}

func NewLeagueService(
	teamRepo repository.TeamRepository,
	matchRepo repository.MatchRepository,
	sim simulator.MatchSimulator,
) *LeagueService {
	return &LeagueService{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
		simulator: sim,
	}
}

// PlayNextWeek simulates the next unplayed week and returns results + updated standings
func (s *LeagueService) PlayNextWeek() ([]domain.Match, []domain.Standing, error) {
	unplayed, err := s.matchRepo.GetUnplayed()
	if err != nil {
		return nil, nil, err
	}
	if len(unplayed) == 0 {
		return nil, nil, nil
	}

	nextWeek := unplayed[0].Week // first unplayed match's week number
	var weekMatches []domain.Match
	for _, m := range unplayed {
		if m.Week == nextWeek {
			weekMatches = append(weekMatches, m)
		}
	}

	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, nil, err
	}

	teamMap := make(map[int]domain.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	for i := range weekMatches {
		result := s.simulator.Simulate(teamMap[weekMatches[i].HomeTeamID], teamMap[weekMatches[i].AwayTeamID])
		if err := s.matchRepo.UpdateResult(weekMatches[i].ID, result); err != nil {
			return nil, nil, err
		}
		hg, ag := result.HomeGoals, result.AwayGoals
		weekMatches[i].HomeGoals = &hg
		weekMatches[i].AwayGoals = &ag
		weekMatches[i].IsPlayed = true
	}

	allMatches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, nil, err
	}

	standings := CalculateStandings(teams, allMatches)
	return weekMatches, standings, nil
}

// PlayAll simulates all remaining weeks and returns results grouped by week
func (s *LeagueService) PlayAll() (map[int][]domain.Match, []domain.Standing, error) {
	results := make(map[int][]domain.Match)

	for {
		weekMatches, standings, err := s.PlayNextWeek()
		if err != nil {
			return nil, nil, err
		}
		if weekMatches == nil {
			allMatches, _ := s.matchRepo.GetAll()
			teams, _ := s.teamRepo.GetAll()
			return results, CalculateStandings(teams, allMatches), nil
		}
		results[weekMatches[0].Week] = weekMatches
		_ = standings
	}
}

// GetStandings returns the current league table
func (s *LeagueService) GetStandings() ([]domain.Standing, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return CalculateStandings(teams, matches), nil
}
