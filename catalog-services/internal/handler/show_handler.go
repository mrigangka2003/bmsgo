package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/models"
	"github.com/mrigangka2003/bms/catalog-service/internal/repository"
	"github.com/mrigangka2003/bms/catalog-service/internal/utils"
)

type ShowHandler struct {
	repo *repository.ShowRepo
}

func NewShowHandler(repo *repository.ShowRepo) *ShowHandler {
	return &ShowHandler{repo: repo}
}

// POST /shows
func (h *ShowHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	var show models.Show

	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Because of Postgres Foreign Keys, if the MovieID or TheaterID doesn't exist, 
	// this repo.Create function will safely return an error!
	if err := h.repo.Create(r.Context(), &show); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to create show. Make sure movie_id and theater_id are valid.")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, show)
}


// GET /movies/{id}/shows
func (h *ShowHandler) GetShowsForMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("id")
	if movieID == "" {
		utils.ErrorJSON(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	// USE OUR NEW JOIN FUNCTION HERE!
	shows, err := h.repo.GetShowDetails(r.Context(), movieID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch shows: " + err.Error())
		return
	}

	if shows == nil {
		shows = []models.ShowDetails{}
	}

	utils.WriteJSON(w, http.StatusOK, shows)
}


