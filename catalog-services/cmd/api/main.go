package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrigangka2003/bms/catalog-service/internal/config"
	"github.com/mrigangka2003/bms/catalog-service/internal/database"
	"github.com/mrigangka2003/bms/catalog-service/internal/handler"
	"github.com/mrigangka2003/bms/catalog-service/internal/middleware"
	"github.com/mrigangka2003/bms/catalog-service/internal/repository"
	"github.com/mrigangka2003/bms/catalog-service/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	dbPool := database.ConnectDB(cfg.DATABASE_URL)

	// Initialize repositories
	movieRepo := repository.NewMovieRepo(dbPool)
	theaterRepo := repository.NewTheaterRepo(dbPool)
	showRepo := repository.NewShowRepo(dbPool)

	// Initialize handlers
	movieHandler := handler.NewMovieHandler(movieRepo)
	theaterHandler := handler.NewTheaterHandler(theaterRepo)
	showHandler := handler.NewShowHandler(showRepo)

	// Register routes
	mux := routes.RegisterRoutes(movieHandler, theaterHandler, showHandler)

	// Wrap with CORS middleware
	handlerWithCORS := middleware.CORS(mux)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handlerWithCORS,
	}

	// Start server in background goroutine
	go func() {
		log.Printf("Catalog Service starting on port %s...\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v\n", err)
	}

	// Close DB after HTTP shutdown
	dbPool.Close()

	log.Println("Database connection closed.")
	log.Println("Server exited cleanly.")
}
