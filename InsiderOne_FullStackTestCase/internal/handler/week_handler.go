package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
)

type WeekHandler struct {
	leagueService *service.LeagueService
}

func NewWeekHandler(leagueService *service.LeagueService) *WeekHandler {
	return &WeekHandler{leagueService: leagueService}
}

func (h *WeekHandler) PlayNextWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	matches, standings, err := h.leagueService.PlayNextWeek()
	if err != nil {
		http.Error(w, `{"error":"failed to play next week"}`, http.StatusInternalServerError)
		return
	}
	if matches == nil {
		writeJSON(w, http.StatusOK, map[string]string{"message": "all weeks have been played"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"week":      matches[0].Week,
		"matches":   matches,
		"standings": standings,
	})
}
