package models

import "time"

type Movie struct {
	Id            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Rate          float64   `json:"rate" db:"rate"`
	Cover         string    `json:"cover" db:"cover"`
	UserCreatedId int       `json:"userCreatedId" db:"user_created_id"`
	Date          time.Time `json:"date" db:"date"`
}

type AddMovieRequest struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Cover       string `json:"cover" db:"cover"`
	Date        string `json:"date" db:"date"`
}
