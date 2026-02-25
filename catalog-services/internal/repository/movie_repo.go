package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrigangka2003/bms/catalog-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type MovieRepo struct {
	db *pgxpool.Pool
}

// new movie creates a new repo
func NewMovieRepo(db *pgxpool.Pool) *MovieRepo {
	return &MovieRepo{db: db}
}

//create inserts a movie into DB

func (r *MovieRepo) Create(ctx context.Context, movie *models.Movie) error {
	query := `
		INSERT INTO movies (title,description,duration_mins,release_date)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at
	`
	row := r.db.QueryRow(ctx, query, movie.Title, movie.Description, movie.Duration, movie.ReleaseDate)
	
	err:= row.Scan(&movie.ID, &movie.CreatedAt)

	return err
}

// GetAll fetches all movies from the DB 
func (r *MovieRepo) GetAll(ctx context.Context) ([]models.Movie, error) {
	query := `SELECT id, title, description, duration_mins, release_date, created_at FROM movies`
	
	// Use "Query" instead of "QueryRow" because we expect MULTIPLE rows (a list of movies).
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err // If the query failed, return nothing (nil) and the error.
	}
	defer rows.Close() // ALWAYS close rows when done so we don't leak memory.

	// Create an empty list (slice) to hold the movies we find.
	var movies []models.Movie

	// Loop through the results, one row at a time.
	for rows.Next() {
		var m models.Movie // Create an empty movie struct
		
		// Copy the columns from the DB row into the empty struct variables.
		// MUST be in the exact same order as the SELECT query above!
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseDate, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		// Add the filled-out movie struct to our list
		movies = append(movies, m)
	}

	return movies, nil // Return the full list of movies!
}


//Get Movies By Id 

func (r *MovieRepo) GetById(ctx context.Context, id string) (*models.Movie, error) {
	query := `
		SELECT id, title, description, duration_mins, release_date, created_at 
		FROM movies 
		WHERE id = $1
	`

	var m models.Movie 

	// Execute the query and scan the result into our struct
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseDate, &m.CreatedAt,
	)
	
	// If the database says "I couldn't find any rows with that ID"
	if err == pgx.ErrNoRows {
		return nil, nil // Return nil, meaning "No movie found, but no system error"
	}

	return &m, nil
}