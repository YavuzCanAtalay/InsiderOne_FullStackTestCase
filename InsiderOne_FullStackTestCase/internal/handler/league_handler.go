package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
)

type LeagueHandler struct {
	leagueService *service.LeagueService
}

func NewLeagueHandler(leagueService *service.LeagueService) *LeagueHandler {
	return &LeagueHandler{leagueService: leagueService}
}

func (h *LeagueHandler) PlayAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	results, standings, err := h.leagueService.PlayAll()
	if err != nil {
		http.Error(w, `{"error":"failed to play all weeks"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"results":   results,
		"standings": standings,
	})
}
