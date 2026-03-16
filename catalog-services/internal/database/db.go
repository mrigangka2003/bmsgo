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
	theatersTable := `
	CREATE TABLE IF NOT EXISTS theaters (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		location TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	showsTable := `
	CREATE TABLE IF NOT EXISTS shows (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
		theater_id UUID NOT NULL REFERENCES theaters(id) ON DELETE CASCADE,
		start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		price INT NOT NULL, 
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	queries := []string{query, theatersTable, showsTable}

	for _, query := range queries {
		_, err := db.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to create tables: %v\n", err)
		}
	}

	log.Println("Database tables (movies, theaters,shows) verified/created successfully")
}
