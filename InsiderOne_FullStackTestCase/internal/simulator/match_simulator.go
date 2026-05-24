package simulator

import (
	"math"
	"math/rand"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
)

type MatchSimulator interface {
	Simulate(home domain.Team, away domain.Team) domain.MatchResult
}

type BasicMatchSimulator struct {
	HomeAdvantage float64
}

func NewMatchSimulator() MatchSimulator {
	return &BasicMatchSimulator{HomeAdvantage: 0.1}
}

func (s *BasicMatchSimulator) Simulate(home domain.Team, away domain.Team) domain.MatchResult {
	strengthDiff := float64(home.Strength-away.Strength) / 100.0

	homeExpected := 1.2 + strengthDiff + s.HomeAdvantage
	awayExpected := 1.2 - strengthDiff

	if homeExpected < 0.3 {
		homeExpected = 0.3
	}
	if awayExpected < 0.3 {
		awayExpected = 0.3
	}

	return domain.MatchResult{
		HomeGoals: poisson(homeExpected),
		AwayGoals: poisson(awayExpected),
	}
}

// poisson generates a random number following a Poisson distribution (Knuth algorithm)
func poisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}
