package main

import (
	"log"
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/config"
	"github.com/mrigangka2003/bms/catalog-service/internal/database"
	"github.com/mrigangka2003/bms/catalog-service/internal/handler"
	"github.com/mrigangka2003/bms/catalog-service/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	dbPool := database.ConnectDB(cfg.DATABASE_URL)
	defer dbPool.Close()

	repo := repository.NewMovieRepo(dbPool)
	theaterRepo := repository.NewTheaterRepo(dbPool)

	// handlers
	movieHandler := handler.NewMovieHandler(repo)
	theaterHandler := handler.NewTheaterHandler(theaterRepo)

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("catalog service is healthy and connected to the db"))
	})
  
  //Movie Routes
	http.HandleFunc("POST /movies", movieHandler.CreateMovie)
	http.HandleFunc("GET /movies", movieHandler.GetMovies)
	http.HandleFunc("GET /movies/{id}", movieHandler.GetMovieById)

	//Theater Routes
	http.HandleFunc("POST /theaters", theaterHandler.CreateTheater)
	http.HandleFunc("GET /theaters", theaterHandler.GetTheaters)

	addr:= ":" + cfg.Port 
	log.Printf("🚀 Catalog Service starting on port %s...\n", cfg.Port)
	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
