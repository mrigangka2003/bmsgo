package routes

import (
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/handler"
)

// RegisterRoutes creates a new ServeMux and registers all versioned API routes.
func RegisterRoutes(
	movieHandler *handler.MovieHandler,
	theaterHandler *handler.TheaterHandler,
	showHandler *handler.ShowHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check (unversioned)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("catalog service is healthy and connected to the db"))
	})


	registerV1Routes(mux, movieHandler, theaterHandler, showHandler)

	return mux
}

// registerV1Routes groups all /api/v1 routes.
func registerV1Routes(
	mux *http.ServeMux,
	movieHandler *handler.MovieHandler,
	theaterHandler *handler.TheaterHandler,
	showHandler *handler.ShowHandler,
) {
	// Movie Routes
	mux.HandleFunc("POST /api/v1/movies", movieHandler.CreateMovie)
	mux.HandleFunc("GET /api/v1/movies", movieHandler.GetMovies)
	mux.HandleFunc("GET /api/v1/movies/{id}", movieHandler.GetMovieById)

	// Theater Routes
	mux.HandleFunc("POST /api/v1/theaters", theaterHandler.CreateTheater)
	mux.HandleFunc("GET /api/v1/theaters", theaterHandler.GetTheaters)

	// Show Routes
	mux.HandleFunc("POST /api/v1/shows", showHandler.CreateShow)
	mux.HandleFunc("GET /api/v1/movies/{id}/shows", showHandler.GetShowsForMovie)
}
