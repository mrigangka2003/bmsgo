package models

import "time"

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration_mins"` 
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Theater struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Location string `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type Show struct {
	ID string `json:"id"`
	MovieID string `json:"movie_id"`
	TheaterID string `json:"theater_id"`
	StartTime time.Time `json:"start_time"`
	Price int `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ShowDetails struct {
	ID          string    `json:"id"`
	MovieTitle  string    `json:"movie_title"`
	TheaterName string    `json:"theater_name"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	Price       int       `json:"price"`
}
