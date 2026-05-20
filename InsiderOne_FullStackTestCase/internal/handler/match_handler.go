package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/domain"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
)

type MatchHandler struct {
	matchRepo repository.MatchRepository
}

func NewMatchHandler(matchRepo repository.MatchRepository) *MatchHandler {
	return &MatchHandler{matchRepo: matchRepo}
}

func (h *MatchHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.matchRepo.GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch matches"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

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
}

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
}
