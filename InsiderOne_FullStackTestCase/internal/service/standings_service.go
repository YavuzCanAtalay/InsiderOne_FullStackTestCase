package service

import (
	"sort"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
)

func CalculateStandings(teams []domain.Team, matches []domain.Match) []domain.Standing {
	index := make(map[int]*domain.Standing)

	for _, t := range teams {
		t := t
		index[t.ID] = &domain.Standing{
			TeamID:   t.ID,
			TeamName: t.Name,
		}
	}

	for _, m := range matches {
		if !m.IsPlayed || m.HomeGoals == nil || m.AwayGoals == nil {
			continue
		}

		home := index[m.HomeTeamID]
		away := index[m.AwayTeamID]
		hg, ag := *m.HomeGoals, *m.AwayGoals

		home.Played++
		away.Played++
		home.GoalsFor += hg
		home.GoalsAgainst += ag
		away.GoalsFor += ag
		away.GoalsAgainst += hg
		home.GoalDifference = home.GoalsFor - home.GoalsAgainst
		away.GoalDifference = away.GoalsFor - away.GoalsAgainst

		switch {
		case hg > ag:
			home.Won++
			home.Points += 3
			away.Lost++
		case ag > hg:
			away.Won++
			away.Points += 3
			home.Lost++
		default:
			home.Drawn++
			away.Drawn++
			home.Points++
			away.Points++
		}
	}

	standings := make([]domain.Standing, 0, len(teams))
	for _, s := range index {
		standings = append(standings, *s)
	}

	sort.Slice(standings, func(i, j int) bool {
		a, b := standings[i], standings[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDifference != b.GoalDifference {
			return a.GoalDifference > b.GoalDifference
		}
		return a.GoalsFor > b.GoalsFor
	})

	return standings
}
