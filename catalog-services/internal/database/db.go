package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbUrl string) *pgxpool.Pool {
	// context.Background() is used for initialization steps that don't have a timeout
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v\n", err)
	}

	// Ping the database to ensure the credentials are correct and it is actually running
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database ping failed! Is PostgreSQL running? Error: %v\n", err)
	}

	log.Println("Successfully connected to PostgreSQL Pool!")

	RunMigrations(pool)
	return pool
}

func RunMigrations(db *pgxpool.Pool) {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		duration_mins INT NOT NULL,
		release_date DATE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v\n", err)
	}

	log.Println("database table verified/created successfully")
}
