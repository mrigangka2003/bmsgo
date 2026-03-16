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
func (r *MovieRepo) GetAll(ctx context.Context, page int, limit int) ([]models.Movie, error) {
	// Calculate how many records to skip
	// Example: Page 1 skips 0. Page 2 skips 10.
	offset := (page - 1) * limit

	// Notice the ORDER BY, LIMIT, and OFFSET! 
	// $1 is limit, $2 is offset.
	query := `
		SELECT id, title, description, duration_mins, release_date, created_at 
		FROM movies 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseDate, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
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