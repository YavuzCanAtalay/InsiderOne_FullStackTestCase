package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/prediction"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
)

type WeekHandler struct {
	leagueService  *service.LeagueService
	predictionEngine prediction.PredictionEngine
}

func NewWeekHandler(leagueService *service.LeagueService, predEngine prediction.PredictionEngine) *WeekHandler {
	return &WeekHandler{leagueService: leagueService, predictionEngine: predEngine}
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

	week := matches[0].Week
	resp := map[string]any{
		"week":      week,
		"matches":   matches,
		"standings": standings,
	}

	if week >= 4 {
		predictions, err := h.predictionEngine.Predict(week)
		if err == nil {
			resp["predictions"] = predictions
		}
	}

	writeJSON(w, http.StatusOK, resp)
}
