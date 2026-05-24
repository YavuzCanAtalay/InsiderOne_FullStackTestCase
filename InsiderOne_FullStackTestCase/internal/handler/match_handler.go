package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
)

type MatchHandler struct { // stores a match repository used to fetch match data from db
	matchRepo repository.MatchRepository
}

func NewMatchHandler(matchRepo repository.MatchRepository) *MatchHandler {
	return &MatchHandler{matchRepo: matchRepo}
} // fetches all matches from db and returns as JSON response

func (h *MatchHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.matchRepo.GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch matches"}`, http.StatusInternalServerError)
		return
	}

	byWeek := make(map[int][]domain.Match)
	for _, m := range matches {
		byWeek[m.Week] = append(byWeek[m.Week], m)
	}

	type weekGroup struct {
		Week    int            `json:"week"`
		Matches []domain.Match `json:"matches"`
	}

	var firstLeg, secondLeg []weekGroup
	for week := 1; week <= 6; week++ {
		g := weekGroup{Week: week, Matches: byWeek[week]}
		if week <= 3 {
			firstLeg = append(firstLeg, g)
		} else {
			secondLeg = append(secondLeg, g)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"first_leg":  firstLeg,
		"second_leg": secondLeg,
	})
} // fetches matches for a specific week, week number is read from URL path

func (h *MatchHandler) GetMatchesByWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := strings.TrimPrefix(r.URL.Path, "/matches/week/")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, `{"error":"invalid week number"}`, http.StatusBadRequest)
		return
	}

	matches, err := h.matchRepo.GetByWeek(week)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch matches"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, matches)
} // updates match result, match ID is read from URL path and new result is read from request body

func (h *MatchHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/matches/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid match id"}`, http.StatusBadRequest)
		return
	}

	var result domain.MatchResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.matchRepo.UpdateResult(id, result); err != nil {
		http.Error(w, `{"error":"failed to update match"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
} // reads match ID from URL path, new result from request body, updates match result in db and returns success status as JSON response
