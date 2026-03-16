package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrigangka2003/bms/catalog-service/internal/models"
)

type ShowRepo struct {
	db *pgxpool.Pool
}

func NewShowRepo(db *pgxpool.Pool) *ShowRepo {
	return &ShowRepo{db: db}
}

// Create inserts a new show linking a movie and a theater.
func (r *ShowRepo) Create(ctx context.Context, show *models.Show) error {
	query := `
		INSERT INTO shows (movie_id, theater_id, start_time, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, show.MovieID, show.TheaterID, show.StartTime, show.Price).
		Scan(&show.ID, &show.CreatedAt)

	return err
}

// GetShowsByMovie fetches all shows for a specific Movie ID.
func (r *ShowRepo) GetShowsByMovie(ctx context.Context, movieID string) ([]models.Show, error) {
	// Notice the WHERE clause: we only want shows for this exact movie!
	query := `
		SELECT id, movie_id, theater_id, start_time, price, created_at 
		FROM shows 
		WHERE movie_id = $1
		ORDER BY start_time ASC
	`

	rows, err := r.db.Query(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var s models.Show
		err := rows.Scan(&s.ID, &s.MovieID, &s.TheaterID, &s.StartTime, &s.Price, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}

	return shows, nil
}

// GetShowDetails fetches shows but ALSO fetches the Movie Title and Theater Name!
func (r *ShowRepo) GetShowDetails(ctx context.Context, movieID string) ([]models.ShowDetails, error) {
	// Look at this beautiful SQL JOIN!
	// 's' stands for shows, 'm' for movies, 't' for theaters.
	query := `
		SELECT 
			s.id, 
			m.title, 
			t.name, 
			t.location, 
			s.start_time, 
			s.price
		FROM shows s
		JOIN movies m ON s.movie_id = m.id
		JOIN theaters t ON s.theater_id = t.id
		WHERE s.movie_id = $1
		ORDER BY s.start_time ASC
	`

	rows, err := r.db.Query(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []models.ShowDetails
	for rows.Next() {
		var s models.ShowDetails
		// Must match the exact order of the SELECT statement above
		err := rows.Scan(&s.ID, &s.MovieTitle, &s.TheaterName, &s.Location, &s.StartTime, &s.Price)
		if err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}

	return shows, nil
}