package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrigangka2003/bms/catalog-service/internal/models"
)

type TheaterRepo struct {
	db *pgxpool.Pool
}

func NewTheaterRepo(db *pgxpool.Pool) *TheaterRepo {
	return &TheaterRepo{db: db}
}

// inserts a new theater into the database
func (r *TheaterRepo) Create(ctx context.Context, theater *models.Theater) error {
	query := `
	INSERT INTO theaters (name, location)
	VALUES ($1, $2)
	RETURNING id ,created_at
	`
	err := r.db.QueryRow(ctx, query, theater.Name, theater.Location).Scan(&theater.ID, &theater.CreatedAt)
	return err
}

func (r *TheaterRepo) GetAll(ctx context.Context) ([]models.Theater, error) {
	query := `SELECT id, name, location, created_at FROM theaters`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theaters []models.Theater
	for rows.Next() {
		var t models.Theater
		err := rows.Scan(&t.ID, &t.Name, &t.Location, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		theaters = append(theaters, t)
	}

	return theaters, nil
}