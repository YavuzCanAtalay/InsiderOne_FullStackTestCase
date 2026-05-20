package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/prediction"
	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
)

type PredictionHandler struct {
	engine    prediction.PredictionEngine
	matchRepo repository.MatchRepository
}

func NewPredictionHandler(engine prediction.PredictionEngine, matchRepo repository.MatchRepository) *PredictionHandler {
	return &PredictionHandler{engine: engine, matchRepo: matchRepo}
}

func (h *PredictionHandler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	played, err := h.matchRepo.GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch matches"}`, http.StatusInternalServerError)
		return
	}

	currentWeek := 0
	for _, m := range played {
		if m.IsPlayed && m.Week > currentWeek {
			currentWeek = m.Week
		}
	}

	if currentWeek < 4 {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "predictions are available after week 4",
		})
		return
	}

	predictions, err := h.engine.Predict(currentWeek)
	if err != nil {
		http.Error(w, `{"error":"failed to calculate predictions"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"week":        currentWeek,
		"predictions": predictions,
	})
}
