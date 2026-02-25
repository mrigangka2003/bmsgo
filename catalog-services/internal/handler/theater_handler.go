package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/models"
	"github.com/mrigangka2003/bms/catalog-service/internal/repository"
	"github.com/mrigangka2003/bms/catalog-service/internal/utils"
)

type TheaterHandler struct {
	repo *repository.TheaterRepo
}

func NewTheaterHandler(repo *repository.TheaterRepo) *TheaterHandler {
	return &TheaterHandler{repo: repo}
}

// POST /theaters
func (h *TheaterHandler) CreateTheater(w http.ResponseWriter, r *http.Request) {
	var theater models.Theater

	// 1. Read JSON
	if err := json.NewDecoder(r.Body).Decode(&theater); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// 2. Save to DB
	if err := h.repo.Create(r.Context(), &theater); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to save theater")
		return
	}

	// 3. Send Success JSON
	utils.WriteJSON(w, http.StatusCreated, theater)
}

// GET /theaters
func (h *TheaterHandler) GetTheaters(w http.ResponseWriter, r *http.Request) {
	theaters, err := h.repo.GetAll(r.Context())
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch theaters")
		return
	}

	if theaters == nil {
		theaters = []models.Theater{}
	}

	utils.WriteJSON(w, http.StatusOK, theaters)
}