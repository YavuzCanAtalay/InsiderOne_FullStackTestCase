package handler

import (
	"net/http"

	"github.com/YavuzCanAtalay/InsiderOne_FullStackTestCase/internal/repository"
)

type TeamHandler struct { // stores a team repository used to fetch team data from db
	teamRepo repository.TeamRepository
}

func NewTeamHandler(teamRepo repository.TeamRepository) *TeamHandler {
	return &TeamHandler{teamRepo: teamRepo}
}

func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamRepo.GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch teams"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, teams)
}
