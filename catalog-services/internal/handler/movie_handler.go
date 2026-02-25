package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/models"
	"github.com/mrigangka2003/bms/catalog-service/internal/repository"
	"github.com/mrigangka2003/bms/catalog-service/internal/utils"
)

type MovieHandler struct {
	repo *repository.MovieRepo
}

func NewMovieHandler(repo *repository.MovieRepo) *MovieHandler {
	return &MovieHandler{repo: repo}
}

// POST /movies
func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	if err := h.repo.Create(r.Context(), &movie); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to save movie to database")
		return
	}


	utils.WriteJSON(w, http.StatusCreated, movie)
}

// GET /movies
func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.repo.GetAll(r.Context())
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch movies")
		return
	}

	if movies == nil {
		movies = []models.Movie{}
	}

	
	utils.WriteJSON(w, http.StatusOK, movies)
}

// GET /movies/{id}
func (h *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.ErrorJSON(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	movie, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Database error")
		return
	}

	if movie == nil {
		utils.ErrorJSON(w, http.StatusNotFound, "Movie not found")
		return
	}


	utils.WriteJSON(w, http.StatusOK, movie)
}
