package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/service"
)

type StandingsHandler struct { // stores a league service used to fetch standings data from db
	leagueService *service.LeagueService
}

func NewStandingsHandler(leagueService *service.LeagueService) *StandingsHandler {
	return &StandingsHandler{leagueService: leagueService}
}

func (h *StandingsHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	standings, err := h.leagueService.GetStandings()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch standings"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, standings)
} // fetches current league standings from league service and returns as JSON response
