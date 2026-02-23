package main

import (
	"log"
	"net/http"

	"github.com/mrigangka2003/bms/catalog-service/internal/config"
	"github.com/mrigangka2003/bms/catalog-service/internal/database"
)

func main() {
	cfg := config.LoadConfig()

	dbPool := database.ConnectDB(cfg.DATABASE_URL)
	defer dbPool.Close()

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("catalog service is healthy and connected to the db"))
	})


	addr:= ":" + cfg.Port 
	log.Printf("🚀 Catalog Service starting on port %s...\n", cfg.Port)
	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
